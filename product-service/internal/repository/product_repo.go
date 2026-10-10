package repository

import (
	"context"
	"errors"
	"log"

	"com.ecommerce/product-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) InitSchema(ctx context.Context) {
	if r.db == nil {
		return
	}
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			category_id SERIAL PRIMARY KEY,
			category_title VARCHAR(255) NOT NULL,
			image_url VARCHAR(255),
			parent_category_id INT REFERENCES categories(category_id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS products (
			product_id SERIAL PRIMARY KEY,
			product_title VARCHAR(255) NOT NULL,
			image_url VARCHAR(255),
			sku VARCHAR(100) UNIQUE,
			price_unit DECIMAL(12, 2) DEFAULT 0.0,
			quantity INT DEFAULT 0,
			category_id INT REFERENCES categories(category_id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := r.db.Exec(ctx, query); err != nil {
		log.Printf("[Product Service] Notice during schema verification: %v", err)
	}
}

// =================== Categories ===================

func (r *ProductRepository) FindAllCategories(ctx context.Context) ([]model.CategoryDto, error) {
	query := `SELECT category_id, category_title, COALESCE(image_url, '') FROM categories ORDER BY category_id ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CategoryDto
	for rows.Next() {
		var c model.CategoryDto
		if err := rows.Scan(&c.CategoryId, &c.CategoryTitle, &c.ImageUrl); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if list == nil {
		list = []model.CategoryDto{}
	}
	return list, nil
}

func (r *ProductRepository) FindCategoryByID(ctx context.Context, id int) (*model.CategoryDto, error) {
	query := `SELECT category_id, category_title, COALESCE(image_url, '') FROM categories WHERE category_id = $1`
	var c model.CategoryDto
	err := r.db.QueryRow(ctx, query, id).Scan(&c.CategoryId, &c.CategoryTitle, &c.ImageUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ProductRepository) SaveCategory(ctx context.Context, c *model.CategoryDto) (int, error) {
	var parentID *int
	if c.ParentCategory != nil && c.ParentCategory.CategoryId > 0 {
		parentID = &c.ParentCategory.CategoryId
	}
	var id int
	query := `INSERT INTO categories (category_title, image_url, parent_category_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING category_id`
	err := r.db.QueryRow(ctx, query, c.CategoryTitle, c.ImageUrl, parentID).Scan(&id)
	return id, err
}

func (r *ProductRepository) UpdateCategory(ctx context.Context, id int, c *model.CategoryDto) error {
	var parentID *int
	if c.ParentCategory != nil && c.ParentCategory.CategoryId > 0 {
		parentID = &c.ParentCategory.CategoryId
	}
	query := `UPDATE categories SET category_title = $1, image_url = $2, parent_category_id = $3, updated_at = NOW() WHERE category_id = $4`
	_, err := r.db.Exec(ctx, query, c.CategoryTitle, c.ImageUrl, parentID, id)
	return err
}

func (r *ProductRepository) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, id)
	return err
}

// =================== Products ===================

func (r *ProductRepository) FindAllProducts(ctx context.Context) ([]model.ProductDto, error) {
	query := `
		SELECT p.product_id, p.product_title, COALESCE(p.image_url, ''), COALESCE(p.sku, ''),
		       COALESCE(p.price_unit, 0), COALESCE(p.quantity, 0),
		       COALESCE(c.category_id, 0), COALESCE(c.category_title, ''), COALESCE(c.image_url, '')
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.category_id
		ORDER BY p.product_id DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ProductDto
	for rows.Next() {
		var p model.ProductDto
		var catID int
		var catTitle, catImage string
		if err := rows.Scan(&p.ProductId, &p.ProductTitle, &p.ImageUrl, &p.Sku, &p.PriceUnit, &p.Quantity, &catID, &catTitle, &catImage); err != nil {
			return nil, err
		}
		if catID > 0 {
			p.Category = &model.CategoryDto{
				CategoryId:    catID,
				CategoryTitle: catTitle,
				ImageUrl:      catImage,
			}
		}
		list = append(list, p)
	}
	if list == nil {
		list = []model.ProductDto{}
	}
	return list, nil
}

func (r *ProductRepository) FindProductByID(ctx context.Context, id int) (*model.ProductDto, error) {
	query := `
		SELECT p.product_id, p.product_title, COALESCE(p.image_url, ''), COALESCE(p.sku, ''),
		       COALESCE(p.price_unit, 0), COALESCE(p.quantity, 0),
		       COALESCE(c.category_id, 0), COALESCE(c.category_title, ''), COALESCE(c.image_url, '')
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.category_id
		WHERE p.product_id = $1`
	var p model.ProductDto
	var catID int
	var catTitle, catImage string
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ProductId, &p.ProductTitle, &p.ImageUrl, &p.Sku, &p.PriceUnit, &p.Quantity, &catID, &catTitle, &catImage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if catID > 0 {
		p.Category = &model.CategoryDto{
			CategoryId:    catID,
			CategoryTitle: catTitle,
			ImageUrl:      catImage,
		}
	}
	return &p, nil
}

func (r *ProductRepository) SaveProduct(ctx context.Context, p *model.ProductDto) (int, error) {
	var catID *int
	if p.Category != nil && p.Category.CategoryId > 0 {
		catID = &p.Category.CategoryId
	}
	var id int
	query := `
		INSERT INTO products (product_title, image_url, sku, price_unit, quantity, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING product_id`
	err := r.db.QueryRow(ctx, query, p.ProductTitle, p.ImageUrl, p.Sku, p.PriceUnit, p.Quantity, catID).Scan(&id)
	return id, err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, p *model.ProductDto) error {
	var catID *int
	if p.Category != nil && p.Category.CategoryId > 0 {
		catID = &p.Category.CategoryId
	}
	query := `
		UPDATE products
		SET product_title = $1, image_url = $2, sku = $3, price_unit = $4, quantity = $5, category_id = $6, updated_at = NOW()
		WHERE product_id = $7`
	_, err := r.db.Exec(ctx, query, p.ProductTitle, p.ImageUrl, p.Sku, p.PriceUnit, p.Quantity, catID, id)
	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE product_id = $1`, id)
	return err
}

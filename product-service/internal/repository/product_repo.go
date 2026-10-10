package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"com.ecommerce/product-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/fm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type ProductRepository struct {
	db pgxdriver.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: pgxdriver.NewPool(pool)}
}

func (r *ProductRepository) InitSchema(ctx context.Context) {
	if r.db.Pool == nil {
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

		CREATE TABLE IF NOT EXISTS favourites (
			user_id INT NOT NULL,
			product_id INT NOT NULL,
			like_date TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, product_id, like_date)
		);
	`
	if _, err := r.db.Exec(ctx, query); err != nil {
		log.Printf("[Product Service] Notice during schema verification: %v", err)
	}
}

// =================== Categories ===================

func categoryColumns() []any {
	return []any{
		psql.Quote("category_id"),
		psql.Quote("category_title"),
		psql.F("COALESCE", psql.Quote("image_url"), psql.S(""))(fm.As("image_url")),
	}
}

func (r *ProductRepository) FindAllCategories(ctx context.Context) ([]model.CategoryDto, error) {
	if r.db.Pool == nil {
		return []model.CategoryDto{}, nil
	}
	query := psql.Select(
		sm.Columns(categoryColumns()...),
		sm.From("categories"),
		sm.OrderBy(psql.Quote("category_id")).Asc(),
	)
	list, err := bob.All(ctx, r.db, query, scan.StructMapper[model.CategoryDto]())
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.CategoryDto{}
	}
	return list, nil
}

func (r *ProductRepository) FindCategoryByID(ctx context.Context, id int) (*model.CategoryDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns(categoryColumns()...),
		sm.From("categories"),
		sm.Where(psql.Quote("category_id").EQ(psql.Arg(id))),
	)
	c, err := bob.One(ctx, r.db, query, scan.StructMapper[model.CategoryDto]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ProductRepository) SaveCategory(ctx context.Context, c *model.CategoryDto) (int, error) {
	if r.db.Pool == nil {
		return 1, nil
	}
	var parentID *int
	if c.ParentCategory != nil && c.ParentCategory.CategoryId > 0 {
		parentID = &c.ParentCategory.CategoryId
	}
	query := psql.Insert(
		im.Into("categories", "category_title", "image_url", "parent_category_id", "created_at", "updated_at"),
		im.Values(
			psql.Arg(c.CategoryTitle),
			psql.Arg(c.ImageUrl),
			psql.Arg(parentID),
			psql.Raw("NOW()"),
			psql.Raw("NOW()"),
		),
		im.Returning("category_id"),
	)
	return bob.One(ctx, r.db, query, scan.SingleColumnMapper[int])
}

func (r *ProductRepository) UpdateCategory(ctx context.Context, id int, c *model.CategoryDto) error {
	if r.db.Pool == nil {
		return nil
	}
	var parentID *int
	if c.ParentCategory != nil && c.ParentCategory.CategoryId > 0 {
		parentID = &c.ParentCategory.CategoryId
	}
	query := psql.Update(
		um.Table("categories"),
		um.Set(
			psql.Quote("category_title").EQ(psql.Arg(c.CategoryTitle)),
			psql.Quote("image_url").EQ(psql.Arg(c.ImageUrl)),
			psql.Quote("parent_category_id").EQ(psql.Arg(parentID)),
			psql.Quote("updated_at").EQ(psql.Raw("NOW()")),
		),
		um.Where(psql.Quote("category_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

func (r *ProductRepository) DeleteCategory(ctx context.Context, id int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("categories"),
		dm.Where(psql.Quote("category_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

// =================== Products ===================

// productRow mirrors the joined products/categories projection used for scanning.
type productRow struct {
	ProductId        int     `db:"product_id"`
	ProductTitle     string  `db:"product_title"`
	ImageUrl         string  `db:"image_url"`
	Sku              string  `db:"sku"`
	PriceUnit        float64 `db:"price_unit"`
	Quantity         int     `db:"quantity"`
	CategoryId       int     `db:"category_id"`
	CategoryTitle    string  `db:"category_title"`
	CategoryImageUrl string  `db:"category_image_url"`
}

func (row productRow) toDto() model.ProductDto {
	p := model.ProductDto{
		ProductId:    row.ProductId,
		ProductTitle: row.ProductTitle,
		ImageUrl:     row.ImageUrl,
		Sku:          row.Sku,
		PriceUnit:    row.PriceUnit,
		Quantity:     row.Quantity,
	}
	if row.CategoryId > 0 {
		p.Category = &model.CategoryDto{
			CategoryId:    row.CategoryId,
			CategoryTitle: row.CategoryTitle,
			ImageUrl:      row.CategoryImageUrl,
		}
	}
	return p
}

func productColumns() []any {
	return []any{
		psql.Quote("p", "product_id"),
		psql.Quote("p", "product_title"),
		psql.F("COALESCE", psql.Quote("p", "image_url"), psql.S(""))(fm.As("image_url")),
		psql.F("COALESCE", psql.Quote("p", "sku"), psql.S(""))(fm.As("sku")),
		psql.F("COALESCE", psql.Quote("p", "price_unit"), psql.Arg(0))(fm.As("price_unit")),
		psql.F("COALESCE", psql.Quote("p", "quantity"), psql.Arg(0))(fm.As("quantity")),
		psql.F("COALESCE", psql.Quote("c", "category_id"), psql.Arg(0))(fm.As("category_id")),
		psql.F("COALESCE", psql.Quote("c", "category_title"), psql.S(""))(fm.As("category_title")),
		psql.F("COALESCE", psql.Quote("c", "image_url"), psql.S(""))(fm.As("category_image_url")),
	}
}

func (r *ProductRepository) FindAllProducts(ctx context.Context) ([]model.ProductDto, error) {
	if r.db.Pool == nil {
		return []model.ProductDto{}, nil
	}
	query := psql.Select(
		sm.Columns(productColumns()...),
		sm.From("products p", sm.LeftJoin("categories c").On(psql.Raw("c.category_id = p.category_id"))),
		sm.OrderBy(psql.Quote("p", "product_id")).Desc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[productRow]())
	if err != nil {
		return nil, err
	}
	list := make([]model.ProductDto, 0, len(rows))
	for _, row := range rows {
		list = append(list, row.toDto())
	}
	return list, nil
}

func (r *ProductRepository) FindProductByID(ctx context.Context, id int) (*model.ProductDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns(productColumns()...),
		sm.From("products p", sm.LeftJoin("categories c").On(psql.Raw("c.category_id = p.category_id"))),
		sm.Where(psql.Quote("p", "product_id").EQ(psql.Arg(id))),
	)
	row, err := bob.One(ctx, r.db, query, scan.StructMapper[productRow]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	p := row.toDto()
	return &p, nil
}

func (r *ProductRepository) SaveProduct(ctx context.Context, p *model.ProductDto) (int, error) {
	if r.db.Pool == nil {
		return 1, nil
	}
	var catID *int
	if p.Category != nil && p.Category.CategoryId > 0 {
		catID = &p.Category.CategoryId
	}
	query := psql.Insert(
		im.Into("products", "product_title", "image_url", "sku", "price_unit", "quantity", "category_id", "created_at", "updated_at"),
		im.Values(
			psql.Arg(p.ProductTitle),
			psql.Arg(p.ImageUrl),
			psql.Arg(p.Sku),
			psql.Arg(p.PriceUnit),
			psql.Arg(p.Quantity),
			psql.Arg(catID),
			psql.Raw("NOW()"),
			psql.Raw("NOW()"),
		),
		im.Returning("product_id"),
	)
	return bob.One(ctx, r.db, query, scan.SingleColumnMapper[int])
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, p *model.ProductDto) error {
	if r.db.Pool == nil {
		return nil
	}
	var catID *int
	if p.Category != nil && p.Category.CategoryId > 0 {
		catID = &p.Category.CategoryId
	}
	query := psql.Update(
		um.Table("products"),
		um.Set(
			psql.Quote("product_title").EQ(psql.Arg(p.ProductTitle)),
			psql.Quote("image_url").EQ(psql.Arg(p.ImageUrl)),
			psql.Quote("sku").EQ(psql.Arg(p.Sku)),
			psql.Quote("price_unit").EQ(psql.Arg(p.PriceUnit)),
			psql.Quote("quantity").EQ(psql.Arg(p.Quantity)),
			psql.Quote("category_id").EQ(psql.Arg(catID)),
			psql.Quote("updated_at").EQ(psql.Raw("NOW()")),
		),
		um.Where(psql.Quote("product_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("products"),
		dm.Where(psql.Quote("product_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

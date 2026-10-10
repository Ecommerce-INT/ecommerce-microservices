package repository

import (
	"context"
	"errors"
	"log"

	"com.ecommerce/shipping-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShippingRepository struct {
	db *pgxpool.Pool
}

func NewShippingRepository(db *pgxpool.Pool) *ShippingRepository {
	return &ShippingRepository{db: db}
}

func (r *ShippingRepository) InitSchema(ctx context.Context) {
	if r.db == nil {
		return
	}
	query := `
		CREATE TABLE IF NOT EXISTS order_items (
			order_id INT NOT NULL,
			product_id INT NOT NULL,
			ordered_quantity INT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (order_id, product_id)
		);
	`
	if _, err := r.db.Exec(ctx, query); err != nil {
		log.Printf("[Shipping Service] Notice during schema verification: %v", err)
	}
}

func (r *ShippingRepository) FindAll(ctx context.Context) ([]model.OrderItemDto, error) {
	if r.db == nil {
		return []model.OrderItemDto{}, nil
	}
	query := `SELECT order_id, product_id, ordered_quantity, created_at, updated_at FROM order_items ORDER BY order_id DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.OrderItemDto
	for rows.Next() {
		var item model.OrderItemDto
		if err := rows.Scan(&item.OrderId, &item.ProductId, &item.OrderedQuantity, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.OrderItemDto{}
	}
	return list, nil
}

func (r *ShippingRepository) FindByID(ctx context.Context, orderId, productId int) (*model.OrderItemDto, error) {
	if r.db == nil {
		return nil, nil
	}
	query := `SELECT order_id, product_id, ordered_quantity, created_at, updated_at FROM order_items WHERE order_id = $1 AND product_id = $2`
	var item model.OrderItemDto
	err := r.db.QueryRow(ctx, query, orderId, productId).Scan(&item.OrderId, &item.ProductId, &item.OrderedQuantity, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ShippingRepository) Save(ctx context.Context, item *model.OrderItemDto) error {
	if r.db == nil {
		return nil
	}
	query := `
		INSERT INTO order_items (order_id, product_id, ordered_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (order_id, product_id) DO UPDATE
		SET ordered_quantity = EXCLUDED.ordered_quantity, updated_at = NOW()`
	_, err := r.db.Exec(ctx, query, item.OrderId, item.ProductId, item.OrderedQuantity)
	return err
}

func (r *ShippingRepository) DeleteByID(ctx context.Context, orderId, productId int) error {
	if r.db == nil {
		return nil
	}
	query := `DELETE FROM order_items WHERE order_id = $1 AND product_id = $2`
	_, err := r.db.Exec(ctx, query, orderId, productId)
	return err
}

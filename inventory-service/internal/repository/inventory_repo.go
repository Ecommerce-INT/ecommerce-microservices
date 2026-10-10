package repository

import (
	"context"
	"log"

	"com.ecommerce/inventory-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) InitSchema(ctx context.Context) {
	if r.db == nil {
		return
	}
	query := `
		CREATE TABLE IF NOT EXISTS inventory (
			id BIGSERIAL PRIMARY KEY,
			product_name VARCHAR(255) NOT NULL,
			quantity INT DEFAULT 0
		);
	`
	if _, err := r.db.Exec(ctx, query); err != nil {
		log.Printf("[Inventory Service] Notice during schema verification: %v", err)
	}
}

func (r *InventoryRepository) FindByProductNames(ctx context.Context, productNames []string) ([]model.Inventory, error) {
	if r.db == nil || len(productNames) == 0 {
		return []model.Inventory{}, nil
	}

	query := `SELECT id, product_name, quantity FROM inventory WHERE product_name = ANY($1)`
	rows, err := r.db.Query(ctx, query, productNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Inventory
	for rows.Next() {
		var item model.Inventory
		if err := rows.Scan(&item.ID, &item.ProductName, &item.Quantity); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.Inventory{}
	}
	return list, nil
}

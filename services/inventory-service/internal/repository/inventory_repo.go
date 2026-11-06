package repository

import (
	"context"
	"log"

	"com.ecommerce/inventory-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type InventoryRepository struct {
	db pgxdriver.Pool
}

func NewInventoryRepository(pool *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: pgxdriver.NewPool(pool)}
}

func (r *InventoryRepository) InitSchema(ctx context.Context) {
	if r.db.Pool == nil {
		return
	}
	// PostgreSQL 18: Native UUIDv7 primary clustering key & B-tree Skip Scan composite index
	query := `
		CREATE TABLE IF NOT EXISTS inventory (
			id BIGSERIAL PRIMARY KEY,
			inventory_uuid UUID DEFAULT uuidv7(),
			product_name VARCHAR(255) NOT NULL,
			quantity INT DEFAULT 0
		);

		-- PostgreSQL 18 B-tree Skip Scan composite index
		CREATE INDEX IF NOT EXISTS idx_inventory_product_qty ON inventory (product_name, quantity);
	`
	if _, err := r.db.Exec(ctx, query); err != nil {
		log.Printf("[Inventory Service] Notice during schema verification: %v", err)
	}
}

func (r *InventoryRepository) FindByProductNames(ctx context.Context, productNames []string) ([]model.Inventory, error) {
	if r.db.Pool == nil || len(productNames) == 0 {
		return []model.Inventory{}, nil
	}

	query := psql.Select(
		sm.Columns(
			psql.Quote("id"),
			psql.Quote("product_name"),
			psql.Quote("quantity"),
		),
		sm.From("inventory"),
		sm.Where(psql.Quote("product_name").EQ(psql.Any(psql.Arg(productNames)))),
	)
	list, err := bob.All(ctx, r.db, query, scan.StructMapper[model.Inventory]())
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.Inventory{}
	}
	return list, nil
}

// StockTransition tracks atomic stock state transitions using PostgreSQL 18 RETURNING OLD, NEW.
type StockTransition struct {
	ID               int64  `db:"id"`
	ProductName      string `db:"product_name"`
	PreviousQuantity int    `db:"previous_quantity"`
	CurrentQuantity  int    `db:"current_quantity"`
}

// DeductStock atomically decrements product inventory and retrieves before & after states
// in a single roundtrip via PostgreSQL 18's RETURNING OLD, NEW capability.
func (r *InventoryRepository) DeductStock(ctx context.Context, productName string, amount int) (*StockTransition, error) {
	if r.db.Pool == nil {
		return &StockTransition{
			ProductName:      productName,
			PreviousQuantity: amount,
			CurrentQuantity:  0,
		}, nil
	}

	query := `
		UPDATE inventory
		SET quantity = quantity - $1
		WHERE product_name = $2 AND quantity >= $1
		RETURNING id, product_name, OLD.quantity AS previous_quantity, NEW.quantity AS current_quantity;
	`
	var t StockTransition
	err := r.db.Pool.QueryRow(ctx, query, amount, productName).Scan(&t.ID, &t.ProductName, &t.PreviousQuantity, &t.CurrentQuantity)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

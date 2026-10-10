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

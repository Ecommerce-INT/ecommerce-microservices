package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"com.ecommerce/shipping-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type ShippingRepository struct {
	db pgxdriver.Pool
}

func NewShippingRepository(pool *pgxpool.Pool) *ShippingRepository {
	return &ShippingRepository{db: pgxdriver.NewPool(pool)}
}

func (r *ShippingRepository) InitSchema(ctx context.Context) {
	if r.db.Pool == nil {
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

func orderItemColumns() []any {
	return []any{
		psql.Quote("order_id"),
		psql.Quote("product_id"),
		psql.Quote("ordered_quantity"),
		psql.Quote("created_at"),
		psql.Quote("updated_at"),
	}
}

func (r *ShippingRepository) FindAll(ctx context.Context) ([]model.OrderItemDto, error) {
	if r.db.Pool == nil {
		return []model.OrderItemDto{}, nil
	}
	query := psql.Select(
		sm.Columns(orderItemColumns()...),
		sm.From("order_items"),
		sm.OrderBy(psql.Quote("order_id")).Desc(),
	)
	list, err := bob.All(ctx, r.db, query, scan.StructMapper[model.OrderItemDto]())
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.OrderItemDto{}
	}
	return list, nil
}

func (r *ShippingRepository) FindByID(ctx context.Context, orderId, productId int) (*model.OrderItemDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns(orderItemColumns()...),
		sm.From("order_items"),
		sm.Where(psql.And(
			psql.Quote("order_id").EQ(psql.Arg(orderId)),
			psql.Quote("product_id").EQ(psql.Arg(productId)),
		)),
	)
	item, err := bob.One(ctx, r.db, query, scan.StructMapper[model.OrderItemDto]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ShippingRepository) Save(ctx context.Context, item *model.OrderItemDto) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Insert(
		im.Into("order_items", "order_id", "product_id", "ordered_quantity", "created_at", "updated_at"),
		im.Values(
			psql.Arg(item.OrderId),
			psql.Arg(item.ProductId),
			psql.Arg(item.OrderedQuantity),
			psql.Raw("NOW()"),
			psql.Raw("NOW()"),
		),
		im.OnConflict("order_id", "product_id").
			DoUpdate(im.Set(
				psql.Quote("ordered_quantity").EQ(psql.Raw("EXCLUDED.ordered_quantity")),
				psql.Quote("updated_at").EQ(psql.Raw("NOW()")),
			)),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

func (r *ShippingRepository) DeleteByID(ctx context.Context, orderId, productId int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("order_items"),
		dm.Where(psql.And(
			psql.Quote("order_id").EQ(psql.Arg(orderId)),
			psql.Quote("product_id").EQ(psql.Arg(productId)),
		)),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

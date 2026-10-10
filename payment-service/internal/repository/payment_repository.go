package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"com.ecommerce/payment-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type PaymentRepository struct {
	db pgxdriver.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: pgxdriver.NewPool(pool)}
}

func (r *PaymentRepository) InitSchema(ctx context.Context) error {
	if r.db.Pool == nil {
		return nil
	}
	schema := `
	CREATE TABLE IF NOT EXISTS payments (
		payment_id SERIAL PRIMARY KEY,
		order_id INT,
		user_id BIGINT,
		is_payed BOOLEAN DEFAULT FALSE,
		payment_status VARCHAR(50) DEFAULT 'NOT_STARTED',
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(ctx, schema)
	return err
}

type paymentRow struct {
	PaymentId     int    `db:"payment_id"`
	OrderId       *int   `db:"order_id"`
	UserId        *int64 `db:"user_id"`
	IsPayed       *bool  `db:"is_payed"`
	PaymentStatus string `db:"payment_status"`
}

func (row paymentRow) toDto() model.PaymentDto {
	paymentId := row.PaymentId
	p := model.PaymentDto{
		PaymentId:     &paymentId,
		OrderId:       row.OrderId,
		UserId:        row.UserId,
		IsPayed:       row.IsPayed,
		PaymentStatus: row.PaymentStatus,
	}
	if row.OrderId != nil {
		p.Order = &model.OrderDto{OrderId: row.OrderId}
	}
	if row.UserId != nil {
		p.User = &model.UserDto{ID: row.UserId}
	}
	return p
}

func paymentColumns() []any {
	return []any{
		psql.Quote("payment_id"),
		psql.Quote("order_id"),
		psql.Quote("user_id"),
		psql.Quote("is_payed"),
		psql.Quote("payment_status"),
	}
}

func orderByMod(col string, order string) dialect.OrderBy[*dialect.SelectQuery] {
	mod := sm.OrderBy(psql.Quote(col))
	if order == "DESC" {
		return mod.Desc()
	}
	return mod.Asc()
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]model.PaymentDto, error) {
	if r.db.Pool == nil {
		return []model.PaymentDto{}, nil
	}

	query := psql.Select(
		sm.Columns(paymentColumns()...),
		sm.From("payments"),
		sm.OrderBy(psql.Quote("payment_id")).Asc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[paymentRow]())
	if err != nil {
		return nil, err
	}

	payments := make([]model.PaymentDto, 0, len(rows))
	for _, row := range rows {
		payments = append(payments, row.toDto())
	}
	return payments, nil
}

func (r *PaymentRepository) FindAllPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.PaymentDto, int64, error) {
	if r.db.Pool == nil {
		return []model.PaymentDto{}, 0, nil
	}

	col := "payment_id"
	switch strings.ToLower(sortBy) {
	case "orderid", "order_id":
		col = "order_id"
	case "userid", "user_id":
		col = "user_id"
	case "ispayed", "is_payed":
		col = "is_payed"
	default:
		col = "payment_id"
	}

	order := "ASC"
	if strings.ToUpper(sortOrder) == "DESC" {
		order = "DESC"
	}

	countQuery := psql.Select(
		sm.Columns("COUNT(*)"),
		sm.From("payments"),
	)
	total, err := bob.One(ctx, r.db, countQuery, scan.SingleColumnMapper[int64])
	if err != nil {
		return nil, 0, err
	}

	query := psql.Select(
		sm.Columns(paymentColumns()...),
		sm.From("payments"),
		orderByMod(col, order),
		sm.Limit(size),
		sm.Offset(page*size),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[paymentRow]())
	if err != nil {
		return nil, 0, err
	}

	payments := make([]model.PaymentDto, 0, len(rows))
	for _, row := range rows {
		payments = append(payments, row.toDto())
	}
	return payments, total, nil
}

func (r *PaymentRepository) FindById(ctx context.Context, id int) (*model.PaymentDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns(paymentColumns()...),
		sm.From("payments"),
		sm.Where(psql.Quote("payment_id").EQ(psql.Arg(id))),
	)
	row, err := bob.One(ctx, r.db, query, scan.StructMapper[paymentRow]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	payment := row.toDto()
	return &payment, nil
}

func (r *PaymentRepository) ExistsByOrderIdAndIsPayed(ctx context.Context, orderId *int) (bool, error) {
	if r.db.Pool == nil || orderId == nil {
		return false, nil
	}
	query := psql.Select(
		sm.Columns(psql.Exists(psql.Select(
			sm.Columns("1"),
			sm.From("payments"),
			sm.Where(psql.And(
				psql.Quote("order_id").EQ(psql.Arg(*orderId)),
				psql.Quote("is_payed").EQ(psql.Raw("TRUE")),
			)),
		))),
	)
	return bob.One(ctx, r.db, query, scan.SingleColumnMapper[bool])
}

func (r *PaymentRepository) Save(ctx context.Context, p model.PaymentDto) (*model.PaymentDto, error) {
	if r.db.Pool == nil {
		id := 1
		p.PaymentId = &id
		return &p, nil
	}

	query := psql.Insert(
		im.Into("payments", "order_id", "user_id", "is_payed", "payment_status", "created_at", "updated_at"),
		im.Values(
			psql.Arg(p.OrderId),
			psql.Arg(p.UserId),
			psql.F("COALESCE", psql.Arg(p.IsPayed), psql.Raw("FALSE")),
			psql.F("COALESCE", psql.Arg(p.PaymentStatus), psql.S("NOT_STARTED")),
			psql.Raw("CURRENT_TIMESTAMP"),
			psql.Raw("CURRENT_TIMESTAMP"),
		),
		im.Returning("payment_id"),
	)
	id, err := bob.One(ctx, r.db, query, scan.SingleColumnMapper[int])
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *PaymentRepository) Update(ctx context.Context, id int, p model.PaymentDto) (*model.PaymentDto, error) {
	if r.db.Pool == nil {
		p.PaymentId = &id
		return &p, nil
	}

	query := psql.Update(
		um.Table("payments"),
		um.Set(
			psql.Quote("order_id").EQ(psql.F("COALESCE", psql.Arg(p.OrderId), psql.Quote("order_id"))),
			psql.Quote("user_id").EQ(psql.F("COALESCE", psql.Arg(p.UserId), psql.Quote("user_id"))),
			psql.Quote("is_payed").EQ(psql.F("COALESCE", psql.Arg(p.IsPayed), psql.Quote("is_payed"))),
			psql.Quote("payment_status").EQ(psql.F("COALESCE", psql.Arg(p.PaymentStatus), psql.Quote("payment_status"))),
			psql.Quote("updated_at").EQ(psql.Raw("CURRENT_TIMESTAMP")),
		),
		um.Where(psql.Quote("payment_id").EQ(psql.Arg(id))),
	)
	if _, err := bob.Exec(ctx, r.db, query); err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *PaymentRepository) DeleteById(ctx context.Context, id int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("payments"),
		dm.Where(psql.Quote("payment_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

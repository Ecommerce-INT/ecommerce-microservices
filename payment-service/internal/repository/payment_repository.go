package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"com.ecommerce/payment-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) InitSchema(ctx context.Context) error {
	if r.pool == nil {
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
	_, err := r.pool.Exec(ctx, schema)
	return err
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]model.PaymentDto, error) {
	if r.pool == nil {
		return []model.PaymentDto{}, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT payment_id, order_id, user_id, is_payed, payment_status
		FROM payments
		ORDER BY payment_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.PaymentDto
	for rows.Next() {
		var p model.PaymentDto
		if err := rows.Scan(&p.PaymentId, &p.OrderId, &p.UserId, &p.IsPayed, &p.PaymentStatus); err != nil {
			return nil, err
		}
		if p.OrderId != nil {
			p.Order = &model.OrderDto{OrderId: p.OrderId}
		}
		if p.UserId != nil {
			p.User = &model.UserDto{ID: p.UserId}
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PaymentRepository) FindAllPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.PaymentDto, int64, error) {
	if r.pool == nil {
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

	var total int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM payments").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := page * size
	query := fmt.Sprintf(`
		SELECT payment_id, order_id, user_id, is_payed, payment_status
		FROM payments
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, col, order)

	rows, err := r.pool.Query(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []model.PaymentDto
	for rows.Next() {
		var p model.PaymentDto
		if err := rows.Scan(&p.PaymentId, &p.OrderId, &p.UserId, &p.IsPayed, &p.PaymentStatus); err != nil {
			return nil, 0, err
		}
		if p.OrderId != nil {
			p.Order = &model.OrderDto{OrderId: p.OrderId}
		}
		if p.UserId != nil {
			p.User = &model.UserDto{ID: p.UserId}
		}
		payments = append(payments, p)
	}
	return payments, total, nil
}

func (r *PaymentRepository) FindById(ctx context.Context, id int) (*model.PaymentDto, error) {
	if r.pool == nil {
		return nil, nil
	}
	var p model.PaymentDto
	err := r.pool.QueryRow(ctx, `
		SELECT payment_id, order_id, user_id, is_payed, payment_status
		FROM payments
		WHERE payment_id = $1
	`, id).Scan(&p.PaymentId, &p.OrderId, &p.UserId, &p.IsPayed, &p.PaymentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if p.OrderId != nil {
		p.Order = &model.OrderDto{OrderId: p.OrderId}
	}
	if p.UserId != nil {
		p.User = &model.UserDto{ID: p.UserId}
	}
	return &p, nil
}

func (r *PaymentRepository) ExistsByOrderIdAndIsPayed(ctx context.Context, orderId *int) (bool, error) {
	if r.pool == nil || orderId == nil {
		return false, nil
	}
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM payments WHERE order_id = $1 AND is_payed = TRUE)
	`, *orderId).Scan(&exists)
	return exists, err
}

func (r *PaymentRepository) Save(ctx context.Context, p model.PaymentDto) (*model.PaymentDto, error) {
	if r.pool == nil {
		id := 1
		p.PaymentId = &id
		return &p, nil
	}

	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO payments (order_id, user_id, is_payed, payment_status, created_at, updated_at)
		VALUES ($1, $2, COALESCE($3, FALSE), COALESCE($4, 'NOT_STARTED'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING payment_id
	`, p.OrderId, p.UserId, p.IsPayed, p.PaymentStatus).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *PaymentRepository) Update(ctx context.Context, id int, p model.PaymentDto) (*model.PaymentDto, error) {
	if r.pool == nil {
		p.PaymentId = &id
		return &p, nil
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE payments
		SET order_id = COALESCE($1, order_id),
		    user_id = COALESCE($2, user_id),
		    is_payed = COALESCE($3, is_payed),
		    payment_status = COALESCE($4, payment_status),
		    updated_at = CURRENT_TIMESTAMP
		WHERE payment_id = $5
	`, p.OrderId, p.UserId, p.IsPayed, p.PaymentStatus, id)
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *PaymentRepository) DeleteById(ctx context.Context, id int) error {
	if r.pool == nil {
		return nil
	}
	_, err := r.pool.Exec(ctx, `DELETE FROM payments WHERE payment_id = $1`, id)
	return err
}

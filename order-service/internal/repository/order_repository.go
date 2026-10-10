package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"com.ecommerce/order-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) InitSchema(ctx context.Context) error {
	if r.pool == nil {
		return nil
	}
	schema := `
	CREATE TABLE IF NOT EXISTS carts (
		cart_id SERIAL PRIMARY KEY,
		user_id BIGINT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		order_id SERIAL PRIMARY KEY,
		order_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		order_desc TEXT,
		order_fee DOUBLE PRECISION,
		product_id INT,
		cart_id INT REFERENCES carts(cart_id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.pool.Exec(ctx, schema)
	return err
}

// -------------------------------------------------------------
// Cart Methods
// -------------------------------------------------------------

func (r *OrderRepository) FindAllCarts(ctx context.Context) ([]model.CartDto, error) {
	if r.pool == nil {
		return []model.CartDto{}, nil
	}
	rows, err := r.pool.Query(ctx, `SELECT cart_id, user_id FROM carts ORDER BY cart_id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []model.CartDto
	for rows.Next() {
		var c model.CartDto
		if err := rows.Scan(&c.CartId, &c.UserId); err != nil {
			return nil, err
		}
		if c.UserId != nil {
			c.User = &model.UserDto{ID: c.UserId}
		}
		c.Orders, _ = r.findOrdersByCartId(ctx, *c.CartId)
		carts = append(carts, c)
	}
	return carts, nil
}

func (r *OrderRepository) FindAllCartsPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.CartDto, int64, error) {
	if r.pool == nil {
		return []model.CartDto{}, 0, nil
	}

	col := "cart_id"
	if strings.ToLower(sortBy) == "userid" || strings.ToLower(sortBy) == "user_id" {
		col = "user_id"
	}

	order := "ASC"
	if strings.ToUpper(sortOrder) == "DESC" {
		order = "DESC"
	}

	var total int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM carts").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := page * size
	query := fmt.Sprintf(`SELECT cart_id, user_id FROM carts ORDER BY %s %s LIMIT $1 OFFSET $2`, col, order)
	rows, err := r.pool.Query(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var carts []model.CartDto
	for rows.Next() {
		var c model.CartDto
		if err := rows.Scan(&c.CartId, &c.UserId); err != nil {
			return nil, 0, err
		}
		if c.UserId != nil {
			c.User = &model.UserDto{ID: c.UserId}
		}
		c.Orders, _ = r.findOrdersByCartId(ctx, *c.CartId)
		carts = append(carts, c)
	}
	return carts, total, nil
}

func (r *OrderRepository) FindCartById(ctx context.Context, id int) (*model.CartDto, error) {
	if r.pool == nil {
		return nil, nil
	}
	var c model.CartDto
	err := r.pool.QueryRow(ctx, `SELECT cart_id, user_id FROM carts WHERE cart_id = $1`, id).Scan(&c.CartId, &c.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if c.UserId != nil {
		c.User = &model.UserDto{ID: c.UserId}
	}
	c.Orders, _ = r.findOrdersByCartId(ctx, *c.CartId)
	return &c, nil
}

func (r *OrderRepository) SaveCart(ctx context.Context, cart model.CartDto) (*model.CartDto, error) {
	if r.pool == nil {
		id := 1
		cart.CartId = &id
		return &cart, nil
	}

	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO carts (user_id, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING cart_id
	`, cart.UserId).Scan(&id)
	if err != nil {
		return nil, err
	}
	cart.CartId = &id
	if cart.UserId != nil {
		cart.User = &model.UserDto{ID: cart.UserId}
	}

	if len(cart.Orders) > 0 {
		for _, o := range cart.Orders {
			o.Cart = &model.CartDto{CartId: &id, UserId: cart.UserId}
			_, _ = r.SaveOrder(ctx, o)
		}
		cart.Orders, _ = r.findOrdersByCartId(ctx, id)
	}

	return &cart, nil
}

func (r *OrderRepository) UpdateCart(ctx context.Context, id int, cart model.CartDto) (*model.CartDto, error) {
	if r.pool == nil {
		cart.CartId = &id
		return &cart, nil
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE carts
		SET user_id = COALESCE($1, user_id), updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $2
	`, cart.UserId, id)
	if err != nil {
		return nil, err
	}
	return r.FindCartById(ctx, id)
}

func (r *OrderRepository) DeleteCartById(ctx context.Context, id int) error {
	if r.pool == nil {
		return nil
	}
	_, err := r.pool.Exec(ctx, `DELETE FROM carts WHERE cart_id = $1`, id)
	return err
}

func (r *OrderRepository) findOrdersByCartId(ctx context.Context, cartId int) ([]model.OrderDto, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT order_id, to_char(order_date, 'YYYY-MM-DD"T"HH24:MI:SS'), order_desc, order_fee, product_id
		FROM orders
		WHERE cart_id = $1
		ORDER BY order_id ASC
	`, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderDto
	for rows.Next() {
		var o model.OrderDto
		var desc *string
		if err := rows.Scan(&o.OrderId, &o.OrderDate, &desc, &o.OrderFee, &o.ProductId); err != nil {
			return nil, err
		}
		if desc != nil {
			o.OrderDesc = *desc
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// -------------------------------------------------------------
// Order Methods
// -------------------------------------------------------------

func (r *OrderRepository) FindAllOrders(ctx context.Context) ([]model.OrderDto, error) {
	if r.pool == nil {
		return []model.OrderDto{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT o.order_id, to_char(o.order_date, 'YYYY-MM-DD"T"HH24:MI:SS'), o.order_desc, o.order_fee, o.product_id,
		       c.cart_id, c.user_id
		FROM orders o
		LEFT JOIN carts c ON c.cart_id = o.cart_id
		ORDER BY o.order_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderDto
	for rows.Next() {
		var o model.OrderDto
		var desc *string
		var cartId *int
		var userId *int64
		if err := rows.Scan(&o.OrderId, &o.OrderDate, &desc, &o.OrderFee, &o.ProductId, &cartId, &userId); err != nil {
			return nil, err
		}
		if desc != nil {
			o.OrderDesc = *desc
		}
		if cartId != nil {
			o.Cart = &model.CartDto{CartId: cartId, UserId: userId}
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) FindAllOrdersPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.OrderDto, int64, error) {
	if r.pool == nil {
		return []model.OrderDto{}, 0, nil
	}

	col := "order_id"
	switch strings.ToLower(sortBy) {
	case "orderfee", "order_fee":
		col = "order_fee"
	case "orderdate", "order_date":
		col = "order_date"
	default:
		col = "order_id"
	}

	order := "ASC"
	if strings.ToUpper(sortOrder) == "DESC" {
		order = "DESC"
	}

	var total int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := page * size
	query := fmt.Sprintf(`
		SELECT o.order_id, to_char(o.order_date, 'YYYY-MM-DD"T"HH24:MI:SS'), o.order_desc, o.order_fee, o.product_id,
		       c.cart_id, c.user_id
		FROM orders o
		LEFT JOIN carts c ON c.cart_id = o.cart_id
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, col, order)

	rows, err := r.pool.Query(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []model.OrderDto
	for rows.Next() {
		var o model.OrderDto
		var desc *string
		var cartId *int
		var userId *int64
		if err := rows.Scan(&o.OrderId, &o.OrderDate, &desc, &o.OrderFee, &o.ProductId, &cartId, &userId); err != nil {
			return nil, 0, err
		}
		if desc != nil {
			o.OrderDesc = *desc
		}
		if cartId != nil {
			o.Cart = &model.CartDto{CartId: cartId, UserId: userId}
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *OrderRepository) FindOrderById(ctx context.Context, id int) (*model.OrderDto, error) {
	if r.pool == nil {
		return nil, nil
	}
	var o model.OrderDto
	var desc *string
	var cartId *int
	var userId *int64
	err := r.pool.QueryRow(ctx, `
		SELECT o.order_id, to_char(o.order_date, 'YYYY-MM-DD"T"HH24:MI:SS'), o.order_desc, o.order_fee, o.product_id,
		       c.cart_id, c.user_id
		FROM orders o
		LEFT JOIN carts c ON c.cart_id = o.cart_id
		WHERE o.order_id = $1
	`, id).Scan(&o.OrderId, &o.OrderDate, &desc, &o.OrderFee, &o.ProductId, &cartId, &userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if desc != nil {
		o.OrderDesc = *desc
	}
	if cartId != nil {
		o.Cart = &model.CartDto{CartId: cartId, UserId: userId}
	}
	return &o, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order model.OrderDto) (*model.OrderDto, error) {
	if r.pool == nil {
		id := 1
		order.OrderId = &id
		return &order, nil
	}

	var cartId *int
	if order.Cart != nil && order.Cart.CartId != nil {
		cartId = order.Cart.CartId
	}

	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO orders (order_date, order_desc, order_fee, product_id, cart_id, created_at, updated_at)
		VALUES (COALESCE($1::timestamptz, CURRENT_TIMESTAMP), $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING order_id
	`, nullableTime(order.OrderDate), order.OrderDesc, order.OrderFee, order.ProductId, cartId).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindOrderById(ctx, id)
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, id int, order model.OrderDto) (*model.OrderDto, error) {
	if r.pool == nil {
		order.OrderId = &id
		return &order, nil
	}

	var cartId *int
	if order.Cart != nil && order.Cart.CartId != nil {
		cartId = order.Cart.CartId
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE orders
		SET order_desc = COALESCE($1, order_desc),
		    order_fee = COALESCE($2, order_fee),
		    product_id = COALESCE($3, product_id),
		    cart_id = COALESCE($4, cart_id),
		    updated_at = CURRENT_TIMESTAMP
		WHERE order_id = $5
	`, order.OrderDesc, order.OrderFee, order.ProductId, cartId, id)
	if err != nil {
		return nil, err
	}
	return r.FindOrderById(ctx, id)
}

func (r *OrderRepository) DeleteOrderById(ctx context.Context, id int) error {
	if r.pool == nil {
		return nil
	}
	_, err := r.pool.Exec(ctx, `DELETE FROM orders WHERE order_id = $1`, id)
	return err
}

func (r *OrderRepository) ExistsByOrderId(ctx context.Context, id int) (bool, error) {
	if r.pool == nil {
		return false, nil
	}
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE order_id = $1)`, id).Scan(&exists)
	return exists, err
}

func nullableTime(s string) any {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05", s)
	}
	if err != nil {
		return nil
	}
	return t
}

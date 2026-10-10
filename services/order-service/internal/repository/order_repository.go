package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"com.ecommerce/order-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/fm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type OrderRepository struct {
	db pgxdriver.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: pgxdriver.NewPool(pool)}
}

func (r *OrderRepository) InitSchema(ctx context.Context) error {
	if r.db.Pool == nil {
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
	_, err := r.db.Exec(ctx, schema)
	return err
}

func orderByMod(col bob.Expression, order string) dialect.OrderBy[*dialect.SelectQuery] {
	mod := sm.OrderBy(col)
	if order == "DESC" {
		return mod.Desc()
	}
	return mod.Asc()
}

// -------------------------------------------------------------
// Cart Methods
// -------------------------------------------------------------

type cartRow struct {
	CartId int    `db:"cart_id"`
	UserId *int64 `db:"user_id"`
}

func (row cartRow) toDto(orders []model.OrderDto) model.CartDto {
	cart := model.CartDto{
		CartId: &row.CartId,
		UserId: row.UserId,
		Orders: orders,
	}
	if row.UserId != nil {
		cart.User = &model.UserDto{ID: row.UserId}
	}
	return cart
}

func (r *OrderRepository) FindAllCarts(ctx context.Context) ([]model.CartDto, error) {
	if r.db.Pool == nil {
		return []model.CartDto{}, nil
	}
	query := psql.Select(
		sm.Columns("cart_id", "user_id"),
		sm.From("carts"),
		sm.OrderBy(psql.Quote("cart_id")).Asc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[cartRow]())
	if err != nil {
		return nil, err
	}

	carts := make([]model.CartDto, 0, len(rows))
	for _, row := range rows {
		orders, _ := r.findOrdersByCartId(ctx, row.CartId)
		carts = append(carts, row.toDto(orders))
	}
	return carts, nil
}

func (r *OrderRepository) FindAllCartsPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.CartDto, int64, error) {
	if r.db.Pool == nil {
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

	countQuery := psql.Select(
		sm.Columns("COUNT(*)"),
		sm.From("carts"),
	)
	total, err := bob.One(ctx, r.db, countQuery, scan.SingleColumnMapper[int64])
	if err != nil {
		return nil, 0, err
	}

	query := psql.Select(
		sm.Columns("cart_id", "user_id"),
		sm.From("carts"),
		orderByMod(psql.Quote(col), order),
		sm.Limit(size),
		sm.Offset(page*size),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[cartRow]())
	if err != nil {
		return nil, 0, err
	}

	carts := make([]model.CartDto, 0, len(rows))
	for _, row := range rows {
		orders, _ := r.findOrdersByCartId(ctx, row.CartId)
		carts = append(carts, row.toDto(orders))
	}
	return carts, total, nil
}

func (r *OrderRepository) FindCartById(ctx context.Context, id int) (*model.CartDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns("cart_id", "user_id"),
		sm.From("carts"),
		sm.Where(psql.Quote("cart_id").EQ(psql.Arg(id))),
	)
	row, err := bob.One(ctx, r.db, query, scan.StructMapper[cartRow]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	orders, _ := r.findOrdersByCartId(ctx, row.CartId)
	cart := row.toDto(orders)
	return &cart, nil
}

func (r *OrderRepository) SaveCart(ctx context.Context, cart model.CartDto) (*model.CartDto, error) {
	if r.db.Pool == nil {
		id := 1
		cart.CartId = &id
		return &cart, nil
	}

	query := psql.Insert(
		im.Into("carts", "user_id", "created_at", "updated_at"),
		im.Values(
			psql.Arg(cart.UserId),
			psql.Raw("CURRENT_TIMESTAMP"),
			psql.Raw("CURRENT_TIMESTAMP"),
		),
		im.Returning("cart_id"),
	)
	id, err := bob.One(ctx, r.db, query, scan.SingleColumnMapper[int])
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
	if r.db.Pool == nil {
		cart.CartId = &id
		return &cart, nil
	}

	query := psql.Update(
		um.Table("carts"),
		um.Set(
			psql.Quote("user_id").EQ(psql.F("COALESCE", psql.Arg(cart.UserId), psql.Quote("user_id"))),
			psql.Quote("updated_at").EQ(psql.Raw("CURRENT_TIMESTAMP")),
		),
		um.Where(psql.Quote("cart_id").EQ(psql.Arg(id))),
	)
	if _, err := bob.Exec(ctx, r.db, query); err != nil {
		return nil, err
	}
	return r.FindCartById(ctx, id)
}

func (r *OrderRepository) DeleteCartById(ctx context.Context, id int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("carts"),
		dm.Where(psql.Quote("cart_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

func (r *OrderRepository) findOrdersByCartId(ctx context.Context, cartId int) ([]model.OrderDto, error) {
	query := psql.Select(
		sm.Columns(
			psql.Quote("order_id"),
			psql.F("to_char", psql.Quote("order_date"), psql.S(`YYYY-MM-DD"T"HH24:MI:SS`))(fm.As("order_date")),
			psql.Quote("order_desc"),
			psql.Quote("order_fee"),
			psql.Quote("product_id"),
		),
		sm.From("orders"),
		sm.Where(psql.Quote("cart_id").EQ(psql.Arg(cartId))),
		sm.OrderBy(psql.Quote("order_id")).Asc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[orderRow]())
	if err != nil {
		return nil, err
	}

	var orders []model.OrderDto
	for _, row := range rows {
		orders = append(orders, row.toDto())
	}
	return orders, nil
}

// -------------------------------------------------------------
// Order Methods
// -------------------------------------------------------------

type orderRow struct {
	OrderId   int      `db:"order_id"`
	OrderDate string   `db:"order_date"`
	OrderDesc *string  `db:"order_desc"`
	OrderFee  *float64 `db:"order_fee"`
	ProductId *int     `db:"product_id"`
	CartId    *int     `db:"cart_id"`
	UserId    *int64   `db:"user_id"`
}

func (row orderRow) toDto() model.OrderDto {
	orderId := row.OrderId
	order := model.OrderDto{
		OrderId:   &orderId,
		OrderDate: row.OrderDate,
		OrderFee:  row.OrderFee,
		ProductId: row.ProductId,
	}
	if row.OrderDesc != nil {
		order.OrderDesc = *row.OrderDesc
	}
	if row.CartId != nil {
		order.Cart = &model.CartDto{CartId: row.CartId, UserId: row.UserId}
	}
	return order
}

func orderColumns() []any {
	return []any{
		psql.Quote("o", "order_id"),
		psql.F("to_char", psql.Quote("o", "order_date"), psql.S(`YYYY-MM-DD"T"HH24:MI:SS`))(fm.As("order_date")),
		psql.Quote("o", "order_desc"),
		psql.Quote("o", "order_fee"),
		psql.Quote("o", "product_id"),
		psql.Quote("c", "cart_id"),
		psql.Quote("c", "user_id"),
	}
}

func (r *OrderRepository) FindAllOrders(ctx context.Context) ([]model.OrderDto, error) {
	if r.db.Pool == nil {
		return []model.OrderDto{}, nil
	}
	query := psql.Select(
		sm.Columns(orderColumns()...),
		sm.From("orders o", sm.LeftJoin("carts c").On(psql.Raw("c.cart_id = o.cart_id"))),
		sm.OrderBy(psql.Quote("o", "order_id")).Asc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[orderRow]())
	if err != nil {
		return nil, err
	}

	orders := make([]model.OrderDto, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, row.toDto())
	}
	return orders, nil
}

func (r *OrderRepository) FindAllOrdersPaged(ctx context.Context, page, size int, sortBy, sortOrder string) ([]model.OrderDto, int64, error) {
	if r.db.Pool == nil {
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

	countQuery := psql.Select(
		sm.Columns("COUNT(*)"),
		sm.From("orders"),
	)
	total, err := bob.One(ctx, r.db, countQuery, scan.SingleColumnMapper[int64])
	if err != nil {
		return nil, 0, err
	}

	query := psql.Select(
		sm.Columns(orderColumns()...),
		sm.From("orders o", sm.LeftJoin("carts c").On(psql.Raw("c.cart_id = o.cart_id"))),
		orderByMod(psql.Quote("o", col), order),
		sm.Limit(size),
		sm.Offset(page*size),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[orderRow]())
	if err != nil {
		return nil, 0, err
	}

	orders := make([]model.OrderDto, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, row.toDto())
	}
	return orders, total, nil
}

func (r *OrderRepository) FindOrderById(ctx context.Context, id int) (*model.OrderDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	query := psql.Select(
		sm.Columns(orderColumns()...),
		sm.From("orders o", sm.LeftJoin("carts c").On(psql.Raw("c.cart_id = o.cart_id"))),
		sm.Where(psql.Quote("o", "order_id").EQ(psql.Arg(id))),
	)
	row, err := bob.One(ctx, r.db, query, scan.StructMapper[orderRow]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	order := row.toDto()
	return &order, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order model.OrderDto) (*model.OrderDto, error) {
	if r.db.Pool == nil {
		id := 1
		order.OrderId = &id
		return &order, nil
	}

	var cartId *int
	if order.Cart != nil && order.Cart.CartId != nil {
		cartId = order.Cart.CartId
	}

	query := psql.Insert(
		im.Into("orders", "order_date", "order_desc", "order_fee", "product_id", "cart_id", "created_at", "updated_at"),
		im.Values(
			psql.F("COALESCE", psql.Cast(psql.Arg(nullableTime(order.OrderDate)), "timestamptz"), psql.Raw("CURRENT_TIMESTAMP")),
			psql.Arg(order.OrderDesc),
			psql.Arg(order.OrderFee),
			psql.Arg(order.ProductId),
			psql.Arg(cartId),
			psql.Raw("CURRENT_TIMESTAMP"),
			psql.Raw("CURRENT_TIMESTAMP"),
		),
		im.Returning("order_id"),
	)
	id, err := bob.One(ctx, r.db, query, scan.SingleColumnMapper[int])
	if err != nil {
		return nil, err
	}
	return r.FindOrderById(ctx, id)
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, id int, order model.OrderDto) (*model.OrderDto, error) {
	if r.db.Pool == nil {
		order.OrderId = &id
		return &order, nil
	}

	var cartId *int
	if order.Cart != nil && order.Cart.CartId != nil {
		cartId = order.Cart.CartId
	}

	query := psql.Update(
		um.Table("orders"),
		um.Set(
			psql.Quote("order_desc").EQ(psql.F("COALESCE", psql.Arg(order.OrderDesc), psql.Quote("order_desc"))),
			psql.Quote("order_fee").EQ(psql.F("COALESCE", psql.Arg(order.OrderFee), psql.Quote("order_fee"))),
			psql.Quote("product_id").EQ(psql.F("COALESCE", psql.Arg(order.ProductId), psql.Quote("product_id"))),
			psql.Quote("cart_id").EQ(psql.F("COALESCE", psql.Arg(cartId), psql.Quote("cart_id"))),
			psql.Quote("updated_at").EQ(psql.Raw("CURRENT_TIMESTAMP")),
		),
		um.Where(psql.Quote("order_id").EQ(psql.Arg(id))),
	)
	if _, err := bob.Exec(ctx, r.db, query); err != nil {
		return nil, err
	}
	return r.FindOrderById(ctx, id)
}

func (r *OrderRepository) DeleteOrderById(ctx context.Context, id int) error {
	if r.db.Pool == nil {
		return nil
	}
	query := psql.Delete(
		dm.From("orders"),
		dm.Where(psql.Quote("order_id").EQ(psql.Arg(id))),
	)
	_, err := bob.Exec(ctx, r.db, query)
	return err
}

func (r *OrderRepository) ExistsByOrderId(ctx context.Context, id int) (bool, error) {
	if r.db.Pool == nil {
		return false, nil
	}
	query := psql.Select(
		sm.Columns(psql.Exists(psql.Select(
			sm.Columns("1"),
			sm.From("orders"),
			sm.Where(psql.Quote("order_id").EQ(psql.Arg(id))),
		))),
	)
	return bob.One(ctx, r.db, query, scan.SingleColumnMapper[bool])
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

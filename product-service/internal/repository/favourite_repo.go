package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"com.ecommerce/product-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	pgxdriver "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/scan"
)

type FavouriteRepository struct {
	db pgxdriver.Pool
}

func NewFavouriteRepository(pool *pgxpool.Pool) *FavouriteRepository {
	return &FavouriteRepository{db: pgxdriver.NewPool(pool)}
}

type favouriteRow struct {
	UserID    int       `db:"user_id"`
	ProductID int       `db:"product_id"`
	LikeDate  time.Time `db:"like_date"`
}

func (row favouriteRow) toDto() model.FavouriteDto {
	return model.FavouriteDto{
		UserID:    row.UserID,
		ProductID: row.ProductID,
		LikeDate:  formatLikeDate(row.LikeDate),
	}
}

func parseLikeDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Parse(time.RFC3339, dateStr)
}

func formatLikeDate(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func (r *FavouriteRepository) FindAll(ctx context.Context) ([]model.FavouriteDto, error) {
	if r.db.Pool == nil {
		return []model.FavouriteDto{}, nil
	}
	query := psql.Select(
		sm.Columns(
			psql.Quote("user_id"),
			psql.Quote("product_id"),
			psql.Quote("like_date"),
		),
		sm.From("favourites"),
		sm.OrderBy(psql.Quote("like_date")).Desc(),
	)
	rows, err := bob.All(ctx, r.db, query, scan.StructMapper[favouriteRow]())
	if err != nil {
		return nil, err
	}
	result := make([]model.FavouriteDto, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toDto())
	}
	return result, nil
}

func (r *FavouriteRepository) FindByID(ctx context.Context, id model.FavouriteID) (*model.FavouriteDto, error) {
	if r.db.Pool == nil {
		return nil, nil
	}
	t, err := parseLikeDate(id.LikeDate)
	if err != nil {
		return nil, fmt.Errorf("invalid likeDate: %w", err)
	}
	query := psql.Select(
		sm.Columns(
			psql.Quote("user_id"),
			psql.Quote("product_id"),
			psql.Quote("like_date"),
		),
		sm.From("favourites"),
		sm.Where(psql.And(
			psql.Quote("user_id").EQ(psql.Arg(id.UserID)),
			psql.Quote("product_id").EQ(psql.Arg(id.ProductID)),
			psql.Quote("like_date").EQ(psql.Arg(t)),
		)),
		sm.Limit(1),
	)
	row, err := bob.One(ctx, r.db, query, scan.StructMapper[favouriteRow]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	dto := row.toDto()
	return &dto, nil
}

func (r *FavouriteRepository) Save(ctx context.Context, dto model.FavouriteDto) (*model.FavouriteDto, error) {
	if r.db.Pool == nil {
		return &dto, nil
	}
	t, err := parseLikeDate(dto.LikeDate)
	if err != nil {
		t = time.Now()
		dto.LikeDate = formatLikeDate(t)
	}

	query := psql.Insert(
		im.Into("favourites", "user_id", "product_id", "like_date"),
		im.Values(
			psql.Arg(dto.UserID),
			psql.Arg(dto.ProductID),
			psql.Arg(t),
		),
		im.OnConflict("user_id", "product_id", "like_date").
			DoUpdate(im.Set(psql.Quote("updated_at").EQ(psql.Raw("CURRENT_TIMESTAMP")))),
	)
	if _, err := bob.Exec(ctx, r.db, query); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (r *FavouriteRepository) DeleteByID(ctx context.Context, id model.FavouriteID) (bool, error) {
	if r.db.Pool == nil {
		return true, nil
	}
	t, err := parseLikeDate(id.LikeDate)
	if err != nil {
		return false, fmt.Errorf("invalid likeDate: %w", err)
	}

	query := psql.Delete(
		dm.From("favourites"),
		dm.Where(psql.And(
			psql.Quote("user_id").EQ(psql.Arg(id.UserID)),
			psql.Quote("product_id").EQ(psql.Arg(id.ProductID)),
			psql.Quote("like_date").EQ(psql.Arg(t)),
		)),
	)
	result, err := bob.Exec(ctx, r.db, query)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected >= 0, nil
}

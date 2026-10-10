package repository

import (
	"context"
	"fmt"
	"time"

	"com.ecommerce/product-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FavouriteRepository struct {
	db *pgxpool.Pool
}

func NewFavouriteRepository(db *pgxpool.Pool) *FavouriteRepository {
	return &FavouriteRepository{db: db}
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
	if r.db == nil {
		return []model.FavouriteDto{}, nil
	}
	rows, err := r.db.Query(ctx, `SELECT user_id, product_id, like_date FROM favourites ORDER BY like_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.FavouriteDto
	for rows.Next() {
		var (
			userID    int
			productID int
			likeDate  time.Time
		)
		if err := rows.Scan(&userID, &productID, &likeDate); err != nil {
			return nil, err
		}
		result = append(result, model.FavouriteDto{
			UserID:    userID,
			ProductID: productID,
			LikeDate:  formatLikeDate(likeDate),
		})
	}
	if result == nil {
		result = []model.FavouriteDto{}
	}
	return result, nil
}

func (r *FavouriteRepository) FindByID(ctx context.Context, id model.FavouriteID) (*model.FavouriteDto, error) {
	if r.db == nil {
		return nil, nil
	}
	t, err := parseLikeDate(id.LikeDate)
	if err != nil {
		return nil, fmt.Errorf("invalid likeDate: %w", err)
	}

	var (
		userID    int
		productID int
		likeDate  time.Time
	)
	err = r.db.QueryRow(ctx, `SELECT user_id, product_id, like_date FROM favourites WHERE user_id = $1 AND product_id = $2 AND like_date = $3 LIMIT 1`, id.UserID, id.ProductID, t).Scan(&userID, &productID, &likeDate)
	if err != nil {
		return nil, nil
	}

	return &model.FavouriteDto{
		UserID:    userID,
		ProductID: productID,
		LikeDate:  formatLikeDate(likeDate),
	}, nil
}

func (r *FavouriteRepository) Save(ctx context.Context, dto model.FavouriteDto) (*model.FavouriteDto, error) {
	if r.db == nil {
		return &dto, nil
	}
	t, err := parseLikeDate(dto.LikeDate)
	if err != nil {
		t = time.Now()
		dto.LikeDate = formatLikeDate(t)
	}

	query := `
		INSERT INTO favourites (user_id, product_id, like_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id, like_date) DO UPDATE
		SET updated_at = CURRENT_TIMESTAMP
	`
	_, err = r.db.Exec(ctx, query, dto.UserID, dto.ProductID, t)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (r *FavouriteRepository) DeleteByID(ctx context.Context, id model.FavouriteID) (bool, error) {
	if r.db == nil {
		return true, nil
	}
	t, err := parseLikeDate(id.LikeDate)
	if err != nil {
		return false, fmt.Errorf("invalid likeDate: %w", err)
	}

	cmd, err := r.db.Exec(ctx, `DELETE FROM favourites WHERE user_id = $1 AND product_id = $2 AND like_date = $3`, id.UserID, id.ProductID, t)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() >= 0, nil
}

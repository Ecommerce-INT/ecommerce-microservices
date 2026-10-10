package repository

import (
	"context"
	"errors"
	"math"

	"com.ecommerce/tax-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxRepository struct {
	db *pgxpool.Pool
}

func NewTaxRepository(db *pgxpool.Pool) *TaxRepository {
	return &TaxRepository{db: db}
}

// =================== TaxClass ===================

func (r *TaxRepository) FindAllTaxClasses(ctx context.Context) ([]model.TaxClassVm, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM tax_class ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TaxClassVm
	for rows.Next() {
		var item model.TaxClassVm
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.TaxClassVm{}
	}
	return list, nil
}

func (r *TaxRepository) FindTaxClassByID(ctx context.Context, id int64) (*model.TaxClassVm, error) {
	var item model.TaxClassVm
	err := r.db.QueryRow(ctx, "SELECT id, name FROM tax_class WHERE id = $1", id).Scan(&item.ID, &item.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *TaxRepository) ExistsTaxClassByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_class WHERE name = $1)", name).Scan(&exists)
	return exists, err
}

func (r *TaxRepository) ExistsTaxClassByNameNotID(ctx context.Context, name string, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_class WHERE name = $1 AND id != $2)", name, id).Scan(&exists)
	return exists, err
}

func (r *TaxRepository) ExistsTaxClassByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_class WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

func (r *TaxRepository) CreateTaxClass(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, "INSERT INTO tax_class (name, created_on) VALUES ($1, NOW()) RETURNING id", name).Scan(&id)
	return id, err
}

func (r *TaxRepository) UpdateTaxClass(ctx context.Context, id int64, name string) error {
	_, err := r.db.Exec(ctx, "UPDATE tax_class SET name = $1, last_modified_on = NOW() WHERE id = $2", name, id)
	return err
}

func (r *TaxRepository) DeleteTaxClass(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tax_class WHERE id = $1", id)
	return err
}

func (r *TaxRepository) GetPageableTaxClasses(ctx context.Context, pageNo, pageSize int) (*model.TaxClassListGetVm, error) {
	var totalElements int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM tax_class").Scan(&totalElements)
	if err != nil {
		return nil, err
	}

	offset := pageNo * pageSize
	rows, err := r.db.Query(ctx, "SELECT id, name FROM tax_class ORDER BY id ASC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TaxClassVm
	for rows.Next() {
		var item model.TaxClassVm
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.TaxClassVm{}
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(pageSize)))
	}
	isLast := (pageNo + 1) >= totalPages

	return &model.TaxClassListGetVm{
		TaxClassContent: list,
		PageNo:          pageNo,
		PageSize:        pageSize,
		TotalElements:   totalElements,
		TotalPages:      totalPages,
		IsLast:          isLast,
	}, nil
}

// =================== TaxRate ===================

func (r *TaxRepository) FindAllTaxRates(ctx context.Context) ([]model.TaxRateVm, error) {
	query := `
		SELECT tr.id, tr.rate, tr.zip_code, tr.tax_class_id, COALESCE(tc.name, ''), tr.state_or_province_id, tr.country_id
		FROM tax_rate tr
		LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
		ORDER BY tr.id ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TaxRateVm
	for rows.Next() {
		var item model.TaxRateVm
		if err := rows.Scan(&item.ID, &item.Rate, &item.ZipCode, &item.TaxClassId, &item.TaxClassName, &item.StateOrProvinceId, &item.CountryId); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.TaxRateVm{}
	}
	return list, nil
}

func (r *TaxRepository) FindTaxRateByID(ctx context.Context, id int64) (*model.TaxRateVm, error) {
	query := `
		SELECT tr.id, tr.rate, tr.zip_code, tr.tax_class_id, COALESCE(tc.name, ''), tr.state_or_province_id, tr.country_id
		FROM tax_rate tr
		LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
		WHERE tr.id = $1`
	var item model.TaxRateVm
	err := r.db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Rate, &item.ZipCode, &item.TaxClassId, &item.TaxClassName, &item.StateOrProvinceId, &item.CountryId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *TaxRepository) ExistsTaxRateByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_rate WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

func (r *TaxRepository) CreateTaxRate(ctx context.Context, req *model.TaxRatePostVm) (int64, error) {
	query := `
		INSERT INTO tax_rate (rate, zip_code, tax_class_id, state_or_province_id, country_id, created_on)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id`
	var id int64
	err := r.db.QueryRow(ctx, query, req.Rate, req.ZipCode, req.TaxClassId, req.StateOrProvinceId, req.CountryId).Scan(&id)
	return id, err
}

func (r *TaxRepository) UpdateTaxRate(ctx context.Context, id int64, req *model.TaxRatePostVm) error {
	query := `
		UPDATE tax_rate
		SET rate = $1, zip_code = $2, tax_class_id = $3, state_or_province_id = $4, country_id = $5, last_modified_on = NOW()
		WHERE id = $6`
	_, err := r.db.Exec(ctx, query, req.Rate, req.ZipCode, req.TaxClassId, req.StateOrProvinceId, req.CountryId, id)
	return err
}

func (r *TaxRepository) DeleteTaxRate(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tax_rate WHERE id = $1", id)
	return err
}

func (r *TaxRepository) GetPageableTaxRates(ctx context.Context, pageNo, pageSize int) (*model.TaxRateListGetVm, error) {
	var totalElements int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM tax_rate").Scan(&totalElements)
	if err != nil {
		return nil, err
	}

	offset := pageNo * pageSize
	query := `
		SELECT tr.id, tr.rate, tr.zip_code, tr.tax_class_id, COALESCE(tc.name, ''), tr.state_or_province_id, tr.country_id
		FROM tax_rate tr
		LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
		ORDER BY tr.id ASC
		LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TaxRateVm
	for rows.Next() {
		var item model.TaxRateVm
		if err := rows.Scan(&item.ID, &item.Rate, &item.ZipCode, &item.TaxClassId, &item.TaxClassName, &item.StateOrProvinceId, &item.CountryId); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.TaxRateVm{}
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(pageSize)))
	}
	isLast := (pageNo + 1) >= totalPages

	return &model.TaxRateListGetVm{
		TaxRateContent: list,
		PageNo:         pageNo,
		PageSize:       pageSize,
		TotalElements:  totalElements,
		TotalPages:     totalPages,
		IsLast:         isLast,
	}, nil
}

func (r *TaxRepository) GetTaxPercent(ctx context.Context, taxClassID, countryID int64, stateOrProvinceID *int64, zipCode *string) (float64, error) {
	query := `
		SELECT rate FROM tax_rate tr
		WHERE tr.country_id = $1
		  AND (tr.state_or_province_id = $2 OR tr.state_or_province_id IS NULL)
		  AND tr.tax_class_id = $3
		  AND (tr.zip_code = $4 OR TRIM(COALESCE(tr.zip_code, '')) = '' OR tr.zip_code IS NULL)
		FETCH FIRST 1 ROWS ONLY`
	var rate float64
	err := r.db.QueryRow(ctx, query, countryID, stateOrProvinceID, taxClassID, zipCode).Scan(&rate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0.0, nil
		}
		return 0.0, err
	}
	return rate, nil
}

func (r *TaxRepository) GetBatchTaxRates(ctx context.Context, taxClassIDs []int64, countryID int64, stateOrProvinceID *int64, zipCode *string) ([]model.TaxRateVm, error) {
	query := `
		SELECT tr.id, tr.rate, tr.zip_code, tr.tax_class_id, COALESCE(tc.name, ''), tr.state_or_province_id, tr.country_id
		FROM tax_rate tr
		LEFT JOIN tax_class tc ON tr.tax_class_id = tc.id
		WHERE tr.country_id = $1
		  AND (tr.state_or_province_id = $2 OR tr.state_or_province_id IS NULL)
		  AND (tr.zip_code = $3 OR tr.zip_code IS NULL)
		  AND tr.tax_class_id = ANY($4)`
	rows, err := r.db.Query(ctx, query, countryID, stateOrProvinceID, zipCode, taxClassIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TaxRateVm
	for rows.Next() {
		var item model.TaxRateVm
		if err := rows.Scan(&item.ID, &item.Rate, &item.ZipCode, &item.TaxClassId, &item.TaxClassName, &item.StateOrProvinceId, &item.CountryId); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.TaxRateVm{}
	}
	return list, nil
}

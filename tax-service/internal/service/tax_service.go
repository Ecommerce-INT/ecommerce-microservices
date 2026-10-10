package service

import (
	"context"
	"errors"
	"fmt"

	"com.ecommerce/tax-service/internal/model"
	"com.ecommerce/tax-service/internal/repository"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrDuplicated = errors.New("duplicated")
)

type TaxService struct {
	repo *repository.TaxRepository
}

func NewTaxService(repo *repository.TaxRepository) *TaxService {
	return &TaxService{repo: repo}
}

// =================== TaxClass ===================

func (s *TaxService) FindAllTaxClasses(ctx context.Context) ([]model.TaxClassVm, error) {
	return s.repo.FindAllTaxClasses(ctx)
}

func (s *TaxService) FindTaxClassByID(ctx context.Context, id int64) (*model.TaxClassVm, error) {
	tc, err := s.repo.FindTaxClassByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tc == nil {
		return nil, fmt.Errorf("%w: tax class with id %d not found", ErrNotFound, id)
	}
	return tc, nil
}

func (s *TaxService) CreateTaxClass(ctx context.Context, req *model.TaxClassPostVm) (*model.TaxClassVm, error) {
	exists, err := s.repo.ExistsTaxClassByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: name %s already exists", ErrDuplicated, req.Name)
	}

	id, err := s.repo.CreateTaxClass(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &model.TaxClassVm{ID: id, Name: req.Name}, nil
}

func (s *TaxService) UpdateTaxClass(ctx context.Context, id int64, req *model.TaxClassPostVm) error {
	existsID, err := s.repo.ExistsTaxClassByID(ctx, id)
	if err != nil {
		return err
	}
	if !existsID {
		return fmt.Errorf("%w: tax class with id %d not found", ErrNotFound, id)
	}

	existsName, err := s.repo.ExistsTaxClassByNameNotID(ctx, req.Name, id)
	if err != nil {
		return err
	}
	if existsName {
		return fmt.Errorf("%w: name %s already exists", ErrDuplicated, req.Name)
	}

	return s.repo.UpdateTaxClass(ctx, id, req.Name)
}

func (s *TaxService) DeleteTaxClass(ctx context.Context, id int64) error {
	existsID, err := s.repo.ExistsTaxClassByID(ctx, id)
	if err != nil {
		return err
	}
	if !existsID {
		return fmt.Errorf("%w: tax class with id %d not found", ErrNotFound, id)
	}
	return s.repo.DeleteTaxClass(ctx, id)
}

func (s *TaxService) GetPageableTaxClasses(ctx context.Context, pageNo, pageSize int) (*model.TaxClassListGetVm, error) {
	if pageNo < 0 {
		pageNo = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPageableTaxClasses(ctx, pageNo, pageSize)
}

// =================== TaxRate ===================

func (s *TaxService) FindAllTaxRates(ctx context.Context) ([]model.TaxRateVm, error) {
	return s.repo.FindAllTaxRates(ctx)
}

func (s *TaxService) FindTaxRateByID(ctx context.Context, id int64) (*model.TaxRateVm, error) {
	tr, err := s.repo.FindTaxRateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tr == nil {
		return nil, fmt.Errorf("%w: tax rate with id %d not found", ErrNotFound, id)
	}
	return tr, nil
}

func (s *TaxService) CreateTaxRate(ctx context.Context, req *model.TaxRatePostVm) (*model.TaxRateVm, error) {
	existsClass, err := s.repo.ExistsTaxClassByID(ctx, req.TaxClassId)
	if err != nil {
		return nil, err
	}
	if !existsClass {
		return nil, fmt.Errorf("%w: tax class with id %d not found", ErrNotFound, req.TaxClassId)
	}

	id, err := s.repo.CreateTaxRate(ctx, req)
	if err != nil {
		return nil, err
	}

	tc, _ := s.repo.FindTaxClassByID(ctx, req.TaxClassId)
	className := ""
	if tc != nil {
		className = tc.Name
	}

	return &model.TaxRateVm{
		ID:                id,
		Rate:              req.Rate,
		ZipCode:           req.ZipCode,
		TaxClassId:        req.TaxClassId,
		TaxClassName:      className,
		StateOrProvinceId: req.StateOrProvinceId,
		CountryId:         req.CountryId,
	}, nil
}

func (s *TaxService) UpdateTaxRate(ctx context.Context, id int64, req *model.TaxRatePostVm) error {
	existsRate, err := s.repo.ExistsTaxRateByID(ctx, id)
	if err != nil {
		return err
	}
	if !existsRate {
		return fmt.Errorf("%w: tax rate with id %d not found", ErrNotFound, id)
	}

	existsClass, err := s.repo.ExistsTaxClassByID(ctx, req.TaxClassId)
	if err != nil {
		return err
	}
	if !existsClass {
		return fmt.Errorf("%w: tax class with id %d not found", ErrNotFound, req.TaxClassId)
	}

	return s.repo.UpdateTaxRate(ctx, id, req)
}

func (s *TaxService) DeleteTaxRate(ctx context.Context, id int64) error {
	existsRate, err := s.repo.ExistsTaxRateByID(ctx, id)
	if err != nil {
		return err
	}
	if !existsRate {
		return fmt.Errorf("%w: tax rate with id %d not found", ErrNotFound, id)
	}
	return s.repo.DeleteTaxRate(ctx, id)
}

func (s *TaxService) GetPageableTaxRates(ctx context.Context, pageNo, pageSize int) (*model.TaxRateListGetVm, error) {
	if pageNo < 0 {
		pageNo = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPageableTaxRates(ctx, pageNo, pageSize)
}

func (s *TaxService) GetTaxPercent(ctx context.Context, taxClassID, countryID int64, stateOrProvinceID *int64, zipCode *string) (float64, error) {
	return s.repo.GetTaxPercent(ctx, taxClassID, countryID, stateOrProvinceID, zipCode)
}

func (s *TaxService) GetBatchTaxRates(ctx context.Context, taxClassIDs []int64, countryID int64, stateOrProvinceID *int64, zipCode *string) ([]model.TaxRateVm, error) {
	return s.repo.GetBatchTaxRates(ctx, taxClassIDs, countryID, stateOrProvinceID, zipCode)
}

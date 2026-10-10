package service

import (
	"context"
	"errors"
	"fmt"

	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/repository"
)

var ErrNotFound = errors.New("not found")

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// =================== Categories ===================

func (s *ProductService) FindAllCategories(ctx context.Context) ([]model.CategoryDto, error) {
	return s.repo.FindAllCategories(ctx)
}

func (s *ProductService) FindCategoryByID(ctx context.Context, id int) (*model.CategoryDto, error) {
	c, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("%w: category with id %d not found", ErrNotFound, id)
	}
	return c, nil
}

func (s *ProductService) SaveCategory(ctx context.Context, dto *model.CategoryDto) (*model.CategoryDto, error) {
	id, err := s.repo.SaveCategory(ctx, dto)
	if err != nil {
		return nil, err
	}
	dto.CategoryId = id
	return dto, nil
}

func (s *ProductService) UpdateCategory(ctx context.Context, id int, dto *model.CategoryDto) (*model.CategoryDto, error) {
	err := s.repo.UpdateCategory(ctx, id, dto)
	if err != nil {
		return nil, err
	}
	dto.CategoryId = id
	return dto, nil
}

func (s *ProductService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

// =================== Products ===================

func (s *ProductService) FindAllProducts(ctx context.Context) ([]model.ProductDto, error) {
	return s.repo.FindAllProducts(ctx)
}

func (s *ProductService) FindProductByID(ctx context.Context, id int) (*model.ProductDto, error) {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("%w: product with id %d not found", ErrNotFound, id)
	}
	return p, nil
}

func (s *ProductService) SaveProduct(ctx context.Context, dto *model.ProductDto) (*model.ProductDto, error) {
	id, err := s.repo.SaveProduct(ctx, dto)
	if err != nil {
		return nil, err
	}
	dto.ProductId = id
	return dto, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, dto *model.ProductDto) (*model.ProductDto, error) {
	err := s.repo.UpdateProduct(ctx, id, dto)
	if err != nil {
		return nil, err
	}
	dto.ProductId = id
	return dto, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.DeleteProduct(ctx, id)
}

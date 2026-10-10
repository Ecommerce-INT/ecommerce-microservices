package service

import (
	"context"
	"errors"
	"fmt"

	"com.ecommerce/shipping-service/internal/model"
	"com.ecommerce/shipping-service/internal/repository"
)

var ErrNotFound = errors.New("not found")

type ShippingService struct {
	repo *repository.ShippingRepository
}

func NewShippingService(repo *repository.ShippingRepository) *ShippingService {
	return &ShippingService{repo: repo}
}

func (s *ShippingService) FindAll(ctx context.Context) ([]model.OrderItemDto, error) {
	return s.repo.FindAll(ctx)
}

func (s *ShippingService) FindByID(ctx context.Context, orderId, productId int) (*model.OrderItemDto, error) {
	item, err := s.repo.FindByID(ctx, orderId, productId)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("%w: order item not found for orderId: %d, productId: %d", ErrNotFound, orderId, productId)
	}
	return item, nil
}

func (s *ShippingService) Save(ctx context.Context, dto *model.OrderItemDto) (*model.OrderItemDto, error) {
	err := s.repo.Save(ctx, dto)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *ShippingService) Update(ctx context.Context, dto *model.OrderItemDto) (*model.OrderItemDto, error) {
	return s.Save(ctx, dto)
}

func (s *ShippingService) DeleteByID(ctx context.Context, orderId, productId int) error {
	return s.repo.DeleteByID(ctx, orderId, productId)
}

package service

import (
	"context"
	"errors"
	"math"

	"com.ecommerce/order-service/internal/model"
	"com.ecommerce/order-service/internal/repository"
)

var (
	ErrCartNotFound  = errors.New("cart not found")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// -------------------------------------------------------------
// Cart Services
// -------------------------------------------------------------

func (s *OrderService) FindAllCarts(ctx context.Context) ([]model.CartDto, error) {
	return s.repo.FindAllCarts(ctx)
}

func (s *OrderService) FindAllCartsPaged(ctx context.Context, page, size int, sortBy, sortOrder string) (*model.PageResponse[model.CartDto], error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}

	carts, total, err := s.repo.FindAllCartsPaged(ctx, page, size, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total == 0 {
		totalPages = 0
	}

	return &model.PageResponse[model.CartDto]{
		Content:          carts,
		TotalElements:    total,
		TotalPages:       totalPages,
		Size:             size,
		Number:           page,
		NumberOfElements: len(carts),
		First:            page == 0,
		Last:             page >= totalPages-1,
		Empty:            len(carts) == 0,
	}, nil
}

func (s *OrderService) FindCartById(ctx context.Context, id int) (*model.CartDto, error) {
	cart, err := s.repo.FindCartById(ctx, id)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}
	return cart, nil
}

func (s *OrderService) SaveCart(ctx context.Context, cart model.CartDto) (*model.CartDto, error) {
	return s.repo.SaveCart(ctx, cart)
}

func (s *OrderService) UpdateCart(ctx context.Context, id int, cart model.CartDto) (*model.CartDto, error) {
	updated, err := s.repo.UpdateCart(ctx, id, cart)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrCartNotFound
	}
	return updated, nil
}

func (s *OrderService) DeleteCartById(ctx context.Context, id int) error {
	return s.repo.DeleteCartById(ctx, id)
}

// -------------------------------------------------------------
// Order Services
// -------------------------------------------------------------

func (s *OrderService) FindAllOrders(ctx context.Context) ([]model.OrderDto, error) {
	return s.repo.FindAllOrders(ctx)
}

func (s *OrderService) FindAllOrdersPaged(ctx context.Context, page, size int, sortBy, sortOrder string) (*model.PageResponse[model.OrderDto], error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}

	orders, total, err := s.repo.FindAllOrdersPaged(ctx, page, size, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total == 0 {
		totalPages = 0
	}

	return &model.PageResponse[model.OrderDto]{
		Content:          orders,
		TotalElements:    total,
		TotalPages:       totalPages,
		Size:             size,
		Number:           page,
		NumberOfElements: len(orders),
		First:            page == 0,
		Last:             page >= totalPages-1,
		Empty:            len(orders) == 0,
	}, nil
}

func (s *OrderService) FindOrderById(ctx context.Context, id int) (*model.OrderDto, error) {
	order, err := s.repo.FindOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderService) SaveOrder(ctx context.Context, order model.OrderDto) (*model.OrderDto, error) {
	return s.repo.SaveOrder(ctx, order)
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, order model.OrderDto) (*model.OrderDto, error) {
	updated, err := s.repo.UpdateOrder(ctx, id, order)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrOrderNotFound
	}
	return updated, nil
}

func (s *OrderService) DeleteOrderById(ctx context.Context, id int) error {
	return s.repo.DeleteOrderById(ctx, id)
}

func (s *OrderService) ExistsByOrderId(ctx context.Context, id int) (bool, error) {
	return s.repo.ExistsByOrderId(ctx, id)
}

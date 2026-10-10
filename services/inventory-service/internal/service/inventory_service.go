package service

import (
	"context"

	"com.ecommerce/inventory-service/internal/model"
	"com.ecommerce/inventory-service/internal/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) IsInStock(ctx context.Context, productNames []string) ([]model.InventoryResponse, error) {
	items, err := s.repo.FindByProductNames(ctx, productNames)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]int)
	for _, it := range items {
		itemMap[it.ProductName] = it.Quantity
	}

	responses := make([]model.InventoryResponse, len(productNames))
	for i, name := range productNames {
		qty, exists := itemMap[name]
		responses[i] = model.InventoryResponse{
			ProductName: name,
			IsInStock:   exists && qty > 0,
			Quantity:    qty,
		}
	}

	return responses, nil
}

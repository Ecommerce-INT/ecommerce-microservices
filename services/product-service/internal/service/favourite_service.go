package service

import (
	"context"

	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/repository"
)

type FavouriteService struct {
	repo *repository.FavouriteRepository
}

func NewFavouriteService(repo *repository.FavouriteRepository) *FavouriteService {
	return &FavouriteService{repo: repo}
}

func (s *FavouriteService) FindAll(ctx context.Context) ([]model.FavouriteDto, error) {
	return s.repo.FindAll(ctx)
}

func (s *FavouriteService) FindByID(ctx context.Context, id model.FavouriteID) (*model.FavouriteDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FavouriteService) Save(ctx context.Context, dto model.FavouriteDto) (*model.FavouriteDto, error) {
	return s.repo.Save(ctx, dto)
}

func (s *FavouriteService) DeleteByID(ctx context.Context, id model.FavouriteID) (bool, error) {
	return s.repo.DeleteByID(ctx, id)
}

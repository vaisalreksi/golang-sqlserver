package services

import (
	"context"
	"fmt"
	"golang-sqlserver/internal/models"
	"golang-sqlserver/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, params SearchParams) ([]models.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

type SearchParams struct {
	Keyword         string
	ProductCategory string
	Tier            string
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) Update(ctx context.Context, product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) Search(ctx context.Context, params SearchParams) ([]models.Product, error) {
	return s.repo.Search(ctx, params.Keyword, params.ProductCategory, params.Tier)
}

func validateProduct(product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if !isValidProductCategory(product.ProductCategory) {
		return fmt.Errorf("invalid product category")
	}
	return nil
}

func isValidProductCategory(category models.ProductCategory) bool {
	switch category {
	case models.Rokok, models.Obat, models.Lainnya:
		return true
	}
	return false
}

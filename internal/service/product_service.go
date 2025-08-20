package service

import (
	"context"
	"errors"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/repository"
)

type ProductService interface {
	GetProductStock(ctx context.Context, productID int64) (int, error)
	ReserveProductStock(ctx context.Context, productID int64, quantity int) (bool, error)
	ReleaseProductStock(ctx context.Context, productID int64, quantity int) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

// NewProductService creates and returns a new instance of productService.
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (p *productService) GetProductStock(ctx context.Context, productID int64) (int, error) {
	// This is a placeholder implementation.
	// In a real application, this method would interact with a database or other data source.
	// For now, we return a fixed stock value and no error.
	productDetail, err := p.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to get product stock")
		return 0, err
	}

	if productDetail == nil {
		log.Logger.Warn().Int64("productID", productID).Msg("Product not found")
		return 0, errors.New("product not found")
	}

	if productDetail.Stock < 0 {
		log.Logger.Error().Int64("productID", productID).Msg("Product stock is negative")
		return 0, errors.New("product stock is negative")
	}

	return productDetail.Stock, nil
}

func (p *productService) ReserveProductStock(ctx context.Context, productID int64, quantity int) (bool, error) {
	// This is a placeholder implementation.
	// In a real application, this method would interact with a database or other data source.
	// For now, we assume the reservation is always successful and return true with no error.
	productDetail, err := p.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to reserve product stock")
		return false, err
	}

	if productDetail == nil {
		log.Logger.Warn().Int64("productID", productID).Msg("Product not found for reservation")
		return false, errors.New("product not found")
	}

	if productDetail.Stock < quantity {
		log.Logger.Warn().Int64("productID", productID).Int("quantity", quantity).Msg("Insufficient stock for reservation")
		return false, errors.New("insufficient stock for reservation")
	}
	productDetail.Stock -= quantity
	_, err = p.productRepo.UpdateProduct(ctx, productDetail)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to update product stock after reservation")
		return false, err
	}
	return true, nil
}

func (p *productService) ReleaseProductStock(ctx context.Context, productID int64, quantity int) (bool, error) {
	productDetail, err := p.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to release product stock")
		return false, err
	}

	if productDetail == nil {
		log.Logger.Warn().Int64("productID", productID).Msg("Product not found for stock release")
		return false, errors.New("product not found")
	}

	productDetail.Stock += quantity
	_, err = p.productRepo.UpdateProduct(ctx, productDetail)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to update product stock after release")
		return false, err
	}
	return true, nil
}

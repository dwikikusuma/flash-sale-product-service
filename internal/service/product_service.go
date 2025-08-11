package service

import (
	"errors"
	"product-catalog-service/internal/repository"
)

type ProductService interface {
	GetProductStock(productID int64) (int, error)
	ReserveProductStock(productID int64, quantity int) (bool, error)
	ReleaseProductStock(productID int64, quantity int) (bool, error)
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

func (p *productService) GetProductStock(productID int64) (int, error) {
	// This is a placeholder implementation.
	// In a real application, this method would interact with a database or other data source.
	// For now, we return a fixed stock value and no error.
	productDetail, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		return 0, err
	}

	if productDetail == nil {
		return 0, errors.New("product not found")
	}

	if productDetail.Stock < 0 {
		return 0, errors.New("product stock is negative")
	}

	return productDetail.Stock, nil
}

func (p *productService) ReserveProductStock(productID int64, quantity int) (bool, error) {
	// This is a placeholder implementation.
	// In a real application, this method would interact with a database or other data source.
	// For now, we assume the reservation is always successful and return true with no error.
	productDetail, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		return false, err
	}

	if productDetail == nil {
		return false, errors.New("product not found")
	}

	if productDetail.Stock < quantity {
		return false, errors.New("insufficient stock for reservation")
	}
	productDetail.Stock -= quantity
	_, err = p.productRepo.UpdateProduct(productDetail)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *productService) ReleaseProductStock(productID int64, quantity int) (bool, error) {
	productDetail, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		return false, err
	}

	if productDetail == nil {
		return false, errors.New("product not found")
	}

	productDetail.Stock += quantity
	_, err = p.productRepo.UpdateProduct(productDetail)
	if err != nil {
		return false, err
	}
	return true, nil
}

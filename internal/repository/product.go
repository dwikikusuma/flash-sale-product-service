package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/entity"

	"gorm.io/gorm"
)

// ProductRepository defines the interface for product-related database operations.
type ProductRepository interface {
	// GetProductByID retrieves a product by its ID.
	// Parameters:
	//   - id: The ID of the product to retrieve.
	// Returns:
	//   - A pointer to the Product entity if found, or nil if not found.
	//   - An error if any issues occur during retrieval.
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)

	// CreateProduct creates a new product in the repository.
	// Parameters:
	//   - product: A pointer to the Product entity to create.
	// Returns:
	//   - A pointer to the created Product entity.
	//   - An error if any issues occur during creation.
	CreateProduct(ctx context.Context, product *entity.Product) error

	// UpdateProduct updates an existing product in the repository.
	// Parameters:
	//   - product: A pointer to the Product entity with updated data.
	// Returns:
	//   - A pointer to the updated Product entity.
	//   - An error if any issues occur during the update.
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)

	// DeleteProduct deletes a product from the repository by its ID.
	// Parameters:
	//   - id: The ID of the product to delete.
	// Returns:
	//   - An error if any issues occur during deletion.
	DeleteProduct(ctx context.Context, id int64) error
}

// productRepository is a concrete implementation of the ProductRepository interface.
type productRepository struct {
	cache CacheRepository
	db    *gorm.DB
}

// NewProductRepository creates a new instance of productRepository.
// Returns:
//   - A ProductRepository instance.
func NewProductRepository(cacheRepo CacheRepository, db *gorm.DB) ProductRepository {
	return &productRepository{
		cache: cacheRepo,
		db:    db,
	}
}

// GetProductByID retrieves a product by its ID from the in-memory store.
// Parameters:
//   - id: The ID of the product to retrieve.
//
// Returns:
//   - A pointer to the Product entity if found, or nil if not found.
//   - An error if any issues occur during retrieval.
func (r *productRepository) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	key := fmt.Sprintf("product:%d", id)
	productCache, err := r.cache.Get(ctx, key)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to get product from cache")
		return nil, fmt.Errorf("failed to get product from cache: %w", err)
	}

	if productCache != "" {
		var productFromCache entity.Product
		if err := json.Unmarshal([]byte(productCache), &productFromCache); err != nil {
			log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to unmarshal product from cache")
			return nil, fmt.Errorf("failed to unmarshal product from cache: %w", err)
		}
		return &productFromCache, nil
	}

	var product *entity.Product
	err = r.db.Table("products").WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to get product from database")
		return nil, errors.New("failed to get product from database")
	}

	// Cache the product for future requests
	if err := r.cache.Set(ctx, key, product); err != nil {
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to set product in cache")
		return nil, fmt.Errorf("failed to set product in cache: %w", err)
	}

	return product, nil

}

// CreateProduct adds a new product to the in-memory store.
// Parameters:
//   - product: A pointer to the Product entity to create.
//
// Returns:
//   - A pointer to the created Product entity.
//   - An error if any issues occur during creation.
func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	err := r.db.Table("products").WithContext(ctx).Create(product).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create product in database")
		return errors.New("failed to create product in database: %w")
	}
	return nil
}

// UpdateProduct updates an existing product in the in-memory store.
// Parameters:
//   - product: A pointer to the Product entity with updated data.
//
// Returns:
//   - A pointer to the updated Product entity.
//   - An error if any issues occur during the update.
func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	err := r.db.Table("products").WithContext(ctx).Save(product).Error
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", product.ID).Msg("Failed to update product in database")
		return nil, errors.New("failed to update product in database")
	}

	err = r.cache.Set(ctx, fmt.Sprintf("product:%d", product.ID), product)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", product.ID).Msg("Failed to update product in cache")
		return nil, fmt.Errorf("failed to update product in cache: %w", err)
	}

	updatedProduct, err := r.GetProductByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated product: %w", err)
	}
	return updatedProduct, nil
}

// DeleteProduct removes a product from the in-memory store by its ID.
// Parameters:
//   - id: The ID of the product to delete.
//
// Returns:
//   - An error if any issues occur during deletion.
func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	product, err := r.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve product before deletion: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product with ID %d not found", id)
	}

	err = r.db.Table("products").WithContext(ctx).Delete(&entity.Product{}, id).Error
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to delete product from database")
		return fmt.Errorf("failed to delete product from database: %w", err)
	}

	err = r.cache.Delete(ctx, fmt.Sprintf("product:%d", id))
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to delete product from cache")
		return fmt.Errorf("failed to delete product from cache: %w", err)
	}

	return nil
}

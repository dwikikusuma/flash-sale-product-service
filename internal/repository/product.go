package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/entity"
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
	CreateProduct(product *entity.Product) (*entity.Product, error)

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
}

// NewProductRepository creates a new instance of productRepository.
// Returns:
//   - A ProductRepository instance.
func NewProductRepository(cacheRepo CacheRepository) ProductRepository {
	return &productRepository{
		cache: cacheRepo,
	}
}

// products is an in-memory store simulating a database for products.
var products = map[int64]*entity.Product{
	1: {
		ID:          1,
		Name:        "Product A",
		Description: "Description of Product A",
		Price:       100.0,
		Stock:       50,
	},
	2: {
		ID:          2,
		Name:        "Product B",
		Description: "Description of Product B",
		Price:       200.0,
		Stock:       30,
	},
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

	product, ok := products[id]
	if !ok {
		log.Logger.Warn().Int64("productID", id).Msg("Product not found")
		return nil, nil // Product not found, return nil
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
func (r *productRepository) CreateProduct(product *entity.Product) (*entity.Product, error) {
	product.ID = int64(len(products) + 1) // Simulating an auto-generated ID
	products[product.ID] = product
	return product, nil
}

// UpdateProduct updates an existing product in the in-memory store.
// Parameters:
//   - product: A pointer to the Product entity with updated data.
//
// Returns:
//   - A pointer to the updated Product entity.
//   - An error if any issues occur during the update.
func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	key := fmt.Sprintf("product:%d", product.ID)
	err := r.cache.Set(ctx, key, product)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", product.ID).Msg("Failed to update product in cache")
		return nil, fmt.Errorf("failed to update product in cache: %w", err)
	}
	products[product.ID] = product
	return product, nil
}

// DeleteProduct removes a product from the in-memory store by its ID.
// Parameters:
//   - id: The ID of the product to delete.
//
// Returns:
//   - An error if any issues occur during deletion.
func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	delete(products, id)
	err := r.cache.Delete(ctx, fmt.Sprintf("product:%d", id))
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", id).Msg("Failed to delete product from cache")
		return fmt.Errorf("failed to delete product from cache: %w", err)
	}
	return nil
}

package repository

import "product-catalog-service/internal/entity"

// ProductRepository defines the interface for product-related database operations.
type ProductRepository interface {
	// GetProductByID retrieves a product by its ID.
	// Parameters:
	//   - id: The ID of the product to retrieve.
	// Returns:
	//   - A pointer to the Product entity if found, or nil if not found.
	//   - An error if any issues occur during retrieval.
	GetProductByID(id int64) (*entity.Product, error)

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
	UpdateProduct(product *entity.Product) (*entity.Product, error)

	// DeleteProduct deletes a product from the repository by its ID.
	// Parameters:
	//   - id: The ID of the product to delete.
	// Returns:
	//   - An error if any issues occur during deletion.
	DeleteProduct(id int64) error
}

// productRepository is a concrete implementation of the ProductRepository interface.
type productRepository struct {
}

// NewProductRepository creates a new instance of productRepository.
// Returns:
//   - A ProductRepository instance.
func NewProductRepository() ProductRepository {
	return &productRepository{}
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
func (r *productRepository) GetProductByID(id int64) (*entity.Product, error) {
	product, ok := products[id]
	if !ok {
		return nil, nil // Product not found
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
func (r *productRepository) UpdateProduct(product *entity.Product) (*entity.Product, error) {
	products[product.ID] = product
	return product, nil
}

// DeleteProduct removes a product from the in-memory store by its ID.
// Parameters:
//   - id: The ID of the product to delete.
//
// Returns:
//   - An error if any issues occur during deletion.
func (r *productRepository) DeleteProduct(id int64) error {
	delete(products, id)
	return nil
}

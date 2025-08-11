package entity

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// ProductStock represents the stock information for a product.
type StockReservation struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

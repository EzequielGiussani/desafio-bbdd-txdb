package internal

import "errors"

var (
	ServiceProductInternalServerError = errors.New("Product: internal server error")
)

// ServiceProduct is the interface that wraps the basic Product methods.
type ServiceProduct interface {
	// FindAll returns all products.
	FindAll() (p []Product, err error)
	// Save saves a product.
	Save(p *Product) (err error)
	GetTopFiveProducts() (p []ProductDescTotal, err error)
}

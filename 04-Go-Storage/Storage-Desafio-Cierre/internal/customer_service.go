package internal

import "errors"

var (
	ServiceCustomerInternalServerError = errors.New("Customer: internal server error")
)

// ServiceCustomer is the interface that wraps the basic methods that a customer service should implement.
type ServiceCustomer interface {
	// FindAll returns all customers
	FindAll() (c []Customer, err error)
	// Save saves a customer
	Save(c *Customer) (err error)
	GetTotalAmountInEachCondition() (ccTotal []CcTotal, err error)
	GetTopFiveActiveCustomers() (customers []CustomerNameAmount, err error)
}

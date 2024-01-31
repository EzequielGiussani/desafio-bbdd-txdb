package service

import "app/internal"

// NewCustomersDefault creates new default service for customer entity.
func NewCustomersDefault(rp internal.RepositoryCustomer) *CustomersDefault {
	return &CustomersDefault{rp}
}

// CustomersDefault is the default service implementation for customer entity.
type CustomersDefault struct {
	// rp is the repository for customer entity.
	rp internal.RepositoryCustomer
}

// FindAll returns all customers.
func (s *CustomersDefault) FindAll() (c []internal.Customer, err error) {
	c, err = s.rp.FindAll()
	return
}

// Save saves the customer.
func (s *CustomersDefault) Save(c *internal.Customer) (err error) {
	err = s.rp.Save(c)
	return
}

func (s *CustomersDefault) GetTotalAmountInEachCondition() (ccTotal []internal.CcTotal, err error) {
	ccTotal, err = s.rp.GetTotalAmountInEachCondition()
	if err != nil {
		switch err {
		default:
			return nil, internal.ServiceCustomerInternalServerError
		}
	}
	return
}

func (s *CustomersDefault) GetTopFiveActiveCustomers() (customers []internal.CustomerNameAmount, err error) {
	customers, err = s.rp.GetTopFiveActiveCustomers()
	if err != nil {
		switch err {
		default:
			return nil, internal.ServiceCustomerInternalServerError
		}
	}
	return
}

package internal

// RepositoryCustomer is the interface that wraps the basic methods that a customer repository should implement.
type RepositoryCustomer interface {
	// FindAll returns all customers saved in the database.
	FindAll() (c []Customer, err error)
	// Save saves a customer into the database.
	Save(c *Customer) (err error)

	GetTotalAmountInEachCondition() (ccTotal []CcTotal, err error)

	GetTopFiveActiveCustomers() (customers []CustomerNameAmount, err error)
}

type CcTotal struct {
	Condition string
	Total     float64
}

type CustomerNameAmount struct {
	FirstName string
	LastName  string
	Amount    float64
}

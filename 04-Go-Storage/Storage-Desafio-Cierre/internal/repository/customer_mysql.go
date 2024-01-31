package repository

import (
	"database/sql"

	"app/internal"
)

// NewCustomersMySQL creates new mysql repository for customer entity.
func NewCustomersMySQL(db *sql.DB) *CustomersMySQL {
	return &CustomersMySQL{db}
}

// CustomersMySQL is the MySQL repository implementation for customer entity.
type CustomersMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all customers from the database.
func (r *CustomersMySQL) FindAll() (c []internal.Customer, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `first_name`, `last_name`, `condition` FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var cs internal.Customer
		// scan the row into the customer
		err := rows.Scan(&cs.Id, &cs.FirstName, &cs.LastName, &cs.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cs)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the customer into the database.
func (r *CustomersMySQL) Save(c *internal.Customer) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
		(*c).FirstName, (*c).LastName, (*c).Condition,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*c).Id = int(id)

	return
}

func (r *CustomersMySQL) GetTotalAmountInEachCondition() (ccTotal []internal.CcTotal, err error) {
	rows, err := r.db.Query("select c.`condition` , ROUND(sum(i.total), 2) from customers c inner join invoices i on i.customer_id = c.id group by c.`condition` ")

	if err != nil {
		return []internal.CcTotal{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var ccTot internal.CcTotal
		var condition int
		err = rows.Scan(&condition, &ccTot.Total)

		if condition == 1 {
			ccTot.Condition = "Activo"
		} else {
			ccTot.Condition = "Inactivo"
		}

		ccTotal = append(ccTotal, ccTot)
	}

	err = rows.Err()

	if err != nil {
		return []internal.CcTotal{}, err
	}

	return

}

func (r *CustomersMySQL) GetTopFiveActiveCustomers() (customers []internal.CustomerNameAmount, err error) {
	rows, err := r.db.Query("select c.first_name, c.last_name, ROUND(SUM(i.total), 2) as total from customers c inner join invoices i on c.id = i.customer_id where c.`condition` = 1 group by c.id order by total desc limit 5")

	if err != nil {
		return []internal.CustomerNameAmount{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var customer internal.CustomerNameAmount
		err = rows.Scan(&customer.FirstName, &customer.LastName, &customer.Amount)

		customers = append(customers, customer)
	}

	err = rows.Err()

	if err != nil {
		return []internal.CustomerNameAmount{}, err
	}

	return
}

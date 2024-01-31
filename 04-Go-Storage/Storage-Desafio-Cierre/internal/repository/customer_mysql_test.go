package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "root",
		Addr:      "localhost:3306",
		Net:       "tcp",
		DBName:    "test_db_fantasy_products",
		ParseTime: true,
	}

	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestCustomer_GetTotalAmountInEachCondition(t *testing.T) {

	t.Run("should return a slice of CcTotal", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdb", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (1, 'John', 'Doe', 0)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (2, 'Juan', 'Doe', 0)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (3, 'Lucas', 'Tres', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (4, 'Eze', 'Cuatro', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (1, NOW(), 1, 100)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (2, NOW(), 2, 100)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (3, NOW(), 3, 200)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (4, NOW(), 4, 200)")
			require.NoError(t, err)

		}(db)

		rp := repository.NewCustomersMySQL(db)

		//act

		ccTotal, err := rp.GetTotalAmountInEachCondition()

		require.NoError(t, err)

		//assert

		expectedCc := []internal.CcTotal{
			{
				Condition: "Inactivo",
				Total:     200,
			},
			{
				Condition: "Activo",
				Total:     400,
			},
		}

		require.Equal(t, expectedCc, ccTotal)
	})

	t.Run("should return empty", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdb", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewCustomersMySQL(db)

		//act

		ccTotal, err := rp.GetTotalAmountInEachCondition()

		require.NoError(t, err)

		//assert

		expectedCc := []internal.CcTotal(nil)

		require.Equal(t, expectedCc, ccTotal)
	})
}

func TestCustomer_GetTopFiveActiveCustomers(t *testing.T) {

	t.Run("should return a slice of CcTotal", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdb", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (1, 'John', 'Doe', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (2, 'Juan', 'Doe', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (3, 'Lucas', 'Tres', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `customers` (`id`, `first_name`, `last_name`, `condition`) VALUES (4, 'Eze', 'Cuatro', 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (1, NOW(), 1, 400)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (2, NOW(), 2, 300)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (3, NOW(), 3, 200)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `invoices` (`id`, `datetime`, `customer_id`, `total`) VALUES (4, NOW(), 4, 100)")
			require.NoError(t, err)

		}(db)

		rp := repository.NewCustomersMySQL(db)

		//act

		customers, err := rp.GetTopFiveActiveCustomers()

		require.NoError(t, err)

		//assert

		expectedCustomers := []internal.CustomerNameAmount{
			{
				FirstName: "John",
				LastName:  "Doe",
				Amount:    400,
			},
			{
				FirstName: "Juan",
				LastName:  "Doe",
				Amount:    300,
			},
			{
				FirstName: "Lucas",
				LastName:  "Tres",
				Amount:    200,
			},
			{
				FirstName: "Eze",
				LastName:  "Cuatro",
				Amount:    100,
			},
		}

		require.Equal(t, expectedCustomers, customers)
	})

	t.Run("should return empty", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdb", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewCustomersMySQL(db)

		//act

		ccTotal, err := rp.GetTopFiveActiveCustomers()

		require.NoError(t, err)

		//assert

		expectedCustomers := []internal.CustomerNameAmount(nil)

		require.Equal(t, expectedCustomers, ccTotal)
	})
}

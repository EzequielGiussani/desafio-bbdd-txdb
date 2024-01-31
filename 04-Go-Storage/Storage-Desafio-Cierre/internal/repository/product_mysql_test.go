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

	txdb.Register("txdbProducts", "mysql", cfg.FormatDSN())
}

func TestCustomer_GetTopFiveProducts(t *testing.T) {

	t.Run("should return a slice of CcTotal", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdbProducts", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `products` (`id`, `description`, `price`) VALUES (1, 'description1', 1000)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `products` (`id`, `description`, `price`) VALUES (2, 'description2', 1000)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `products` (`id`, `description`, `price`) VALUES (3, 'description3', 1000)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `products` (`id`, `description`, `price`) VALUES (4, 'description4', 1000)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (1, 20, null, 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (2, 20, null, 1)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (3, 20, null, 2)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (4, 20, null, 3)")
			require.NoError(t, err)

			_, err = db.Exec("INSERT INTO `sales` (`id`, `quantity`, `invoice_id`, `product_id`) VALUES (5, 20, null, 1)")
			require.NoError(t, err)

		}(db)

		rp := repository.NewProductsMySQL(db)

		//act

		products, err := rp.GetTopFiveProducts()

		require.NoError(t, err)

		//assert

		expectedProducts := []internal.ProductDescTotal{
			{
				Description: "description1",
				Total:       3,
			},
			{
				Description: "description2",
				Total:       1,
			},
			{
				Description: "description3",
				Total:       1,
			},
		}

		require.Equal(t, expectedProducts, products)
	})

	t.Run("should return empty", func(t *testing.T) {

		//Arrange
		db, err := sql.Open("txdbProducts", "test_db_fantasy_products")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewProductsMySQL(db)

		//act

		products, err := rp.GetTopFiveProducts()

		require.NoError(t, err)

		//assert

		expectedProducts := []internal.ProductDescTotal(nil)

		require.Equal(t, expectedProducts, products)
	})
}

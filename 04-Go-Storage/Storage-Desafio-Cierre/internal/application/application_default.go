package application

import (
	"app/internal"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

// ConfigApplicationDefault is the configuration for NewApplicationDefault.
type ConfigApplicationDefault struct {
	// Db is the database configuration.
	Db *mysql.Config
	// Addr is the server address.
	Addr string
}

// NewApplicationDefault creates a new ApplicationDefault.
func NewApplicationDefault(config *ConfigApplicationDefault) *ApplicationDefault {
	// default values
	defaultCfg := &ConfigApplicationDefault{
		Db:   nil,
		Addr: ":8080",
	}
	if config != nil {
		if config.Db != nil {
			defaultCfg.Db = config.Db
		}
		if config.Addr != "" {
			defaultCfg.Addr = config.Addr
		}
	}

	return &ApplicationDefault{
		cfgDb:   defaultCfg.Db,
		cfgAddr: defaultCfg.Addr,
	}
}

// ApplicationDefault is an implementation of the Application interface.
type ApplicationDefault struct {
	// cfgDb is the database configuration.
	cfgDb *mysql.Config
	// cfgAddr is the server address.
	cfgAddr string
	// db is the database connection.
	db *sql.DB
	// router is the chi router.
	router *chi.Mux
}

// SetUp sets up the application.
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - db: init
	a.db, err = sql.Open("mysql", a.cfgDb.FormatDSN())
	if err != nil {
		return
	}
	// - db: ping
	err = a.db.Ping()
	if err != nil {
		return
	}
	// - repository
	rpCustomer := repository.NewCustomersMySQL(a.db)
	rpProduct := repository.NewProductsMySQL(a.db)
	rpInvoice := repository.NewInvoicesMySQL(a.db)
	rpSale := repository.NewSalesMySQL(a.db)
	// - service
	svCustomer := service.NewCustomersDefault(rpCustomer)
	svProduct := service.NewProductsDefault(rpProduct)
	svInvoice := service.NewInvoicesDefault(rpInvoice)
	svSale := service.NewSalesDefault(rpSale)
	// - handler
	hdCustomer := handler.NewCustomersDefault(svCustomer)
	hdProduct := handler.NewProductsDefault(svProduct)
	hdInvoice := handler.NewInvoicesDefault(svInvoice)
	hdSale := handler.NewSalesDefault(svSale)

	if err = loadDatabase(rpCustomer, rpInvoice, rpProduct, rpSale); err != nil {
		return
	}

	// routes
	// - router
	a.router = chi.NewRouter()
	// - middlewares
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	// - endpoints
	a.router.Route("/customers", func(r chi.Router) {
		// - GET /customers
		r.Get("/", hdCustomer.GetAll())
		// - POST /customers
		r.Post("/", hdCustomer.Create())
		r.Get("/GetTotalAmountInEachCondition", hdCustomer.GetTotalAmountInEachCondition())
		r.Get("/GetTopFiveActiveCustomers", hdCustomer.GetTopFiveActiveCustomers())
	})
	a.router.Route("/products", func(r chi.Router) {
		// - GET /products
		r.Get("/", hdProduct.GetAll())
		// - POST /products
		r.Post("/", hdProduct.Create())
		r.Get("/GetTopFiveProducts", hdProduct.GetTopFiveProducts())
	})
	a.router.Route("/invoices", func(r chi.Router) {
		// - GET /invoices
		r.Get("/", hdInvoice.GetAll())
		// - POST /invoices
		r.Post("/", hdInvoice.Create())
		// - UPDATE /invoices
		r.Patch("/", hdInvoice.UpdateTotal())
	})
	a.router.Route("/sales", func(r chi.Router) {
		// - GET /sales
		r.Get("/", hdSale.GetAll())
		// - POST /sales
		r.Post("/", hdSale.Create())
	})

	return
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	defer a.db.Close()

	err = http.ListenAndServe(a.cfgAddr, a.router)
	return
}

func loadDatabase(rpCustomer *repository.CustomersMySQL, rpInvoice *repository.InvoicesMySQL, rpProduct *repository.ProductsMySQL, rpSale *repository.SalesMySQL) (err error) {
	//check customers
	customers, err := rpCustomer.FindAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(customers) <= 0 {
		//load customers
		if err = loadCustomers(rpCustomer); err != nil {
			fmt.Println(err)
			return
		}
	}

	//check invoices
	invoices, err := rpInvoice.FindAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(invoices) <= 0 {
		//load invoices
		loadInvoices(rpInvoice)
	}

	//check products
	products, err := rpProduct.FindAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(products) <= 0 {
		//load products
		loadProducts(rpProduct)
	}

	//check sales
	sales, err := rpSale.FindAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(sales) <= 0 {
		//load sales
		loadSales(rpSale)
	}

	return nil
}

func loadCustomers(rpCustomer *repository.CustomersMySQL) (err error) {

	file, err := os.Open("./docs/db/json/customers.json")

	if err != nil {
		return
	}

	defer file.Close()

	var customers []handler.CustomerJSON

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&customers); err != nil {
		return
	}

	for _, customerJSON := range customers {
		customer := internal.Customer{
			Id: customerJSON.Id,
			CustomerAttributes: internal.CustomerAttributes{
				FirstName: customerJSON.FirstName,
				LastName:  customerJSON.LastName,
				Condition: customerJSON.Condition,
			},
		}
		err = rpCustomer.Save(&customer)

		if err != nil {
			return
		}
	}

	return nil

}

func loadInvoices(rpInvoice *repository.InvoicesMySQL) (err error) {

	file, err := os.Open("./docs/db/json/invoices.json")

	if err != nil {
		return
	}

	defer file.Close()

	var invoices []handler.InvoiceJSON

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&invoices); err != nil {
		return
	}

	for _, invoicesJSON := range invoices {
		invoice := internal.Invoice{
			Id: invoicesJSON.Id,
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   invoicesJSON.Datetime,
				CustomerId: invoicesJSON.CustomerId,
				Total:      invoicesJSON.Total,
			},
		}
		err = rpInvoice.Save(&invoice)

		if err != nil {
			return
		}
	}

	return nil

}

func loadProducts(rpProduct *repository.ProductsMySQL) (err error) {

	file, err := os.Open("./docs/db/json/products.json")

	if err != nil {
		return
	}

	defer file.Close()

	var products []handler.ProductJSON

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&products); err != nil {
		return
	}

	for _, productsJSON := range products {
		product := internal.Product{
			Id: productsJSON.Id,
			ProductAttributes: internal.ProductAttributes{
				Price:       productsJSON.Price,
				Description: productsJSON.Description,
			},
		}
		err = rpProduct.Save(&product)

		if err != nil {
			return
		}
	}

	return nil

}

func loadSales(rpSales *repository.SalesMySQL) (err error) {

	file, err := os.Open("./docs/db/json/sales.json")

	if err != nil {
		return
	}

	defer file.Close()

	var sales []handler.SaleJSON

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&sales); err != nil {
		return
	}

	for _, salesJSON := range sales {
		sale := internal.Sale{
			Id: salesJSON.Id,
			SaleAttributes: internal.SaleAttributes{
				InvoiceId: salesJSON.InvoiceId,
				ProductId: salesJSON.ProductId,
				Quantity:  salesJSON.Quantity,
			},
		}
		err = rpSales.Save(&sale)

		if err != nil {
			return
		}
	}

	return nil

}

package api

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/NikhilSharmaWe/market/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	app *Application
)

func TestMain(m *testing.M) {
	var err error
	db, err = createDB()

	if err != nil {
		log.Fatal(err)
	}

	app = NewApplication(db)

	result := m.Run()

	os.Exit(result)
}

func createDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("user=miyamoto dbname=test_market password=1234 sslmode=disable"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupTestingEnvironment(db *gorm.DB) error {
	log.Println("seeding")

	seedQuery := ""

	for _, file := range []string{"../db.sql"} {
		filePath, err := filepath.Abs(file)
		if err != nil {
			return err
		}

		seedFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(seedFile)
		if err != nil {
			return err
		}

		seedQuery += string(content) + "\n"

		seedFile.Close()
	}

	return db.Exec(seedQuery).Error
}

func cleanupTestingEnvironment(db *gorm.DB) error {
	return db.Exec("drop table admins, inventory, orders, order_per_product, products, users;").Error
}

func TestProductServices(t *testing.T) {
	if err := setupTestingEnvironment(db); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanupTestingEnvironment(db); err != nil {
			t.Fatal(err)
		}
	}()

	exampleProductReq := models.CreateProductRequest{
		ProductName:     "Handle",
		Category:        "Wood",
		Price:           12,
		InitialQuantity: 20,
	}

	// testing Create method
	_, err := app.ProductService.Get(nil)
	assert.Equal(t, err, models.ErrMatchingRecordNotFound)

	err = app.ProductService.Create(&exampleProductReq)
	assert.Nil(t, err)

	// testing Get method
	resp, err := app.ProductService.Get(nil)
	assert.Nil(t, err)
	assert.Equal(t, map[string]models.ProductResponse{exampleProductReq.ProductName: {
		Category: exampleProductReq.Category,
		Price:    exampleProductReq.Price,
	}}, resp)

	inventory, err := app.InventoryStore.GetOne(map[string]interface{}{"product_name": exampleProductReq.ProductName})
	assert.Nil(t, err)
	assert.Equal(t, models.InventoryDBModel{
		ProductName:     exampleProductReq.ProductName,
		QuantityInStock: exampleProductReq.InitialQuantity,
	}, *inventory)

	// testing Update method
	err = app.ProductService.Update(&models.UpdateProductRequest{
		ProductName: "non_existent_product_name",
		Category:    "Steel",
		Price:       12,
	})
	assert.Equal(t, models.ErrProductNotFound, err)

	err = app.ProductService.Update(&models.UpdateProductRequest{
		ProductName: exampleProductReq.ProductName,
		Category:    "Steel",
		Price:       12,
	})
	assert.Nil(t, err)

	p, err := app.ProductStore.GetOne(map[string]interface{}{"product_name": exampleProductReq.ProductName})
	assert.Nil(t, err)
	assert.Equal(t, models.ProductDBModel{
		ProductName: exampleProductReq.ProductName,
		Category:    "Steel",
		Price:       12,
	}, *p)

	// testing Delete method
	err = app.ProductService.Delete(&models.DeleteProductRequest{
		ProductName: exampleProductReq.ProductName,
	})
	assert.Nil(t, err)

	exists, err := app.ProductStore.IsExists(map[string]interface{}{"product_name": exampleProductReq.ProductName})
	assert.Nil(t, err)
	assert.Equal(t, false, exists)
}

func TestOrderService(t *testing.T) {
	if err := setupTestingEnvironment(db); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanupTestingEnvironment(db); err != nil {
			t.Fatal(err)
		}
	}()

	exampleUser1 := models.UserDBModel{
		Username: "Nikhil",
		Email:    "a@a.com",
		Password: []byte("123"),
	}

	exampleUser2 := models.UserDBModel{
		Username: "Rewak",
		Email:    "r@r.com",
		Password: []byte("321"),
	}

	err := app.UsersStore.Create(exampleUser1)
	assert.Nil(t, err)

	err = app.UsersStore.Create(exampleUser2)
	assert.Nil(t, err)

	exampleProductOneReq := models.CreateProductRequest{
		ProductName:     "Handle",
		Category:        "Wood",
		Price:           12,
		InitialQuantity: 20,
	}

	exampleProductTwoReq := models.CreateProductRequest{
		ProductName:     "Wheel",
		Category:        "Steel",
		Price:           20,
		InitialQuantity: 20,
	}

	exampleOrderReqOne := models.CreateOrderRequest{
		Username:            "Nikhil",
		ProductsAndQuantity: map[string]int{"Handle": 2, "Wheel": 4},
	}

	exampleOrderReqTwo := models.CreateOrderRequest{
		Username:            "Rewak",
		ProductsAndQuantity: map[string]int{"Handle": 3},
	}

	exampleOrderReqThree := models.CreateOrderRequest{
		Username:            "Rewak",
		ProductsAndQuantity: map[string]int{"Wheel": 3},
	}

	_, err = app.OrderService.Create(&exampleOrderReqOne)
	assert.Equal(t, models.ErrProductNotFound, err)

	err = app.ProductService.Create(&exampleProductOneReq)
	assert.Nil(t, err)

	err = app.ProductService.Create(&exampleProductTwoReq)
	assert.Nil(t, err)

	// testing Create method
	order1, err := app.OrderService.Create(&exampleOrderReqOne)
	assert.Nil(t, err)

	_, err = app.OrderStore.GetOne(map[string]interface{}{"order_id": order1.OrderID})
	assert.Nil(t, err)
	assert.Equal(t, "Nikhil", order1.Username)

	expectedTotalCost := exampleProductOneReq.Price*exampleOrderReqOne.ProductsAndQuantity[exampleProductOneReq.ProductName] + exampleProductTwoReq.Price*exampleOrderReqOne.ProductsAndQuantity[exampleProductTwoReq.ProductName]
	assert.Equal(t, expectedTotalCost, order1.TotalCost)

	order2, err := app.OrderService.Create(&exampleOrderReqTwo)
	assert.Nil(t, err)

	_, err = app.OrderStore.GetOne(map[string]interface{}{"order_id": order2.OrderID})
	assert.Nil(t, err)

	order3, err := app.OrderService.Create(&exampleOrderReqThree)
	assert.Nil(t, err)

	_, err = app.OrderStore.GetOne(map[string]interface{}{"order_id": order3.OrderID})
	assert.Nil(t, err)

	// testing GetOrdersByUsername method
	_, err = app.OrderService.GetOrdersByUsername("non_existent_user_name")
	assert.Equal(t, models.ErrMatchingRecordNotFound, err)

	resp, err := app.OrderService.GetOrdersByUsername(exampleOrderReqOne.Username)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp))
	for _, orderResponse := range resp {
		if orderResponse.OrderID != order1.OrderID {
			t.Fatal("unexpected response: expected order id:", order1.OrderID)
		}
	}

	// testing GetOrdersForAdmin method
	resp, err = app.OrderService.GetOrdersForAdmin(&models.GetOrdersForAdminRequest{
		Username:    "Rewak",
		ProductName: "Wheel",
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp))
	for _, orderResponse := range resp {
		if orderResponse.OrderID != order3.OrderID {
			t.Fatal("unexpected response: expected order id:", order3.OrderID)
		}
	}
}

func TestInventoryService(t *testing.T) {
	if err := setupTestingEnvironment(db); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanupTestingEnvironment(db); err != nil {
			t.Fatal(err)
		}
	}()

	exampleProductReq := models.CreateProductRequest{
		ProductName:     "Handle",
		Category:        "Wood",
		Price:           12,
		InitialQuantity: 20,
	}

	err := app.ProductService.Create(&exampleProductReq)
	assert.Nil(t, err)

	// testing Update method
	err = app.InventoryService.Update(&models.UpdateInventoryRequest{
		ProductName: exampleProductReq.ProductName,
		Operation:   "invalid_operation",
		Quantity:    10,
	})
	assert.Equal(t, models.ErrInvalidOperaton, err)

	err = app.InventoryService.Update(&models.UpdateInventoryRequest{
		ProductName: exampleProductReq.ProductName,
		Operation:   "remove",
		Quantity:    100,
	})
	assert.Equal(t, models.ErrInvalidQuantity, err)

	err = app.InventoryService.Update(&models.UpdateInventoryRequest{
		ProductName: exampleProductReq.ProductName,
		Operation:   "remove",
		Quantity:    10,
	})
	assert.Nil(t, err)

	inventory, err := app.InventoryStore.GetOne(map[string]interface{}{"product_name": exampleProductReq.ProductName})
	assert.Nil(t, err)
	assert.Equal(t, 10, inventory.QuantityInStock)
}

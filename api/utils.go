package api

import (
	"net/http"
	"strconv"

	"github.com/NikhilSharmaWe/market/models"
	"github.com/NikhilSharmaWe/market/store"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Application struct {
	CookieStore *sessions.CookieStore
	ProductService
	InventoryService
	OrderService
	store.UsersStore
	store.InventoryStore
	store.OrderStore
	store.ProductStore
	store.OrderPerProductStore
	store.AdminsStore
}

func NewApplication(db *gorm.DB, sessionSecretKey string) *Application {
	userStore := store.NewUsersStore(db)
	inventoryStore := store.NewInventoryStore(db)
	orderStore := store.NewOrdersStore(db)
	orderPerProductStore := store.NewOrderPerProductStore(db)
	productStore := store.NewProductsStore(db)
	adminStore := store.NewAdminsStore(db)

	productService := NewProductService(productStore, inventoryStore)
	inventoryService := NewInventoryService(productStore, inventoryStore)
	orderService := NewOrderService(orderStore, orderPerProductStore, productStore, inventoryStore)

	return &Application{
		CookieStore:      sessions.NewCookieStore([]byte(sessionSecretKey)),
		UsersStore:       userStore,
		InventoryStore:   inventoryStore,
		OrderStore:       orderStore,
		ProductStore:     productStore,
		ProductService:   productService,
		InventoryService: inventoryService,
		OrderService:     orderService,
		AdminsStore:      adminStore,
	}
}

func setSession(c echo.Context) error {
	session := c.Get("session").(*sessions.Session)
	session.ID = uuid.NewV4().String()
	session.Values["username"] = c.FormValue("username")
	session.Values["authenticated"] = true
	return session.Save(c.Request(), c.Response())
}

func clearSessionHandler(c echo.Context) error {
	session := c.Get("session").(*sessions.Session)
	session.Options.MaxAge = -1
	return session.Save(c.Request(), c.Response())
}

func (app *Application) alreadyLoggedIn(c echo.Context) bool {
	session := c.Get("session").(*sessions.Session)

	username, ok := session.Values["username"].(string)
	if !ok {
		return false
	}

	if exists, err := app.UsersStore.IsExists(map[string]interface{}{"username": username}); err != nil || !exists {
		return false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if ok && authenticated {
		return true
	}

	return false
}

func (app *Application) isAdmin(c echo.Context) bool {
	session := c.Get("session").(*sessions.Session)

	username, ok := session.Values["username"].(string)
	if !ok {
		return false
	}

	if exists, err := app.AdminsStore.IsExists(map[string]interface{}{"username": username}); err != nil || !exists {
		return false
	}

	return true
}

func userFromContext(c echo.Context) (*models.UserDBModel, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	return &models.UserDBModel{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: bs,
	}, nil
}

func productFromContext(c echo.Context) (*models.ProductDBModel, error) {
	var (
		price int
		err   error
	)
	priceStr := c.FormValue("price")

	if priceStr == "" {
		price = 0
	} else {
		price, err = strconv.Atoi(priceStr)
		if err != nil {
			return nil, err
		}
	}

	return &models.ProductDBModel{
		ProductName: c.FormValue("product_name"),
		Category:    c.FormValue("category"),
		Price:       price,
	}, nil
}

func createOrderReqFromContext(c echo.Context) (*models.CreateOrderRequest, error) {
	session := c.Get("session").(*sessions.Session)
	username := session.Values["username"].(string)
	productName := c.FormValue("product_name")
	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		return nil, err
	}

	return &models.CreateOrderRequest{
		Username:            username,
		ProductsAndQuantity: map[string]int{productName: quantity},
	}, nil
}

func (app *Application) addProduct(c echo.Context) error {
	product, err := productFromContext(c)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	initialQuantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	req := &models.CreateProductRequest{
		ProductName:     product.ProductName,
		Category:        product.Category,
		Price:           product.Price,
		InitialQuantity: initialQuantity,
	}

	if err := app.ProductService.Create(req); err != nil {
		if err == models.ErrProductAlreadyExists {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) removeProduct(c echo.Context) error {
	productName := c.FormValue("product_name")

	if err := app.ProductService.Delete(&models.DeleteProductRequest{
		ProductName: productName,
	}); err != nil {
		if err == models.ErrProductNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) updateProduct(c echo.Context) error {
	product, err := productFromContext(c)
	if err != nil {
		return err
	}

	req := models.UpdateProductRequest(*product)
	if err := app.ProductService.Update(&req); err != nil {
		if err == models.ErrProductNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		c.Logger().Error(err)
		return err
	}

	return nil
}

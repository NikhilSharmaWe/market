package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NikhilSharmaWe/market/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"

	echo "github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

func (app *Application) Router() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(app.createSessionMiddleware)

	admin := e.Group("/admin")
	admin.Use(app.IfNotLogined, app.IfNotAdmin)

	user := e.Group("/user")
	user.Use(app.IfNotLogined)

	e.Static("/assets", "./public")

	e.GET("/", ServeHTML("./public/login/index.html"), app.IfAlreadyLogined)
	e.POST("/", app.HandleSignIn)

	e.GET("/signup", ServeHTML("./public/signup/index.html"), app.IfAlreadyLogined)
	e.POST("/signup", app.HandleSignUp)

	e.GET("/logout", app.HandleLogout)

	// user apis
	user.GET("/home", ServeHTML("./public/home/index.html"))

	user.GET("/order", ServeHTML("./public/order/index.html"))
	user.POST("/order", app.HandleUserOrder)

	user.GET("/myorders", app.HandleMyOrders)

	user.GET("/products", ServeHTML("./public/product/index.html"))
	user.POST("/products", app.HandleGetProducts)

	// admin apis
	admin.GET("/home", ServeHTML("./public/admin_home/index.html"))

	admin.GET("/product", ServeHTML("./public/admin_product/index.html"))

	admin.GET("/product/:operation", app.HandleAdminProductFiles)
	admin.POST("/product/:operation", app.HandleAdminProductOperations, app.IfNotLogined)

	admin.GET("/inventory", ServeHTML("./public/admin_inventory/index.html"), app.IfNotLogined)
	admin.POST("/inventory", app.HandleInventory, app.IfNotLogined)

	admin.GET("/order", ServeHTML("./public/admin_order/index.html"), app.IfNotLogined)
	admin.POST("/order", app.HandleAdminOrder, app.IfNotLogined)

	return e
}

func ServeHTML(htmlPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File(htmlPath)
	}
}

func (app *Application) HandleSignUp(c echo.Context) error {
	user, err := userFromContext(c)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	exists, err := app.UsersStore.IsExists(map[string]interface{}{"username": user.Username})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}

	if err := app.UsersStore.Create(*user); err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := setSession(c); err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := c.Redirect(http.StatusSeeOther, "/user/home/"); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) HandleSignIn(c echo.Context) error {
	username := c.FormValue("username")
	user, err := app.UsersStore.GetOne(map[string]interface{}{"username": username})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, "user not found")
		}

		c.Logger().Error(err)
		return err
	}

	password := c.FormValue("password")
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "wrong password")
	}

	if err := setSession(c); err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := c.Redirect(http.StatusSeeOther, "/user/home/"); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) HandleLogout(c echo.Context) error {
	if err := clearSessionHandler(c); err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := c.Redirect(http.StatusSeeOther, "/"); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) HandleAdminProductFiles(c echo.Context) error {
	operation := c.Param("operation")
	if operation != "add" && operation != "remove" && operation != "update" {
		return echo.NewHTTPError(http.StatusNotFound, "invalid operation")
	}

	switch operation {
	case "add":
		return ServeHTML("./public/admin_add_product/index.html")(c)
	case "remove":
		return ServeHTML("./public/admin_remove_product/index.html")(c)
	case "update":
		return ServeHTML("./public/admin_update_product/index.html")(c)
	}

	return nil
}

func (app *Application) HandleAdminProductOperations(c echo.Context) error {
	operation := c.Param("operation")
	if operation != "add" && operation != "remove" && operation != "update" {
		return echo.NewHTTPError(http.StatusNotFound, "invalid operation")
	}

	switch operation {
	case "add":
		return app.addProduct(c)
	case "remove":
		return app.removeProduct(c)
	case "update":
		return app.updateProduct(c)
	}

	return nil
}

func (app *Application) HandleInventory(c echo.Context) error {
	productName := c.FormValue("product_name")
	operation := c.FormValue("operation")

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	req := &models.UpdateInventoryRequest{
		ProductName: productName,
		Operation:   operation,
		Quantity:    quantity,
	}

	if err := app.InventoryService.Update(req); err != nil {
		if err == models.ErrInvalidOperaton || err == models.ErrProductNotFound || err == models.ErrInvalidQuantity {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) HandleUserOrder(c echo.Context) error {
	req, err := createOrderReqFromContext(c)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	if _, err := app.OrderService.Create(req); err != nil {
		if err == models.ErrProductNotFound || err == models.ErrInvalidQuantity {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (app *Application) HandleMyOrders(c echo.Context) error {
	session := c.Get("session").(*sessions.Session)
	username := session.Values["username"].(string)

	resp, err := app.OrderService.GetOrdersByUsername(username)
	if err != nil {
		if err == models.ErrMatchingRecordNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		c.Logger().Error(err)
		return err
	}

	if err := c.JSONPretty(http.StatusOK, resp, "    "); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app *Application) HandleGetProducts(c echo.Context) error {
	req := &models.GetProductsRequest{}
	category := c.FormValue("category")
	if category == "" {
		req = nil
	} else {
		req.Category = category
	}

	resp, err := app.ProductService.Get(req)
	if err != nil {
		if err == models.ErrMatchingRecordNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		c.Logger().Error(err)
		return err
	}

	if err := c.JSONPretty(http.StatusOK, resp, "    "); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app *Application) HandleAdminOrder(c echo.Context) error {
	req := models.GetOrdersForAdminRequest{
		Username:    c.FormValue("username"),
		ProductName: c.FormValue("product_name"),
		SortBy:      c.FormValue("sort"),
	}

	resp, err := app.OrderService.GetOrdersForAdmin(&req)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := c.JSONPretty(http.StatusOK, resp, "    "); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

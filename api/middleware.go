package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Application) createSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := app.CookieStore.Get(c.Request(), "signin") // this will also create the cookie if it does not exists
		if err != nil {
			return err
		}

		c.Set("session", session)
		return next(c)
	}
}

func (app *Application) IfAlreadyLogined(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.alreadyLoggedIn(c) {
			return c.Redirect(http.StatusFound, "/user/home")
		}
		return next(c)
	}
}

func (app *Application) IfNotLogined(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		if !app.alreadyLoggedIn(c) {
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}

func (app *Application) IfNotAdmin(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		if !app.isAdmin(c) {
			return echo.NewHTTPError(http.StatusForbidden, "only admins are allowed")
		}
		return next(c)
	}
}

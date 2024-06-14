package server

import "github.com/labstack/echo/v4"

// @Summary	Login
// @Description	Login
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Success	200	{string}	string	"Login"
// @Router		/api/login [post]
func login(c echo.Context) error {
	return c.String(200, "Login")
}

// @Summary	Logout
// @Description	Logout
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Success	200	{string}	string	"Logout"
// @Router		/api/logout [post]
func logout(c echo.Context) error {
	return c.String(200, "Logout")
}

// @Summary	Register
// @Description	Register
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Success	200	{string}	string	"Register"
// @Router		/api/register [post]
func register(c echo.Context) error {
	return c.String(200, "Register")
}

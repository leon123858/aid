package server

import "github.com/labstack/echo/v4"

// @Summary		Login
// @Description	Login
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		LoginRequest	true	"Login Request"
// @Success		200	{object}	Response		"JWT Token"
// @Router			/api/login [post]
func login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Logout
// @Description	Logout
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		LogoutRequest	true	"Logout Request"
// @Success		200	{object}	Response		"empty string"
// @Router			/api/logout [post]
func logout(c echo.Context) error {
	req := LogoutRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Register
// @Description	Register
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		RegisterRequest	true	"Register Request"
// @Success		200	{object}	Response		"empty string"
// @Router			/api/register [post]
func register(c echo.Context) error {
	req := RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Ask
// @Description	Ask
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		AskRequest	true	"Ask Request"
// @Success		200	{object}	Response	"aid string"
// @Router			/api/ask [post]
func ask(c echo.Context) error {
	req := AskRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

// @Summary		Trigger
// @Description	Trigger
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			req	body		TriggerRequest	true	"Trigger Request"
// @Success		200	{object}	Response		"aid string"
// @Router			/api/trigger [post]
func trigger(c echo.Context) error {
	req := TriggerRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return c.JSON(200, generateResponse(true, ""))
}

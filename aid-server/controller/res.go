package controller

import "github.com/labstack/echo/v4"

type ContextWrapper struct {
	echo.Context
}

type Response struct {
	Result bool        `json:"result"`
	Data   interface{} `json:"data"`
}

func (c *ContextWrapper) newError(code int, msg string) error {
	return c.JSON(code, Response{
		Result: false,
		Data:   msg,
	})
}

func (c *ContextWrapper) newSuccess(data interface{}) error {
	return c.JSON(200, Response{
		Result: true,
		Data:   data,
	})
}

func (c *ContextWrapper) newBadRequestError(msg string) error {
	return c.newError(400, msg)
}

func (c *ContextWrapper) newInternalServerError(msg string) error {
	return c.newError(500, msg)
}

func (c *ContextWrapper) newNotFound(msg string) error {
	return c.newError(404, msg)
}

func (c *ContextWrapper) newUnauthorized(msg string) error {
	return c.newError(401, msg)
}

func (c *ContextWrapper) newForbidden(msg string) error {
	return c.newError(403, msg)
}

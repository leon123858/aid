package device

import "github.com/labstack/echo/v4"

func SetRealIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		c.Set("ip", ip)
		return next(c)
	}
}

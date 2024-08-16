package main

import (
	"aid-server/controller"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Route
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.GET("/verify/hash", controller.GetCertHash)
	e.POST("/register/cert", controller.SaveCertHash)
	ac := e.Group("/ac")
	ac.GET("/get/key", controller.GetServerPublicKey)
	ac.POST("/sign/cert", controller.AskServerSignCert)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

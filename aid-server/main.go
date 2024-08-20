package main

import (
	"aid-server/controller"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "aid-server/docs"
	"github.com/swaggo/echo-swagger"
)

//	@title			AID Server API
//	@version		0.1
//	@description	This is the AID Server API DEMO

//	@contact.name	Leon Lin
//	@contact.url	github.com/leon123858
//	@contact.email	a0970785699@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		127.0.0.1:7001
// @schemes	http https
// @BasePath	/
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
	e.Logger.Fatal(e.Start(":7001"))
}

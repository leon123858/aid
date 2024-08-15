package server

import (
	_ "aid-server/docs"
	"aid-server/pkg/device"
	"aid-server/pkg/jwt"
	"aid-server/pkg/res"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			AID API Server
//	@version		1.0
//	@description	This is a AID server implementation for my paper.

//	@contact.name	Leon Lin
//	@contact.url	github.com/leon123858
//	@contact.email	a0970785699@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		127.0.0.1:8080
// @BasePath	/
func generateRouter() *echo.Echo {
	router := echo.New()
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}:${uri} => ${status} , from ${remote_ip}, ${latency_human} (${bytes_in}/${bytes_out})\n",
	}))
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())

	// swagger docs
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	// health check
	router.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	api := router.Group("/api")
	{
		api.Use(device.SetRealIP)
		api.POST("/login", login)
		api.POST("/register", register)
		api.POST("/ask", ask)
		api.POST("/check", check)
		api.POST("/verify", verify, jwt.GenerateParseJwtMiddle(res.GenerateResponse))
	}

	usage := router.Group("/usage")
	{
		usage.POST("/login", loginAlias)
		usage.POST("/register", registerAlias)
	}

	return router
}

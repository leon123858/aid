package main

import (
	"backend/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leon123858/aidgo"
	"net/http"
)

type (
	TodoItem struct {
		ID   int    `json:"id"`
		Task string `json:"task"`
		Done bool   `json:"done"`
	}

	LoginRequest struct {
		Cert aidgo.AidCert `json:"cert"`
	}
)

var (
	aidVerifier     = aidgo.NewVerifier()
	verifyGenerator = aidgo.NewVerifyGenerator()
	aidServerClient = services.NewAIDServerClient()
)

func init() {
	verifyGenerator.P2p = func(aid uuid.UUID, option string, msg interface{}, certOption interface{}) error {
		if option != "rsa" {
			return aidgo.NewNotImplementedError("option not implemented")
		}
		// msg is a ["Hello World!", "signature"]
		originalString := msg.([]string)[0]
		signatureBase64 := msg.([]string)[1]
		publicKeyPemString := certOption.(string)
		// use default rsa verify algo
		return aidgo.DefaultRsaVerifyAlgo(originalString, signatureBase64, publicKeyPemString)
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/todos/:aid", getTodos)
	e.POST("/login/:aid", login)
	e.POST("/logout/:aid", logout, middlewareFunc)
	e.POST("/todos/:aid", createTodos, middlewareFunc)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// check cert bind aid
	if req.Cert.Aid.String() != c.Param("aid") {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "aid not match"})
	}
	// cert hash should match hash in aid server
	hash, err := aidServerClient.RequestHash(c.Param("aid"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	if req.Cert.Hash() != hash {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "hash not match"})
	}
	// save cert
	err = aidVerifier.SaveCert(req.Cert)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

func logout(c echo.Context) error {
	aid := c.Param("aid")
	// parse str to uuid
	aidUUID, err := uuid.Parse(aid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": err.Error()})
	}
	// get cert
	cert, err := aidVerifier.GetCert(aidUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	// clear cert and data
	err = aidVerifier.ClearCert(cert)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	err = aidVerifier.ClearData(aidUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

func getTodos(c echo.Context) error {
	aid := c.Param("aid")
	// parse str to uuid
	aidUUID, err := uuid.Parse(aid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": err.Error()})
	}
	// get data
	data, err := aidVerifier.GetData(aidUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	return c.JSON(http.StatusOK, data.Data["todoList"].(*[]TodoItem))
}

func createTodos(c echo.Context) error {
	aid := c.Param("aid")
	req := new([]TodoItem)
	// bind request
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": err.Error()})
	}
	// parse str to uuid
	aidUUID, err := uuid.Parse(aid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": err.Error()})
	}
	// save data
	err = aidVerifier.SaveData(aidgo.AidData{
		Aid:  aidUUID,
		Data: map[string]interface{}{"todoList": req},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

func middlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sign := c.Request().Header.Get("Sign")
		preSign := c.Request().Header.Get("PreSign")
		aid := c.Param("aid")
		aidUUID, err := uuid.Parse(aid)
		if err != nil {
			return echo.ErrUnauthorized
		}
		// 允許的最大時間差(秒)
		err = aidgo.DefaultTimestampTimeoutAlgo(preSign, 60)
		if err != nil {
			return echo.ErrUnauthorized
		}

		// verify sign, hash preSign and decrypt sign
		err = aidVerifier.VerifyCert(aidUUID, "rsa", []string{preSign, sign}, verifyGenerator)
		if err != nil {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

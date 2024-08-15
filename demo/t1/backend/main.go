package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/leon123858/aidgo"
	"net/http"
	"strconv"
	"time"
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
		// base64 to []byte
		signature, err := base64.StdEncoding.DecodeString(signatureBase64)
		if err != nil {
			return err
		}
		// publicKey to rsa.PublicKey
		block, _ := pem.Decode([]byte(publicKeyPemString))
		if block == nil || block.Type != "PUBLIC KEY" {
			return aidgo.NewInternalServerError("failed to decode PEM block containing public key")
		}
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
		if !ok {
			return aidgo.NewInternalServerError("failed to parse public key")
		}
		// verify the signature
		hashed := sha256.Sum256([]byte(originalString))
		err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signature)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
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
	// save cert
	err := aidVerifier.SaveCert(req.Cert)
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
		const maxTimeDiff int64 = 60

		timestamp, err := strconv.ParseInt(preSign, 10, 64)
		if err != nil {
			return echo.ErrUnauthorized
		}

		now := time.Now().Unix()
		timeDiff := now - timestamp/1000

		if timeDiff > maxTimeDiff || timeDiff < -maxTimeDiff {
			return echo.ErrUnauthorized
		}

		// verify sign, hash preSign and decrypt sign
		err = aidVerifier.VerifyCert(aidUUID, "rsa", []string{preSign, sign}, verifyGenerator)
		println(err)
		if err != nil {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

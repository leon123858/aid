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
	"net/http"
	"strconv"
	"time"

	"github.com/leon123858/aidgo"
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
		// certOption is a map[string]string{"publicKey": "base64 encoded public key"}
		publicKeyBase64 := certOption.(map[string]string)["publicKey"]
		// base64 to []byte
		signature, err := base64.StdEncoding.DecodeString(signatureBase64)
		if err != nil {
			return err
		}
		publicKeyByte, err := base64.StdEncoding.DecodeString(publicKeyBase64)
		if err != nil {
			return err
		}
		// publicKey to rsa.PublicKey
		block, _ := pem.Decode(publicKeyByte)
		if block == nil {
			return aidgo.NewBadRequestError("Failed to decode PEM block")
		}
		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return err
		}
		// verify the signature
		hashedMsg := sha256.Sum256([]byte(originalString))
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedMsg[:], signature)
		return err
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/todos/:aid", getTodos)
	e.POST("/login/:aid", login)
	// verify middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sign := c.Request().Header.Get("Sign")
			preSign := c.Request().Header.Get("PreSign")
			aid := c.Param("aid")
			aidUUID, err := uuid.Parse(aid)
			if err != nil {
				return echo.ErrUnauthorized
			}
			// check preSign is near current time in 5 seconds, so convert string to int64
			const maxTimeDiff int64 = 5 // 允許的最大時間差(秒)

			timestamp, err := strconv.ParseInt(preSign, 10, 64)
			if err != nil {
				return echo.ErrUnauthorized
			}

			now := time.Now().Unix()
			timeDiff := now - timestamp

			if timeDiff > maxTimeDiff || timeDiff < -maxTimeDiff {
				return echo.ErrUnauthorized
			}

			// verify sign, hash preSign and decrypt sign
			err = aidVerifier.VerifyCert(aidUUID, "rsa", []string{preSign, sign}, verifyGenerator)
			if err != nil {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	})
	e.POST("/logout/:aid", logout)
	e.POST("/todos/:aid", createTodos)

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
	return c.JSON(http.StatusOK, data.Data["todoList"].([]TodoItem))
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
	// get data
	err = aidVerifier.SaveData(aidgo.AidData{
		Aid:  aidUUID,
		Data: map[string]interface{}{"todoList": req},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

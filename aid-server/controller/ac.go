package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/leon123858/aidgo"
	"net/http"
)

type CertRequest struct {
	Cert aidgo.AidCert `json:"cert"`
	Info interface{}   `json:"info"`
}

func GetServerPublicKey(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func AskServerSignCert(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

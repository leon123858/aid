package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HashRequest struct {
	Hash string `json:"hash"`
	Aid  string `json:"aid"`
}

func GetCertHash(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func SaveCertHash(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

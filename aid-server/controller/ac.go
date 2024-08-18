package controller

import (
	"aid-server/service/sign"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/leon123858/aidgo"
)

type CertRequest struct {
	Cert aidgo.AidCert `json:"cert"`
	Info interface{}   `json:"info"`
}

var cachePublicKey string = ""

func GetServerPublicKey(c echo.Context) error {
	cw := ContextWrapper{c}
	if cachePublicKey == "" {
		kp := sign.GenerateKey()
		cachePublicKey = kp.PublicKey
	}
	return cw.newSuccess(cachePublicKey)
}

func AskServerSignCert(c echo.Context) error {
	req := new(CertRequest)
	cw := ContextWrapper{c}
	if err := c.Bind(req); err != nil {
		return cw.newBadRequestError(err.Error())
	}
	// write certificate login here
	// ...
	hashAidCert := req.Cert.Hash()
	// use rsa to signInCert the hash
	kp := sign.GenerateKey()
	privateKey, _, err := kp.ToCryptoKeys()
	if err != nil {
		return cw.newInternalServerError(err.Error())
	}
	// hash the hashAidCert
	hashHashAidCert := sha256.Sum256([]byte(hashAidCert))
	signature, err := sign.RsaSign(privateKey, hashHashAidCert[:])
	if err != nil {
		return cw.newInternalServerError(err.Error())
	}
	// convert signature to base64
	signInCert := base64.URLEncoding.EncodeToString(signature)
	req.Cert.ServerInfo.Sign = signInCert
	return cw.newSuccess(req.Cert)
}

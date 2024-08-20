package controller

import (
	"aid-server/service/sign"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/leon123858/aidgo"
)

// CertRequestSwagger swagger: CertRequest
type CertRequestSwagger struct {
	// Cert is the AID certificate to be signed.
	// Note: The actual structure of aidgo.AidCert is defined in an external package.
	Cert interface{} `json:"cert"`
	// Info contains additional information for the certificate signing process.
	Info interface{} `json:"info"`
}

// CertRequest contains the certificate to be signed and additional info.
type CertRequest struct {
	// Cert is the AID certificate to be signed.
	// Note: The actual structure of aidgo.AidCert is defined in an external package.
	Cert aidgo.AidCert `json:"cert"`
	// Info contains additional information for the certificate signing process.
	Info interface{} `json:"info"`
}

var cachePublicKey string = ""

// GetServerPublicKey godoc
//
//	@Summary		Get server public key
//	@Description	Retrieve the server's public key
//	@Tags			certificate
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string	"Server public key"
//	@Router			/ac/get/key [get]
func GetServerPublicKey(c echo.Context) error {
	cw := ContextWrapper{c}
	if cachePublicKey == "" {
		kp := sign.GenerateKey()
		cachePublicKey = kp.PublicKey
	}
	return cw.newSuccess(cachePublicKey)
}

// AskServerSignCert godoc
//
//	@Summary		Request server to sign a certificate
//	@Description	Send a certificate to be signed by the server
//	@Tags			certificate
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CertRequestSwagger	true	"Certificate request"
//	@Success		200		{object}	interface{}			"Signed certificate, ref aidgo.AidCert"
//	@Failure		400		{object}	string				"Bad Request"
//	@Failure		500		{object}	string				"Internal Server Error"
//	@Router			/ac/sign/cert [post]
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

package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// IsValidPublicKey 驗證 RSA 公鑰是否合法
func IsValidPublicKey(publicKey []byte) bool {
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return false
	}

	_, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	return true
}

// VerifySignature 驗證 RSA 簽名是否合法
func VerifySignature(publicKey []byte, data []byte, signature string) (bool, error) {
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to decode signature: %s", err))
	}

	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, errors.New("invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("public key is not an RSA public key")
	}

	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, digest, decodedSignature)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GenerateSignature(privateKey []byte, data []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("invalid private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return privateKey, &privateKey.PublicKey
}

func MarshalPublicKey(publicKey *rsa.PublicKey) []byte {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicKeyPEM
}

func MarshalPrivateKey(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyPEM
}

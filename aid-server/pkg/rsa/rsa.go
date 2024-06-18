package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
func VerifySignature(publicKey []byte, data []byte, signature []byte) (bool, error) {
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

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, digest, signature)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Encrypt 使用 RSA 公鑰加密數據
func Encrypt(publicKey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an RSA public key")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// Decrypt 使用 RSA 私鑰解密數據
func Decrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
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

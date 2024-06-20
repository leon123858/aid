package rsa

import (
	"aid-server/pkg/timestamp"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidPublicKey(t *testing.T) {
	_, publicKey := GenerateRSAKeyPair()
	publicKeyPEM := MarshalPublicKey(publicKey)

	assert.True(t, IsValidPublicKey(publicKeyPEM), "Valid public key is considered invalid")

	invalidPublicKey := []byte("invalid public key")
	assert.False(t, IsValidPublicKey(invalidPublicKey), "Invalid public key is considered valid")
}

func TestVerifySignature(t *testing.T) {
	privateKey, publicKey := GenerateRSAKeyPair()
	publicKeyPEM := MarshalPublicKey(publicKey)

	data := []byte("test data")
	hash := sha256.Sum256(data)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	b64Data := base64.StdEncoding.EncodeToString(signature)

	valid, err := VerifySignature(publicKeyPEM, data, b64Data)
	assert.NoError(t, err, "Unexpected error")
	assert.True(t, valid, "Valid signature is considered invalid")

	invalidSignature := []byte("invalid signature")
	b64Data = base64.StdEncoding.EncodeToString(invalidSignature)
	valid, err = VerifySignature(publicKeyPEM, data, b64Data)
	assert.Error(t, err, "Expected an error for invalid signature")
	assert.False(t, valid, "Invalid signature is considered valid")
}

func TestGenerateSignature(t *testing.T) {
	privateKey, publicKey := GenerateRSAKeyPair()
	privateKeyPEM := MarshalPrivateKey(privateKey)

	data := []byte(timestamp.GetTime().String())
	signature, err := GenerateSignature(privateKeyPEM, data)
	require.NoError(t, err, "Unexpected error during signature generation")

	valid, err := VerifySignature(MarshalPublicKey(publicKey), data, signature)
	assert.NoError(t, err, "Unexpected error during signature verification")
	assert.True(t, valid, "Valid signature is considered invalid")
}

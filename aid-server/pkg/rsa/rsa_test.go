package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

	valid, err := VerifySignature(publicKeyPEM, data, signature)
	assert.NoError(t, err, "Unexpected error")
	assert.True(t, valid, "Valid signature is considered invalid")

	invalidSignature := []byte("invalid signature")
	valid, err = VerifySignature(publicKeyPEM, data, invalidSignature)
	assert.Error(t, err, "Expected an error for invalid signature")
	assert.False(t, valid, "Invalid signature is considered valid")
}

func TestEncryptDecrypt(t *testing.T) {
	privateKey, publicKey := GenerateRSAKeyPair()
	publicKeyPEM := MarshalPublicKey(publicKey)
	privateKeyPEM := MarshalPrivateKey(privateKey)

	data := []byte("test data")
	ciphertext, err := Encrypt(publicKeyPEM, data)
	require.NoError(t, err, "Unexpected error during encryption")

	plaintext, err := Decrypt(privateKeyPEM, ciphertext)
	require.NoError(t, err, "Unexpected error during decryption")

	assert.Equal(t, data, plaintext, "Decrypted data does not match the original data")
}

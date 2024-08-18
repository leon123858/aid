package sign

import (
	"crypto"
	"crypto/rsa"
)

func RsaSign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(nil, key, crypto.SHA256, data)
}

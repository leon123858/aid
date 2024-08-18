package sign

import (
	"aid-server/repository"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"log"
	"sync"
	"time"
)

type PemKeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func (p *PemKeyPair) ToCryptoKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	type keyResult struct {
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		err        error
	}

	privateKeyChan := make(chan keyResult)
	publicKeyChan := make(chan keyResult)

	// 並行轉換私鑰
	go func() {
		defer wg.Done()
		privKey, err := pemToPrivateKey(p.PrivateKey)
		privateKeyChan <- keyResult{privateKey: privKey, err: err}
	}()

	// 並行轉換公鑰
	go func() {
		defer wg.Done()
		pubKey, err := pemToPublicKey(p.PublicKey)
		publicKeyChan <- keyResult{publicKey: pubKey, err: err}
	}()

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(privateKeyChan)
		close(publicKeyChan)
	}()

	// 收集結果
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	for i := 0; i < 2; i++ {
		select {
		case result := <-privateKeyChan:
			if result.err != nil {
				return nil, nil, result.err
			}
			privateKey = result.privateKey
		case result := <-publicKeyChan:
			if result.err != nil {
				return nil, nil, result.err
			}
			publicKey = result.publicKey
		}
	}

	return privateKey, publicKey, nil
}

func pemToPrivateKey(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func pemToPublicKey(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type RSA")
	}

	return publicKey, nil
}

func GenerateKey() PemKeyPair {
	// check if db have keys
	keys, err := getKeysFromDB()
	if err != nil {
		// generate new key
		privateKey, publicKey := generateNewKey()
		// save to db
		err = saveKeysToDB(privateKey, publicKey)
		if err != nil {
			log.Fatalf("Failed to save keys to db: %v", err)
		}
		keys.PrivateKey = privateKey
		keys.PublicKey = publicKey
	}
	if keys.PrivateKey == "" || keys.PublicKey == "" {
		log.Fatalf("Failed to get keys from db: %v", err)
	}
	return keys
}

func getKeysFromDB() (PemKeyPair, error) {
	type keyResult struct {
		pemString string
		err       error
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	publicPemStrChan := make(chan keyResult)
	privatePemStrChan := make(chan keyResult)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		read, err := repository.LDB.Read([]byte("publicKey"))
		if err != nil {
			publicPemStrChan <- keyResult{pemString: "", err: err}
			return
		}
		if read == nil {
			publicPemStrChan <- keyResult{pemString: "", err: nil}
			return
		}
		publicPemStrChan <- keyResult{pemString: string(read), err: nil}
	}()
	go func() {
		defer wg.Done()
		read, err := repository.LDB.Read([]byte("privateKey"))
		if err != nil {
			privatePemStrChan <- keyResult{pemString: "", err: err}
			return
		}
		if read == nil {
			privatePemStrChan <- keyResult{pemString: "", err: nil}
			return
		}
		privatePemStrChan <- keyResult{pemString: string(read), err: nil}
	}()
	go func() {
		wg.Wait()
		close(publicPemStrChan)
		close(privatePemStrChan)
	}()
	// select error or success
	kp := PemKeyPair{}
	for i := 2; i > 0; i-- {
		select {
		case result := <-publicPemStrChan:
			if result.err != nil {
				return kp, result.err
			}
			kp.PublicKey = result.pemString
		case result := <-privatePemStrChan:
			if result.err != nil {
				return kp, result.err
			}
			kp.PrivateKey = result.pemString
		case <-ctx.Done():
			return kp, errors.New("timeout")
		}
	}
	return kp, nil
}

func generateNewKey() (string, string) {
	// 生成RSA密鑰對
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	// 從私鑰中獲取公鑰
	publicKey := &privateKey.PublicKey

	// 將私鑰轉換為PEM格式
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyPEMBytes := pem.EncodeToMemory(privateKeyPEM)

	// 將公鑰轉換為PEM格式
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEM)

	return string(privateKeyPEMBytes), string(publicKeyPEMBytes)
}

func saveKeysToDB(privateKey, publicKey string) error {
	err := repository.LDB.Create([]byte("privateKey"), []byte(privateKey))
	if err != nil {
		return err
	}
	err = repository.LDB.Create([]byte("publicKey"), []byte(publicKey))
	if err != nil {
		return err
	}
	return nil
}

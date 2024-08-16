package auth

import (
	"aid-server/repository"
	"github.com/google/uuid"
)

func LoadHash(aid uuid.UUID) (string, error) {
	hash, err := repository.LDB.Read([]byte(aid.String()))
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func SaveHash(aid uuid.UUID, hash string) error {
	return repository.LDB.Create([]byte(aid.String()), []byte(hash))
}

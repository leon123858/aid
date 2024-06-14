package jwt

import (
	"aid-server/configs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwt_Valid(t *testing.T) {
	configs.Configs.Jwt.Duration = 1 * time.Second
	configs.Configs.Jwt.Secret = "test-secret"
	// should success
	uid := uuid.New().String()
	token, err := GenerateToken(uid)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	claims, err := ParseToken(token)
	assert.Nil(t, err)
	assert.NotNil(t, claims)
	assert.Nil(t, claims.Valid())
	assert.Equal(t, "auth", claims.Subject)
	assert.NotEmpty(t, claims.ID)
	assert.Equal(t, uid, claims.ID)

}

func TestJwt_inValid(t *testing.T) {
	configs.Configs.Jwt.Duration = 1 * time.Second
	configs.Configs.Jwt.Secret = "test-secret"
	uid := uuid.New().String()
	// should fail as expired
	token, err := GenerateToken(uid)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	time.Sleep(time.Millisecond * 1100)
	claims, err := ParseToken(token)
	assert.NotNil(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "token is expired", err.Error())

	// should fail as invalid user id
	token, _ = GenerateToken("invalid-uuid")
	claims, err = ParseToken(token)
	assert.NotNil(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "invalid user id", err.Error())

	// should fail as invalid secret
	configs.Configs.Jwt.Secret = "secret"
	token, _ = GenerateToken(uid)
	configs.Configs.Jwt.Secret = "invalid secret"
	claims, err = ParseToken(token)
	assert.NotNil(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "signature is invalid", err.Error())
}

package jwt

import (
	"aid-server/configs"
	"aid-server/pkg/res"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func TestGenerateParseJwtMiddle(t *testing.T) {
	// 创建一个响应函数，返回自定义的响应格式
	resFunc := func(success bool, message string) res.Response {
		return res.Response{
			Result:  success,
			Content: message,
		}
	}

	// 创建一个 Echo 实例
	e := echo.New()

	// 注册中间件
	e.Use(GenerateParseJwtMiddle(resFunc))

	// 创建一个处理函数，用于测试中间件
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}

	// 测试没有提供 token 的情况
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GenerateParseJwtMiddle(resFunc)(handler)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// 测试提供无效 token 的情况
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "invalid_token")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = GenerateParseJwtMiddle(resFunc)(handler)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// 测试提供有效 token 的情况
	var validToken string
	if validToken, err = GenerateToken(uuid.New().String()); err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", validToken)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = GenerateParseJwtMiddle(resFunc)(handler)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!", rec.Body.String())
}

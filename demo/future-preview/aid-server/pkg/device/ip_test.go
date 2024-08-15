package device

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetRealIP(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 模拟 X-Real-IP 头
	req.Header.Set("X-Real-IP", "192.168.1.1")

	// 创建一个处理函数来检查设置的 IP
	handler := func(c echo.Context) error {
		ip := c.Get("ip").(string)
		assert.Equal(t, "192.168.1.1", ip)
		return c.String(http.StatusOK, "test")
	}

	// 使用 SetRealIP 中间件包装处理函数
	middlewareFunc := SetRealIP(handler)

	// 调用中间件函数
	err := middlewareFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

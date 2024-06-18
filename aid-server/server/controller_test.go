package server

import (
	"aid-server/pkg/rsa"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"aid":"testAID","sign":"testSign","ip":"127.0.0.1","browser":"Chrome"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"result":true,"content":""}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestLogout(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/logout", strings.NewReader(`{"aid":"testAID","ip":"127.0.0.1","browser":"Chrome"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"result":true,"content":""}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestRegister(t *testing.T) {
	e := echo.New()
	_, key := rsa.GenerateRSAKeyPair()
	pubKey := rsa.MarshalPublicKey(key)
	body := RegisterRequest{
		AID:       uuid.New().String(),
		PublicKey: string(pubKey),
		Request: Request{
			Space: Space{
				IP:      "127.0.0.1",
				Browser: "Chrome",
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(bodyBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"result":true,"content":""}`, strings.TrimSpace(rec.Body.String()))
	}

	// test invalid public key
	body.AID = uuid.New().String()
	body.PublicKey = "invalid public key"
	bodyBytes, err = json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(bodyBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, register(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"result":false,"content":"invalid public key"}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestAsk(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/ask", strings.NewReader(`{"ip":"127.0.0.1","browser":"Chrome"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"result":true,"content":""}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestTrigger(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/trigger", strings.NewReader(`{"ip":"127.0.0.1","browser":"Chrome"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, trigger(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"result":true,"content":""}`, strings.TrimSpace(rec.Body.String()))
	}
}

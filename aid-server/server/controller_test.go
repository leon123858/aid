package server

import (
	"aid-server/pkg/jwt"
	"aid-server/pkg/res"
	"aid-server/pkg/rsa"
	"aid-server/pkg/timestamp"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func addUser(e *echo.Echo, t *testing.T) (uuid.UUID, []byte) {
	p, key := rsa.GenerateRSAKeyPair()
	privKey := rsa.MarshalPrivateKey(p)
	pubKey := rsa.MarshalPublicKey(key)
	aidUUID := uuid.New()
	body := RegisterRequest{
		AID:       aidUUID.String(),
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
	}
	return aidUUID, privKey
}

func TestLogin(t *testing.T) {
	e := echo.New()
	aid, key := addUser(e, t)

	ts := timestamp.GetTime()
	sign, err := rsa.GenerateSignature(key, []byte(ts.String()))
	assert.Nil(t, err)
	loginRequest := LoginRequest{
		AID:       aid.String(),
		Sign:      sign,
		Timestamp: ts.String(),
		Request: Request{
			Space: Space{
				IP:      "127.0.0.1",
				Browser: "Chrome",
			},
		},
	}

	jsonData, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonData)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		//println(rec.Body.String())
		var result res.Response
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		if assert.NoError(t, err) {
			assert.True(t, result.Result)
			assert.NotEmpty(t, result.Content)
			claims, err := jwt.ParseToken(result.Content)
			assert.NoError(t, err)
			assert.Equal(t, aid.String(), claims.ID)
		}
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
		k := make(map[string]interface{})
		if err = json.Unmarshal([]byte(rec.Body.String()), &k); err != nil {
			t.Fatal(err)
		}
		str := fmt.Sprintf(`{"result":true,"content":"%s"}`, k["content"])
		assert.Equal(t, str, strings.TrimSpace(rec.Body.String()))
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

func TestAskAndTrigger(t *testing.T) {
	// Setup
	e := echo.New()
	aid, _ := addUser(e, t)

	req := httptest.NewRequest(http.MethodPost, "/api/ask", strings.NewReader(`{
        "ip": "127.0.0.1",
        "browser": "Chrome"
    }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	token, _ := jwt.GenerateToken(aid.String())
	c.Request().Header.Set(echo.HeaderAuthorization, token)
	err := jwt.GenerateParseJwtMiddle(res.GenerateResponse)(ask)(c)
	// Assertions
	var uid string
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp res.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Result)
		assert.NotEmpty(t, resp.Content)
		uid = resp.Content
	}

	req = httptest.NewRequest(http.MethodPost, "/api/trigger", strings.NewReader(`{
        "uid": "`+uid+`",
        "ip": "127.0.0.1",
        "browser": "Chrome"
    }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	// Assertions
	if assert.NoError(t, trigger(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp res.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Result)
		assert.NotEmpty(t, resp.Content)
		assert.Equal(t, string(Online), resp.Content)
	}

	// test invalid
	req = httptest.NewRequest(http.MethodPost, "/api/trigger", strings.NewReader(`{
        "uid": "`+uid+`",
        "ip": "127.0.0.1",
        "browser": "Safari"
    }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	// Assertions
	if assert.NoError(t, trigger(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp res.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Result)
		assert.NotEmpty(t, resp.Content)
		assert.Equal(t, string(Offline), resp.Content)
	}
}

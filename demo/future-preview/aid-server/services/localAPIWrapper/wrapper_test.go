package localAPIWrapper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsk(t *testing.T) {
	// 創建一個模擬的 HTTP 服務器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/ask", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var askRequest AskRequest
		err := json.NewDecoder(r.Body).Decode(&askRequest)
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", askRequest.IP)
		assert.Equal(t, "Chrome", askRequest.Browser)

		response := AIDResponse{
			Result:  true,
			Content: "Ask successful",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 創建一個使用模擬服務器的 API 封裝器
	api := newAPIWrapper(server.URL + "/")

	// 執行 Ask 方法並檢查結果
	response, err := api.Ask("127.0.0.1", "Chrome")
	assert.NoError(t, err)
	assert.True(t, response.Result)
	assert.Equal(t, "Ask successful", response.Content)
}

func TestCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/check", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var checkRequest CheckRequest
		err := json.NewDecoder(r.Body).Decode(&checkRequest)
		assert.NoError(t, err)
		assert.Equal(t, "user123", checkRequest.UID)
		assert.Equal(t, "127.0.0.1", checkRequest.IP)
		assert.Equal(t, "Firefox", checkRequest.Browser)

		response := AIDResponse{
			Result:  true,
			Content: "Check successful",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	api := newAPIWrapper(server.URL + "/")

	response, err := api.Check("user123", "127.0.0.1", "Firefox")
	assert.NoError(t, err)
	assert.True(t, response.Result)
	assert.Equal(t, "Check successful", response.Content)
}

func TestVerify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/verify", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "token123", r.Header.Get("Authorization"))

		var verifyRequest VerifyRequest
		err := json.NewDecoder(r.Body).Decode(&verifyRequest)
		assert.NoError(t, err)
		assert.Equal(t, "user123", verifyRequest.UID)

		response := AIDResponse{
			Result:  true,
			Content: "Verify successful",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	api := newAPIWrapper(server.URL + "/")

	response, err := api.Verify("token123", "user123")
	assert.NoError(t, err)
	assert.True(t, response.Result)
	assert.Equal(t, "Verify successful", response.Content)
}

func TestErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AIDResponse{
			Result:  false,
			Content: "Error occurred",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	api := newAPIWrapper(server.URL + "/")

	_, err := api.Ask("127.0.0.1", "Chrome")
	assert.Error(t, err)
	assert.Equal(t, "error: Error occurred", err.Error())
}

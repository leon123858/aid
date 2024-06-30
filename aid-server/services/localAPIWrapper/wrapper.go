package localAPIWrapper

import (
	"aid-server/configs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type API interface {
	Ask(ip, browser string) (AIDResponse, error)
	Check(uid, ip, browser string) (AIDResponse, error)
	Verify(token, uid string) (AIDResponse, error)
}

type wrapper struct {
	client *http.Client
	host   string
}

func newAPIWrapper(host string) API {
	if host == "" {
		host = os.Getenv("AID_URL")
		if host == "" {
			host = "http://localhost:8080/"
		}
	}
	return &wrapper{
		host:   host,
		client: &http.Client{},
	}
}

type AskRequest struct {
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type CheckRequest struct {
	UID     string `json:"uid"`
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type VerifyRequest struct {
	UID string `json:"uid"`
}

type AIDResponse struct {
	Result  bool   `json:"result"`
	Content string `json:"content"`
}

func (w *wrapper) postToGoService(endpoint string, request interface{}) (AIDResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return AIDResponse{}, err
	}

	resp, err := w.client.Post(fmt.Sprintf("%sapi/%s", w.host, endpoint), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return AIDResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("close with error")
		}
	}(resp.Body)

	var response AIDResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AIDResponse{}, err
	}
	if !response.Result {
		return AIDResponse{}, fmt.Errorf("error: %s", response.Content)
	}
	return response, nil
}

func (w *wrapper) Ask(ip, browser string) (AIDResponse, error) {
	request := AskRequest{IP: ip, Browser: browser}
	return w.postToGoService("ask", request)
}

func (w *wrapper) Check(uid, ip, browser string) (AIDResponse, error) {
	request := CheckRequest{UID: uid, IP: ip, Browser: browser}
	return w.postToGoService("check", request)
}

func (w *wrapper) Verify(token, uid string) (AIDResponse, error) {
	request := VerifyRequest{UID: uid}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return AIDResponse{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/verify", w.host), bytes.NewBuffer(jsonData))
	if err != nil {
		return AIDResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := w.client.Do(req)
	if err != nil {
		return AIDResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("close with error")
		}
	}(resp.Body)

	var response AIDResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AIDResponse{}, err
	}

	if !response.Result {
		return AIDResponse{}, fmt.Errorf("error: %s", response.Content)
	}

	return response, nil
}

func New() API {
	return newAPIWrapper(fmt.Sprintf("http://127.0.0.1:%s/", configs.Configs.Host.Port))
}

package services

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

const AIDServerURL = "http://localhost:7001"

type AIDServerClientInterface interface {
	RequestHash(aid string) (string, error)
}

type AIDServerClient struct {
	client *resty.Client
}

func NewAIDServerClient() AIDServerClientInterface {
	return &AIDServerClient{
		client: resty.New(),
	}
}

func (c *AIDServerClient) RequestHash(aid string) (string, error) {
	resp, err := c.client.R().
		SetQueryParam("aid", aid).
		Get(AIDServerURL + "/verify/hash")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", err
	}
	// parse response to struct
	type Response struct {
		Result bool   `json:"result"`
		Data   string `json:"data"`
	}
	var res Response
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	if !res.Result {
		return "", err
	}
	return res.Data, nil
}

package main

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/skyline93/syncbyte-go/internal/pkg/schema"
)

var client *Client

func InitClient() {
	client = NewClient()
}

type Client struct {
	*resty.Client
}

func NewClient() *Client {
	c := resty.New().SetBaseURL(conf.Server).SetDoNotParseResponse(true)
	return &Client{Client: c}
}

func (c *Client) Get(url string) (interface{}, error) {
	resp, err := c.Client.R().Get(url)
	if err != nil {
		return nil, err
	}

	respBody := schema.ResponseBody{}
	if err := json.NewDecoder(resp.RawResponse.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) Post(url string, body interface{}) (interface{}, error) {
	resp, err := client.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := schema.ResponseBody{}
	if err := json.NewDecoder(resp.RawResponse.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

package pcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	BaseURLUS = "https://api.pcloud.com"
	BaseURLEU = "https://eapi.pcloud.com"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       string
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = BaseURLUS
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) SetAuth(auth string) {
	c.auth = auth
}

func (c *Client) Auth() string {
	return c.auth
}

func (c *Client) do(method string, params url.Values, result any) error {
	if c.auth != "" {
		params.Set("auth", c.auth)
	}

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, method, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) doPost(method string, params url.Values, body io.Reader, contentType string, result any) error {
	if c.auth != "" {
		params.Set("auth", c.auth)
	}

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, method, params.Encode())
	resp, err := c.httpClient.Post(reqURL, contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

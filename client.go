package pcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

const (
	BaseURLUS  = "https://api.pcloud.com"
	BaseURLEU  = "https://eapi.pcloud.com"
	DefaultRPM = 100.0
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       string
	logger     *slog.Logger
	limiter    *rate.Limiter
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = BaseURLUS
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		logger:     newNoopLogger(),
		limiter:    rate.NewLimiter(rate.Limit(DefaultRPM/60.0), 1),
	}
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

func (c *Client) SetRateLimit(rpm float64) {
	if rpm < DefaultRPM {
		rpm = DefaultRPM
	}
	c.limiter = rate.NewLimiter(rate.Limit(rpm/60.0), 1)
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

	c.logger.Debug("request", "method", method)
	c.limiter.Wait(context.Background())

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, method, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		c.logger.Error("request failed", "method", method, "error", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		c.logger.Error("decode failed", "method", method, "error", err)
		return err
	}
	return nil
}

func (c *Client) doPost(method string, params url.Values, body io.Reader, contentType string, result any) error {
	if c.auth != "" {
		params.Set("auth", c.auth)
	}

	c.logger.Debug("request", "method", method)
	c.limiter.Wait(context.Background())

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, method, params.Encode())
	resp, err := c.httpClient.Post(reqURL, contentType, body)
	if err != nil {
		c.logger.Error("request failed", "method", method, "error", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		c.logger.Error("decode failed", "method", method, "error", err)
		return err
	}
	return nil
}

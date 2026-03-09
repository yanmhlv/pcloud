package pcloud

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	// BaseURLUS is the pCloud API endpoint for US-region accounts.
	BaseURLUS = "https://api.pcloud.com"
	// BaseURLEU is the pCloud API endpoint for EU-region accounts.
	BaseURLEU = "https://eapi.pcloud.com"
	// MinRPM is the minimum allowed rate limit in requests per minute.
	MinRPM = 100.0
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is a pCloud API client. Use NewClient to create one.
// All methods are safe for concurrent use.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	auth        string
	tokenSource oauth2.TokenSource
	logger      *slog.Logger
	limiter     *rate.Limiter
}

// NewClient creates a new Client for the given base URL.
// Pass BaseURLUS or BaseURLEU; an empty string defaults to BaseURLUS.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    cmp.Or(baseURL, BaseURLUS),
		httpClient: &http.Client{Timeout: DefaultTimeout},
		logger:     newNoopLogger(),
		limiter:    rate.NewLimiter(rate.Limit(MinRPM/60.0), 10),
	}
}

// SetHTTPClient replaces the default HTTP client.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetLogger attaches a structured logger for request diagnostics.
func (c *Client) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

// SetRateLimit configures the maximum requests per minute.
// Returns an error if rpm is below MinRPM.
func (c *Client) SetRateLimit(rpm float64) error {
	if rpm < MinRPM {
		return fmt.Errorf("rate limit %.1f RPM is below minimum %.1f RPM", rpm, MinRPM)
	}
	c.limiter = rate.NewLimiter(rate.Limit(rpm/60.0), 10)
	return nil
}

// SetTokenSource configures OAuth2 token-based authentication.
// This takes precedence over username/password auth set via Login.
func (c *Client) SetTokenSource(ts oauth2.TokenSource) {
	c.tokenSource = ts
}

func (c *Client) request(ctx context.Context, httpMethod, apiMethod string, params url.Values, body io.Reader, contentType string, result apiError) error {
	if err := c.setAuth(params); err != nil {
		return err
	}

	c.logger.Debug("request", "method", apiMethod)
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, apiMethod, params.Encode())
	req, err := http.NewRequestWithContext(ctx, httpMethod, reqURL, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("request failed", "method", apiMethod, "error", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		c.logger.Error("decode failed", "method", apiMethod, "error", err)
		return err
	}
	return result.Err()
}

func (c *Client) do(ctx context.Context, method string, params url.Values, result apiError) error {
	return c.request(ctx, http.MethodGet, method, params, nil, "", result)
}

func (c *Client) doPost(ctx context.Context, method string, params url.Values, body io.Reader, contentType string, result apiError) error {
	return c.request(ctx, http.MethodPost, method, params, body, contentType, result)
}

func (c *Client) setAuth(params url.Values) error {
	if c.tokenSource != nil {
		token, err := c.tokenSource.Token()
		if err != nil {
			return err
		}
		params.Set("auth", token.AccessToken)
		return nil
	}
	if c.auth != "" {
		params.Set("auth", c.auth)
	}
	return nil
}

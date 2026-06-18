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
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	BaseURLUS = "https://api.pcloud.com"

	BaseURLEU = "https://eapi.pcloud.com"

	MinRPM = 100.0
)

const (
	rateLimiterBurst  = 10
	maxErrorBodyBytes = 512
)

type Client struct {
	mu          sync.RWMutex
	baseURL     string
	httpClient  *http.Client
	auth        string
	tokenSource oauth2.TokenSource
	logger      *slog.Logger
	limiter     *rate.Limiter
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    cmp.Or(baseURL, BaseURLUS),
		httpClient: &http.Client{},
		logger:     newNoopLogger(),
		limiter:    rate.NewLimiter(rate.Limit(MinRPM/60.0), rateLimiterBurst),
	}
}

func (c *Client) SetHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient = httpClient
}

func (c *Client) SetLogger(logger *slog.Logger) {
	if logger == nil {
		logger = newNoopLogger()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

func (c *Client) SetRateLimit(rpm float64) error {
	if rpm < MinRPM {
		return fmt.Errorf("rate limit %.1f RPM is below minimum %.1f RPM", rpm, MinRPM)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.limiter = rate.NewLimiter(rate.Limit(rpm/60.0), rateLimiterBurst)
	return nil
}

func (c *Client) SetTokenSource(ts oauth2.TokenSource) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tokenSource = ts
}

func (c *Client) request(ctx context.Context, httpMethod, apiMethod string, params url.Values, body io.Reader, contentType string, result apiError) error {
	if err := c.setAuth(params); err != nil {
		return err
	}

	c.mu.RLock()
	httpClient := c.httpClient
	logger := c.logger
	limiter := c.limiter
	baseURL := c.baseURL
	c.mu.RUnlock()

	logger.Debug("request", "method", apiMethod)
	if err := limiter.Wait(ctx); err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/%s?%s", baseURL, apiMethod, params.Encode())
	req, err := http.NewRequestWithContext(ctx, httpMethod, reqURL, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("request failed", "method", apiMethod, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		io.Copy(io.Discard, io.LimitReader(resp.Body, maxErrorBodyBytes))
		return fmt.Errorf("pcloud: %s %s: %s", httpMethod, apiMethod, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		logger.Error("decode failed", "method", apiMethod, "error", err)
		return err
	}
	return result.Err()
}

func (c *Client) doGet(ctx context.Context, method string, params url.Values, result apiError) error {
	return c.request(ctx, http.MethodGet, method, params, nil, "", result)
}

func (c *Client) doPost(ctx context.Context, method string, params url.Values, body io.Reader, contentType string, result apiError) error {
	return c.request(ctx, http.MethodPost, method, params, body, contentType, result)
}

func (c *Client) setAuth(params url.Values) error {
	c.mu.RLock()
	ts := c.tokenSource
	auth := c.auth
	c.mu.RUnlock()

	if ts != nil {
		token, err := ts.Token()
		if err != nil {
			return err
		}
		params.Set("auth", token.AccessToken)
		return nil
	}
	if auth != "" {
		params.Set("auth", auth)
	}
	return nil
}

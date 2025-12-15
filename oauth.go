package pcloud

import (
	"context"
	"net/url"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type OAuthToken struct {
	Error
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserID      uint64 `json:"userid"`
}

func (c *Client) OAuthAuthorizeURL(cfg OAuthConfig, state string) string {
	params := url.Values{
		"client_id":     {cfg.ClientID},
		"redirect_uri":  {cfg.RedirectURI},
		"response_type": {"code"},
	}
	if state != "" {
		params.Set("state", state)
	}
	return c.baseURL + "/oauth2_authorize?" + params.Encode()
}

func (c *Client) OAuthExchangeCode(ctx context.Context, cfg OAuthConfig, code string) (*OAuthToken, error) {
	params := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
	}

	var resp OAuthToken
	if err := c.do(ctx, "oauth2_token", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}

	c.auth = resp.AccessToken
	return &resp, nil
}

func (c *Client) SetOAuthToken(token string) {
	c.auth = token
}

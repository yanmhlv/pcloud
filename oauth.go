package pcloud

import (
	"context"
	"fmt"
	"net/url"

	"golang.org/x/oauth2"
)

func EndpointUS() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  BaseURLUS + "/oauth2_authorize",
		TokenURL: BaseURLUS + "/oauth2_token",
	}
}

func EndpointEU() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  BaseURLEU + "/oauth2_authorize",
		TokenURL: BaseURLEU + "/oauth2_token",
	}
}

type exchangeResponse struct {
	Error
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserID      uint64 `json:"userid"`
}

func (c *Client) ExchangeCode(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
	params := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
	}

	var resp exchangeResponse
	if err := c.do(ctx, "oauth2_token", params, &resp); err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	return &oauth2.Token{
		AccessToken: resp.AccessToken,
		TokenType:   resp.TokenType,
	}, nil
}

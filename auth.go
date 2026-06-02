package pcloud

import (
	"context"
	"fmt"
	"net/url"
)

type loginResponse struct {
	Error
	Auth string `json:"auth"`
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	params := url.Values{
		"getauth":  {"1"},
		"username": {username},
		"password": {password},
	}

	var resp loginResponse
	if err := c.do(ctx, "userinfo", params, &resp); err != nil {
		return fmt.Errorf("login: %w", err)
	}

	c.mu.Lock()
	c.auth = resp.Auth
	c.mu.Unlock()
	return nil
}

func (c *Client) Logout(ctx context.Context) error {
	var resp Error
	if err := c.do(ctx, "logout", url.Values{}, &resp); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	c.mu.Lock()
	c.auth = ""
	c.mu.Unlock()
	return nil
}

type userInfoResponse struct {
	Error
	UserInfo
}

func (c *Client) UserInfo(ctx context.Context) (*UserInfo, error) {
	var resp userInfoResponse
	if err := c.do(ctx, "userinfo", url.Values{}, &resp); err != nil {
		return nil, fmt.Errorf("userinfo: %w", err)
	}
	return &resp.UserInfo, nil
}

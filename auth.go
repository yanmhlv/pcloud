package pcloud

import (
	"context"
	"net/url"
)

type loginResponse struct {
	Error
	Auth string `json:"auth"`
}

// Login authenticates with username and password, storing the session token on the client.
func (c *Client) Login(ctx context.Context, username, password string) error {
	params := url.Values{
		"getauth":  {"1"},
		"username": {username},
		"password": {password},
	}

	var resp loginResponse
	if err := c.do(ctx, "userinfo", params, &resp); err != nil {
		return err
	}

	c.auth = resp.Auth
	return nil
}

// Logout invalidates the current session token.
func (c *Client) Logout(ctx context.Context) error {
	var resp Error
	if err := c.do(ctx, "logout", url.Values{}, &resp); err != nil {
		return err
	}

	c.auth = ""
	return nil
}

// UserInfo returns the account details for the authenticated user.
func (c *Client) UserInfo(ctx context.Context) (*UserInfo, error) {
	var resp UserInfo
	if err := c.do(ctx, "userinfo", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

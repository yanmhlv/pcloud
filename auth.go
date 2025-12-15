package pcloud

import "net/url"

type loginResponse struct {
	Error
	Auth string `json:"auth"`
}

func (c *Client) Login(username, password string) error {
	params := url.Values{
		"getauth":  {"1"},
		"username": {username},
		"password": {password},
	}

	var resp loginResponse
	if err := c.do("userinfo", params, &resp); err != nil {
		return err
	}
	if err := resp.Err(); err != nil {
		return err
	}

	c.auth = resp.Auth
	return nil
}

func (c *Client) Logout() error {
	var resp Error
	if err := c.do("logout", url.Values{}, &resp); err != nil {
		return err
	}
	if err := resp.Err(); err != nil {
		return err
	}

	c.auth = ""
	return nil
}

func (c *Client) UserInfo() (*UserInfo, error) {
	var resp UserInfo
	if err := c.do("userinfo", url.Values{}, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

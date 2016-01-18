package pcloud

import (
	"errors"
	"net/http"
)

// Login pcloud login
func (c *Client) Login() error {
	req, err := http.NewRequest("GET", "https://api.pcloud.com/userinfo", nil)
	if err != nil {
		return err
	}

	c.prepareRequest(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	v := struct {
		Authkey string `json:"auth"`
		Result  int    `json:"result"`
	}{}

	if err := c.prepareResponse(resp, &v); err != nil {
		return err
	}

	if v.Result == 0 && v.Authkey != "" {
		c.Authkey = v.Authkey
	}
	return nil
}

// Logout pcloud logout
func (c *Client) Logout() error {
	req, err := http.NewRequest("GET", "https://api.pcloud.com/logout", nil)
	if err != nil {
		return err
	}
	c.prepareRequest(req)

	resp, err := c.client.Do(req)
	if req != nil {
		return err
	}

	v := struct {
		Result  int  `json:"result"`
		Success bool `json:"success"`
	}{}

	if err := c.prepareResponse(resp, &v); err != nil {
		return err
	}

	if v.Success {
		return nil
	}

	if desc, ok := map[int]string{
		1000: "Log in required.",
		2000: "Log in failed.",
		4000: "Too many login tries from this IP address.",
	}[v.Result]; ok {
		return errors.New(desc)
	}
	return errors.New("Unknown logout error")
}

package pcloud

import (
	"encoding/json"
	"errors"
	"net/url"
)

func (c *pCloudClient) Login(username string, password string) error {
	values := url.Values{
		"getauth":  {"1"},
		"username": {username},
		"password": {password},
	}

	resp, err := c.Client.Get(urlBuilder("userinfo", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := struct {
		Auth   string `json:"auth"`
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	c.Auth = &result.Auth
	return nil
}

func (c *pCloudClient) Logout() error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	resp, err := c.Client.Get(urlBuilder("logout", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	c.Auth = nil
	return nil
}

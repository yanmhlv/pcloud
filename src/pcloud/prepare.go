package pcloud

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	host   = "api.pcloud.com"
	scheme = "https"
)

func (c *Client) prepareRequest(r *http.Request) {
	r.URL = &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}
	if c.Authkey != "" && c.Authkey == r.URL.Query().Get("auth") {
		return
	}

	u := r.URL
	query := u.Query()
	query.Del("auth")

	if query.Get("username") == c.Username && query.Get("password") == c.Password && query.Get("getauth") == "1" {
		return
	}

	query.Set("username", c.Username)
	query.Set("password", c.Password)
	query.Set("getauth", "1")
	u.RawQuery = query.Encode()
	r.URL = u
	return
}

func (c *Client) prepareResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := &struct {
		Authkey string `json:"auth"`
	}{}

	if err := json.Unmarshal(data, result); err != nil {
		return err
	}
	c.Authkey = result.Authkey

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

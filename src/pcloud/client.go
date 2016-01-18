package pcloud

import "net/http"

type Client struct {
	client   *http.Client
	Username string
	Password string
	Authkey  string
}

func NewClient(username string, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		client:   &http.Client{},
	}
}

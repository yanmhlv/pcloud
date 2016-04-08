package pcloud

import "net/http"

type pCloudClient struct {
	Auth   *string
	Client *http.Client
}

// NewClient create new pCloudClient
func NewClient() *pCloudClient {
	return &pCloudClient{
		Auth:   nil,
		Client: &http.Client{},
	}
}

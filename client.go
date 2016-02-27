package pcloud

import "net/http"

const (
	apiScheme = "https"
	apiHost   = "api.pcloud.com"
)

type pCloudClient struct {
	Auth   *string
	Client *http.Client
}

func NewClient() *pCloudClient {
	return &pCloudClient{
		Auth:   nil,
		Client: &http.Client{},
	}
}

package pcloud

import "net/url"

const (
	apiScheme = "https"
	apiHost   = "api.pcloud.com"
)

// urlBuilder; return url with GET-params
func urlBuilder(method string, values url.Values) string {
	return (&url.URL{
		Scheme:   apiScheme,
		Host:     apiHost,
		Path:     method,
		RawQuery: values.Encode(),
	}).String()
}

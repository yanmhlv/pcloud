package pcloud

import "net/url"

// urlBuilder; return url with GET-params
func urlBuilder(method string, values url.Values) string {
	const (
		apiScheme = "https"
		apiHost   = "api.pcloud.com"
	)

	u := url.URL{
		Scheme:   apiScheme,
		Host:     apiHost,
		Path:     method,
		RawQuery: values.Encode(),
	}
	return u.String()
}

package pcloud

import "net/url"

func urlBuilder(method string, values url.Values) string {
	return (&url.URL{
		Scheme:   apiScheme,
		Host:     apiHost,
		Path:     method,
		RawQuery: values.Encode(),
	}).String()
}

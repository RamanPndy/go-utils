package goutils

import (
	urlpkg "net/url"
)

func EncodeQueryParams(url string, keyValue map[string]string) string {
	queryParams := urlpkg.Values{}
	for key, value := range keyValue {
		queryParams.Add(key, value)
	}
	return url + queryParams.Encode()
}

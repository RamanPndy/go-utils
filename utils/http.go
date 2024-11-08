package goutils

import (
	urlpkg "net/url"
	"strings"
)

func EncodeQueryParams(url string, keyValue map[string]string) string {
	queryParams := urlpkg.Values{}
	for key, value := range keyValue {
		queryParams.Add(key, value)
	}
	if !strings.HasSuffix(url, "?") {
		url += "?"
	}
	return url + queryParams.Encode()
}

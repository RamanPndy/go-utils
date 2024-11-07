package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestEncodeQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		params   map[string]string
		expected string
	}{
		{
			name:     "Basic parameters",
			url:      "https://example.com?",
			params:   map[string]string{"key1": "value1", "key2": "value2"},
			expected: "https://example.com?key1=value1&key2=value2",
		},
		{
			name:     "Empty map",
			url:      "https://example.com?",
			params:   map[string]string{},
			expected: "https://example.com?",
		},
		{
			name:     "URL without question mark",
			url:      "https://example.com",
			params:   map[string]string{"key1": "value1"},
			expected: "https://example.com?key1=value1",
		},
		{
			name:     "Special characters in parameters",
			url:      "https://example.com?",
			params:   map[string]string{"key": "value with spaces", "symbol": "&special"},
			expected: "https://example.com?key=value+with+spaces&symbol=%26special",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goutils.EncodeQueryParams(tt.url, tt.params)
			if got != tt.expected {
				t.Errorf("EncodeQueryParams() = %v, want %v", got, tt.expected)
			}
		})
	}
}

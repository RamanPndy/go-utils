package goutils

import "strings"

func AsBool(v any) (bool, bool) {
	b, ok := v.(bool)
	if ok {
		return b, true
	}
	if s, ok := v.(string); ok {
		s = strings.TrimSpace(strings.ToLower(s))
		if s == "true" {
			return true, true
		}
		if s == "false" {
			return false, true
		}
	}
	return false, false
}

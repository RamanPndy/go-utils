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

func GetBoolKeysListAndMap(payload map[string]any, keys ...string) bool {
	for _, key := range keys {
		if value, ok := payload[key]; ok {
			if b, ok := value.(bool); ok {
				return b
			}
		}
	}
	return false
}

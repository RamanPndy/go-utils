package goutils

import "strings"

func IndentBlock(text, prefix string) string {
	trimmed := strings.TrimRight(text, "\n")
	if trimmed == "" {
		return ""
	}
	lines := strings.Split(trimmed, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = prefix
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func ContainsFold(items []string, v string) bool {
	for _, item := range items {
		if strings.EqualFold(item, v) {
			return true
		}
	}
	return false
}

func AsStringSlice(v any) ([]string, bool) {
	items, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		s, ok := item.(string)
		if !ok {
			continue
		}
		out = append(out, s)
	}
	return out, true
}

func AsString(v any) (string, bool) {
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(s), true
}

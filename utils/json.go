package goutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// JsonToProperties converts simple, flat json, {"abc":"123","def":"456"} to a string-valued map.
// Field value types in json input are assumed to be of type string. If otherwise, an error is returned.
// A map value of nil is returned with empty json.
func JsonPropertiesToMap(jsonStr []byte) (map[string]string, error) {
	// Unmarshal the json string into map[string]interface{}
	var i interface{}
	var err error

	// Return nil on no properties
	if len(jsonStr) == 0 {
		return nil, nil
	}

	if err = json.Unmarshal(jsonStr, &i); err != nil {
		return nil, fmt.Errorf("properties unmarshal error: %s", err)
	}
	// Interpret this as a case of no properties
	if i == nil {
		return nil, nil
	}

	// Named property map but we don't know what the types are
	unknownTypesMap := i.(map[string]interface{})

	// Create map[string]string to return since package properties can only be of type string. If a property
	// type is a non-string, the package definition is incorrect.
	properties := make(map[string]string)
	for k, v := range unknownTypesMap {
		switch v := v.(type) {
		case string:
			properties[k] = v
		default:
			vt := reflect.TypeOf(v).String()
			return nil, fmt.Errorf("property, %s, invalid type, %v", k, vt)
		}
	}

	return properties, nil
}

func MustJSON(v any) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %s", err)
	}
	return b, nil
}

// EscapeJSONPointer escapes a string for use in a JSON Pointer (RFC 6901).
// The two special characters are ~ (→ ~0) and / (→ ~1).
func EscapeJSONPointer(s string) string {
	// Order matters: escape ~ first to avoid double-escaping.
	out := make([]byte, 0, len(s)+8)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '~':
			out = append(out, '~', '0')
		case '/':
			out = append(out, '~', '1')
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}

func UnescapeJSONPointer(s string) string { // Order matters: unescape ~1 first to avoid double-unescaping.
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '~' && i+1 < len(s) {
			switch s[i+1] {
			case '0':
				out = append(out, '~')
				i++
			case '1':
				out = append(out, '/')
				i++
			default:
				out = append(out, s[i])
			}
		} else {
			out = append(out, s[i])
		}
	}
	return string(out)
}

// Base64Encode encodes a string to base64.
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode decodes a base64-encoded string.
func Base64Decode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %s", err)
	}
	return string(b), nil
}

func JSONEqual(a, b any) (bool, error) {
	aj, err := json.Marshal(a)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	return string(aj) == string(bj), nil
}

func JSONDeepEqual(a, b any) (bool, error) {
	aj, err := json.Marshal(a)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	var am, bm any
	if err := json.Unmarshal(aj, &am); err != nil {
		return false, fmt.Errorf("json unmarshal error: %s", err)
	}
	if err := json.Unmarshal(bj, &bm); err != nil {
		return false, fmt.Errorf("json unmarshal error: %s", err)
	}
	return reflect.DeepEqual(am, bm), nil
}

func JSONContains(a, b any) (bool, error) {
	aj, err := json.Marshal(a)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false, fmt.Errorf("json marshal error: %s", err)
	}
	var am, bm any
	if err := json.Unmarshal(aj, &am); err != nil {
		return false, fmt.Errorf("json unmarshal error: %s", err)
	}
	if err := json.Unmarshal(bj, &bm); err != nil {
		return false, fmt.Errorf("json unmarshal error: %s", err)
	}
	return contains(am, bm), nil
}

func contains(a, b any) bool {
	switch av := a.(type) {
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok {
			return false
		}
		for k, v := range av {
			if !contains(v, bv[k]) {
				return false
			}
		}
		return true
	case []any:
		bv, ok := b.([]any)
		if !ok {
			return false
		}
		for _, av := range av {
			found := false
			for _, bv := range bv {
				if contains(av, bv) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

func JsonEncode(out any, r io.Reader, w io.Writer) error {
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func JsonDecode(s string, v any) error {
	if err := json.Unmarshal([]byte(s), v); err != nil {
		return fmt.Errorf("json unmarshal error: %s", err)
	}
	return nil
}

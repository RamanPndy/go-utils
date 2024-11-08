package goutils_test

import (
	"reflect"
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestJsonPropertiesToMap(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  []byte
		expected map[string]string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Empty JSON input",
			jsonStr:  []byte{},
			expected: nil,
			wantErr:  false,
		},
		{
			name:    "Valid JSON with string properties",
			jsonStr: []byte(`{"name": "John", "city": "New York"}`),
			expected: map[string]string{
				"name": "John",
				"city": "New York",
			},
			wantErr: false,
		},
		{
			name:     "JSON with non-string property",
			jsonStr:  []byte(`{"name": "John", "age": 30}`),
			expected: nil,
			wantErr:  true,
			errMsg:   "property, age, invalid type, float64",
		},
		{
			name:     "JSON with nested structure",
			jsonStr:  []byte(`{"name": "John", "address": {"city": "New York"}}`),
			expected: nil,
			wantErr:  true,
			errMsg:   "property, address, invalid type, map[string]interface {}",
		},
		{
			name:     "JSON with array",
			jsonStr:  []byte(`{"name": "John", "hobbies": ["reading", "traveling"]}`),
			expected: nil,
			wantErr:  true,
			errMsg:   "property, hobbies, invalid type, []interface {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := goutils.JsonPropertiesToMap(tt.jsonStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

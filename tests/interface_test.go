package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestIsNilInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Nil interface", nil, true},
		{"Non-nil int", 42, false},
		{"Nil pointer", (*int)(nil), true},
		{"Non-nil pointer", new(int), false},
		{"Nil map", map[string]int(nil), true},
		{"Empty map", map[string]int{}, false},
		{"Nil slice", []int(nil), true},
		{"Empty slice", []int{}, false},
		{"Nil chan", (chan int)(nil), true},
		{"Non-nil chan", make(chan int), false},
		{"Nil array", (*[3]int)(nil), true},
		{"Non-nil array", &[3]int{}, false},
		{"Non-nil struct", struct{}{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := goutils.IsNilInterface(test.input)
			if result != test.expected {
				t.Errorf("For %s: got %v, expected %v", test.name, result, test.expected)
			}
		})
	}
}

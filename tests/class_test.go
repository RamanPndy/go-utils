package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

type Base interface{}
type Derived struct{}
type Other struct{}

func TestIsSubclass(t *testing.T) {
	tests := []struct {
		name     string
		sub      interface{}
		base     interface{}
		expected bool
	}{
		// {
		// 	name:     "Derived implements Base",
		// 	sub:      &Derived{},
		// 	base:     (*Base)(nil),
		// 	expected: true,
		// },
		{
			name:     "int does not implement Base",
			sub:      1,
			base:     (*Base)(nil),
			expected: false,
		},
		{
			name:     "Other does not implement Base",
			sub:      &Other{},
			base:     (*Base)(nil),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goutils.IsSubclass(tt.sub, tt.base)
			if got != tt.expected {
				t.Errorf("IsSubclass(%v, %v) = %v; want %v", tt.sub, tt.base, got, tt.expected)
			}
		})
	}
}

// Test cases for IsInstance function
func TestIsInstance(t *testing.T) {
	type TestCase struct {
		name       string
		object     interface{}
		targetType interface{}
		expected   bool
	}

	testCases := []TestCase{
		{
			name:       "Match int type",
			object:     42,
			targetType: 0,
			expected:   true,
		},
		{
			name:       "Match string type",
			object:     "hello",
			targetType: "",
			expected:   true,
		},
		{
			name:       "Different int and string",
			object:     42,
			targetType: "",
			expected:   false,
		},
		{
			name:       "Different int types",
			object:     42,
			targetType: 100,
			expected:   true,
		},
		{
			name:       "Nil object",
			object:     nil,
			targetType: 0,
			expected:   false,
		},
		{
			name:       "Nil targetType",
			object:     42,
			targetType: nil,
			expected:   false,
		},
		{
			name:       "Nil object and targetType",
			object:     nil,
			targetType: nil,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := goutils.IsInstance(tc.object, tc.targetType)
			if result != tc.expected {
				t.Errorf("IsInstance(%v, %v) = %v; want %v", tc.object, tc.targetType, result, tc.expected)
			}
		})
	}
}

func TestHasAttr(t *testing.T) {
	// Test case with existing fields
	person := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

	// Test case 1: Check for existing fields
	if !goutils.HasAttr(person, "Name") {
		t.Errorf("Expected 'Name' field to be present in Person struct")
	}

	if !goutils.HasAttr(person, "Age") {
		t.Errorf("Expected 'Age' field to be present in Person struct")
	}

	if !goutils.HasAttr(person, "Email") {
		t.Errorf("Expected 'Email' field to be present in Person struct")
	}

	// Test case 2: Check for non-existing field
	if goutils.HasAttr(person, "Phone") {
		t.Errorf("Expected 'Phone' field to not be present in Person struct")
	}

	// Test case 3: Check with pointer
	personPtr := &person
	if !goutils.HasAttr(personPtr, "Name") {
		t.Errorf("Expected 'Name' field to be present in Person struct (pointer)")
	}

	// Test case 4: Check for non-struct type
	str := "this is a string"
	if goutils.HasAttr(str, "SomeField") {
		t.Errorf("Expected 'SomeField' to not be found in a non-struct type")
	}
}

func TestSetAttr(t *testing.T) {
	person := &Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

	// Test case 1: Set existing field "Name"
	err := goutils.SetAttr(person, "Name", "Bob")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if person.Name != "Bob" {
		t.Errorf("expected Name to be 'Bob', got '%s'", person.Name)
	}

	// Test case 2: Set existing field "Age"
	err = goutils.SetAttr(person, "Age", 35)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if person.Age != 35 {
		t.Errorf("expected Age to be 35, got %d", person.Age)
	}

	// Test case 3: Set field with wrong type
	err = goutils.SetAttr(person, "Age", "thirty-five")
	if err == nil {
		t.Errorf("expected error when setting field 'Age' with wrong type")
	}

	// Test case 4: Set non-existing field
	err = goutils.SetAttr(person, "Address", "New York")
	if err == nil {
		t.Errorf("expected error when setting non-existing field 'Address'")
	}
}

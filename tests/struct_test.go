package goutils_test

import (
	"fmt"
	"reflect"
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestGetFieldNames(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name: "Struct with string and int fields",
			input: Person{
				Name: "John",
				Age:  30,
				City: "123 Street",
			},
			expected: []string{"ID", "FirstName", "LastName", "Name", "Age", "Address", "Email", "City"},
		},
		{
			name: "Pointer to struct",
			input: &Person{
				Name: "Alice",
				Age:  25,
				City: "456 Road",
			},
			expected: []string{"ID", "FirstName", "LastName", "Name", "Age", "Address", "Email", "City"},
		},
		{
			name: "Struct with additional fields",
			input: struct {
				Title  string
				Salary int
				Active bool
			}{
				Title:  "Developer",
				Salary: 100000,
				Active: true,
			},
			expected: []string{"Title", "Salary", "Active"},
		},
	}

	// Run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the function
			got := goutils.GetStructFieldNames(test.input)

			// Compare the result with the expected value
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("For test %s, expected %v, but got %v", test.name, test.expected, got)
			}
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	// Create a Person instance
	person := Person{Name: "John", Age: 30, City: "123 Street"}

	tests := []struct {
		obj       interface{}
		fieldName string
		expected  interface{}
		expectNil bool
	}{
		{person, "Name", "John", false},
		{person, "Age", 30, false},
		{person, "City", "123 Street", false},
		{&person, "Name", "John", false},    // Test with pointer to struct
		{&person, "Age", 30, false},         // Test with pointer to struct
		{&person, "NonExistent", nil, true}, // Non-existent field
		{nil, "Name", nil, true},            // Nil object
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Field: %s", test.fieldName), func(t *testing.T) {
			result := goutils.GetStructFieldValue(test.obj, test.fieldName)
			if test.expectNil && result != nil {
				t.Errorf("Expected nil, got %v", result)
			}
			if !test.expectNil && result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMergeUniqueFields(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		obj1     Person
		obj2     Person
		expected Person
	}{
		{
			name:     "Merge with unique fields",
			obj1:     Person{Name: "John", Age: 30, City: ""},
			obj2:     Person{Name: "", Age: 25, City: "123 Street"},
			expected: Person{Name: "John", Age: 30, City: "123 Street"},
		},
		{
			name:     "Both objects have some fields set",
			obj1:     Person{Name: "Alice", Age: 0, City: "456 Road"},
			obj2:     Person{Name: "", Age: 40, City: ""},
			expected: Person{Name: "Alice", Age: 40, City: "456 Road"},
		},
		{
			name:     "No unique fields",
			obj1:     Person{Name: "", Age: 0, City: ""},
			obj2:     Person{Name: "", Age: 0, City: ""},
			expected: Person{Name: "", Age: 0, City: ""},
		},
		{
			name:     "All fields are unique",
			obj1:     Person{Name: "Bob", Age: 35, City: ""},
			obj2:     Person{Name: "", Age: 0, City: "789 Avenue"},
			expected: Person{Name: "Bob", Age: 35, City: "789 Avenue"},
		},
	}

	// Run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := goutils.MergeUniqueFields(&test.obj1, &test.obj2).(Person)

			// Compare the result with the expected value
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("For test %s, expected %+v, but got %+v", test.name, test.expected, got)
			}
		})
	}
}

func TestSkipMergeUniqueFields(t *testing.T) {
	tests := []struct {
		name     string
		obj1     *Person
		obj2     *Person
		expected Person
	}{
		{
			name: "Different fields, both non-zero",
			obj1: &Person{ID: 1, Name: "John", Age: 25, City: "New York"},
			obj2: &Person{ID: 2, Name: "Alice", Age: 30, City: "Los Angeles"},
			expected: Person{
				ID:   1,          // ID should remain as obj1
				Name: "John",     // Name from obj1
				Age:  25,         // Age from obj1
				City: "New York", // City from obj1
			},
		},
		{
			name: "Same field values in both structs",
			obj1: &Person{ID: 1, Name: "John", Age: 25, City: "New York"},
			obj2: &Person{ID: 2, Name: "John", Age: 25, City: "New York"},
			expected: Person{
				ID:   1,          // ID should remain as obj1
				Name: "John",     // Same Name, so keep from obj1
				Age:  25,         // Same Age, so keep from obj1
				City: "New York", // Same City, so keep from obj1
			},
		},
		{
			name: "Second object has unique values",
			obj1: &Person{ID: 1, Name: "", Age: 0, City: ""},
			obj2: &Person{ID: 2, Name: "Alice", Age: 30, City: "Los Angeles"},
			expected: Person{
				ID:   1,             // ID should remain as obj1
				Name: "Alice",       // Name from obj2
				Age:  30,            // Age from obj2
				City: "Los Angeles", // City from obj2
			},
		},
		{
			name: "First object has unique values",
			obj1: &Person{ID: 1, Name: "John", Age: 25, City: "New York"},
			obj2: &Person{ID: 2, Name: "", Age: 0, City: ""},
			expected: Person{
				ID:   1,          // ID should remain as obj1
				Name: "John",     // Name from obj1
				Age:  25,         // Age from obj1
				City: "New York", // City from obj1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := goutils.SkipMergeUniqueFields(tt.obj1, tt.obj2, []string{"ID"}).(Person)

			if !reflect.DeepEqual(merged, tt.expected) {
				t.Errorf("got %+v, want %+v", merged, tt.expected)
			}
		})
	}
}

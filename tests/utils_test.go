package goutils_test

import (
	"reflect"
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

// Define structs for testing
type Address struct {
	City  string
	State string
}

type Person struct {
	Name    string
	Age     int
	Address *Address
	Email   string
}

// Test case for the Zip function
func TestZip(t *testing.T) {
	tests := []struct {
		slice1 []int
		slice2 []int
		want   [][2]interface{}
	}{
		// Test with equal length slices
		{
			slice1: []int{1, 2, 3},
			slice2: []int{4, 5, 6},
			want:   [][2]interface{}{{1, 4}, {2, 5}, {3, 6}},
		},
		// Test with first slice longer than the second
		{
			slice1: []int{1, 2, 3, 4},
			slice2: []int{5, 6},
			want:   [][2]interface{}{{1, 5}, {2, 6}},
		},
		// Test with second slice longer than the first
		{
			slice1: []int{7, 8},
			slice2: []int{9, 10, 11},
			want:   [][2]interface{}{{7, 9}, {8, 10}},
		},
		// Test with empty slices
		{
			slice1: []int{},
			slice2: []int{},
			want:   [][2]interface{}{},
		},
		// Test with one empty slice
		{
			slice1: []int{1, 2, 3},
			slice2: []int{},
			want:   [][2]interface{}{},
		},
	}

	for _, tt := range tests {
		got := goutils.Zip(tt.slice1, tt.slice2)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Zip(%v, %v) = %v, want %v", tt.slice1, tt.slice2, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	// Define test cases
	tests := []struct {
		name  string
		slice interface{}
		item  interface{}
		want  bool
	}{
		{
			name:  "Test with int slice (item present)",
			slice: []int{1, 2, 3, 4, 5},
			item:  3,
			want:  true,
		},
		{
			name:  "Test with int slice (item not present)",
			slice: []int{1, 2, 3, 4, 5},
			item:  6,
			want:  false,
		},
		{
			name:  "Test with string slice (item present)",
			slice: []string{"apple", "banana", "cherry"},
			item:  "banana",
			want:  true,
		},
		{
			name:  "Test with string slice (item not present)",
			slice: []string{"apple", "banana", "cherry"},
			item:  "grape",
			want:  false,
		},
		{
			name:  "Test with empty slice",
			slice: []int{},
			item:  1,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch s := tt.slice.(type) {
			case []int:
				if got := goutils.Contains(s, tt.item.(int)); got != tt.want {
					t.Errorf("Contains(%v, %v) = %v, want %v", s, tt.item, got, tt.want)
				}
			case []string:
				if got := goutils.Contains(s, tt.item.(string)); got != tt.want {
					t.Errorf("Contains(%v, %v) = %v, want %v", s, tt.item, got, tt.want)
				}
			}
		})
	}
}

func TestCombineSlicesToMap(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		keys   interface{}
		values interface{}
		want   interface{}
	}{
		{
			name:   "Test with int keys and string values (equal length)",
			keys:   []int{1, 2, 3},
			values: []string{"apple", "banana", "cherry"},
			want:   map[int]string{1: "apple", 2: "banana", 3: "cherry"},
		},
		{
			name:   "Test with string keys and float values (equal length)",
			keys:   []string{"a", "b", "c"},
			values: []float64{1.1, 2.2, 3.3},
			want:   map[string]float64{"a": 1.1, "b": 2.2, "c": 3.3},
		},
		{
			name:   "Test with unequal length slices (keys longer)",
			keys:   []int{1, 2, 3, 4},
			values: []string{"apple", "banana"},
			want:   map[int]string{1: "apple", 2: "banana"},
		},
		{
			name:   "Test with unequal length slices (values longer)",
			keys:   []string{"x", "y"},
			values: []float64{10, 20, 30},
			want:   map[string]float64{"x": 10, "y": 20},
		},
		{
			name:   "Test with empty slices",
			keys:   []int{},
			values: []string{},
			want:   map[int]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch k := tt.keys.(type) {
			case []int:
				got := goutils.CombineSlicesToMap(k, tt.values.([]string))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CombineSlicesToMap() = %v, want %v", got, tt.want)
				}
			case []string:
				got := goutils.CombineSlicesToMap(k, tt.values.([]float64))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CombineSlicesToMap() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	// Define test cases
	tests := []struct {
		name  string
		input interface{}
		fn    interface{}
		want  interface{}
	}{
		{
			name:  "Test mapping integers to their squares",
			input: []int{1, 2, 3, 4, 5},
			fn: func(x int) int {
				return x * x
			},
			want: []int{1, 4, 9, 16, 25},
		},
		{
			name:  "Test mapping strings to their lengths",
			input: []string{"apple", "banana", "cherry"},
			fn: func(s string) int {
				return len(s)
			},
			want: []int{5, 6, 6},
		},
		{
			name:  "Test mapping integers to their doubles",
			input: []int{1, 2, 3},
			fn: func(x int) int {
				return x * 2
			},
			want: []int{2, 4, 6},
		},
		{
			name:  "Test mapping floats to their halves",
			input: []float64{2.0, 4.0, 6.0},
			fn: func(x float64) float64 {
				return x / 2.0
			},
			want: []float64{1.0, 2.0, 3.0},
		},
	}

	// Iterate over test cases and run the Map function for each
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				got := goutils.Map(input, tt.fn.(func(int) int))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Map() = %v, want %v", got, tt.want)
				}
			case []string:
				got := goutils.Map(input, tt.fn.(func(string) int))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Map() = %v, want %v", got, tt.want)
				}
			case []float64:
				got := goutils.Map(input, tt.fn.(func(float64) float64))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Map() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// Unit tests for the Filter function
func TestFilter(t *testing.T) {
	// Define test cases
	tests := []struct {
		name  string
		input interface{}
		fn    interface{}
		want  interface{}
	}{
		{
			name:  "Test filtering even numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			fn: func(x int) bool {
				return x%2 == 0
			},
			want: []int{2, 4, 6},
		},
		{
			name:  "Test filtering odd numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			fn: func(x int) bool {
				return x%2 != 0
			},
			want: []int{1, 3, 5},
		},
		{
			name:  "Test filtering strings longer than 4 characters",
			input: []string{"apple", "banana", "fig", "grape"},
			fn: func(s string) bool {
				return len(s) > 4
			},
			want: []string{"apple", "banana", "grape"},
		},
		{
			name:  "Test filtering floats greater than 2.0",
			input: []float64{1.5, 2.5, 3.0, 1.1, 5.5},
			fn: func(f float64) bool {
				return f > 2.0
			},
			want: []float64{2.5, 3.0, 5.5},
		},
	}

	// Iterate over test cases and run the Filter function for each
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				got := goutils.Filter(input, tt.fn.(func(int) bool))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Filter() = %v, want %v", got, tt.want)
				}
			case []string:
				got := goutils.Filter(input, tt.fn.(func(string) bool))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Filter() = %v, want %v", got, tt.want)
				}
			case []float64:
				got := goutils.Filter(input, tt.fn.(func(float64) bool))
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Filter() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDeepCopyJSON(t *testing.T) {
	p1 := Person{
		Name: "Alice",
		Age:  30,
		Address: &Address{
			City:  "Wonderland",
			State: "Imaginary",
		},
	}

	p2, err := goutils.DeepCopyJSON(p1)
	if err != nil {
		t.Fatalf("DeepCopyJSON() error = %v", err)
	}

	// Modify p2 to check if p1 remains unchanged
	p2.Name = "Bob"
	p2.Address.City = "Newland"

	if reflect.DeepEqual(p1, p2) {
		t.Errorf("DeepCopyJSON() = %v, want different from original %v", p2, p1)
	}

	// Ensure the copied object is different from the original
	if p1.Name == p2.Name || p1.Address.City == p2.Address.City {
		t.Errorf("DeepCopyJSON() failed to deep copy, got %v, want different from %v", p2, p1)
	}
}

func TestDeepCopyJSONNil(t *testing.T) {
	var p1 Person

	p2, err := goutils.DeepCopyJSON(p1)
	if err != nil {
		t.Fatalf("DeepCopyJSON() error = %v", err)
	}

	if !reflect.DeepEqual(p1, p2) {
		t.Errorf("DeepCopyJSON() = %v, want %v", p2, p1)
	}
}

// Unit tests for the Any function
func TestAny(t *testing.T) {
	tests := []struct {
		name      string
		slice     interface{}
		predicate interface{}
		want      bool
	}{
		{
			name:  "Has even number",
			slice: []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: true,
		},
		{
			name:  "No even number",
			slice: []int{1, 3, 5},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: false,
		},
		{
			name:  "Has long word",
			slice: []string{"apple", "banana", "cherry"},
			predicate: func(s string) bool {
				return len(s) > 5
			},
			want: true,
		},
		{
			name:  "No long word",
			slice: []string{"apple", "fig", "grape"},
			predicate: func(s string) bool {
				return len(s) > 5
			},
			want: false,
		},
		{
			name:  "Has number greater than 2",
			slice: []float64{1.1, 2.2, 3.3},
			predicate: func(f float64) bool {
				return f > 2.0
			},
			want: true,
		},
		{
			name:  "No number greater than 2",
			slice: []float64{1.1, 1.9},
			predicate: func(f float64) bool {
				return f > 2.0
			},
			want: false,
		},
		{
			name:  "Empty slice",
			slice: []int{},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: false,
		},
		{
			name:  "Nil slice",
			slice: []int(nil),
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch slice := tt.slice.(type) {
			case []int:
				got := goutils.Any(slice, tt.predicate.(func(int) bool))
				if got != tt.want {
					t.Errorf("Any() = %v, want %v", got, tt.want)
				}
			case []string:
				got := goutils.Any(slice, tt.predicate.(func(string) bool))
				if got != tt.want {
					t.Errorf("Any() = %v, want %v", got, tt.want)
				}
			case []float64:
				got := goutils.Any(slice, tt.predicate.(func(float64) bool))
				if got != tt.want {
					t.Errorf("Any() = %v, want %v", got, tt.want)
				}
			default:
				t.Fatalf("Unsupported slice type %T", slice)
			}
		})
	}
}

// Unit tests for the All function
func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		slice     interface{}
		predicate interface{}
		want      bool
	}{
		{
			name:  "All even numbers",
			slice: []int{2, 4, 6, 8},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: true,
		},
		{
			name:  "Not all even numbers",
			slice: []int{2, 4, 6, 7},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: false,
		},
		{
			name:  "All long words",
			slice: []string{"apple", "banana", "cherry"},
			predicate: func(s string) bool {
				return len(s) > 3
			},
			want: true,
		},
		{
			name:  "Not all long words",
			slice: []string{"apple", "fig", "cherry"},
			predicate: func(s string) bool {
				return len(s) > 3
			},
			want: false,
		},
		{
			name:  "All greater than zero",
			slice: []float64{1.1, 2.2, 3.3},
			predicate: func(f float64) bool {
				return f > 0
			},
			want: true,
		},
		{
			name:  "Not all greater than zero",
			slice: []float64{1.1, -2.2, 3.3},
			predicate: func(f float64) bool {
				return f > 0
			},
			want: false,
		},
		{
			name:  "Empty slice",
			slice: []int{},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: true,
		},
		{
			name:  "Nil slice",
			slice: []int(nil),
			predicate: func(n int) bool {
				return n%2 == 0
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch slice := tt.slice.(type) {
			case []int:
				got := goutils.All(slice, tt.predicate.(func(int) bool))
				if got != tt.want {
					t.Errorf("All() = %v, want %v", got, tt.want)
				}
			case []string:
				got := goutils.All(slice, tt.predicate.(func(string) bool))
				if got != tt.want {
					t.Errorf("All() = %v, want %v", got, tt.want)
				}
			case []float64:
				got := goutils.All(slice, tt.predicate.(func(float64) bool))
				if got != tt.want {
					t.Errorf("All() = %v, want %v", got, tt.want)
				}
			default:
				t.Fatalf("Unsupported slice type %T", slice)
			}
		})
	}
}

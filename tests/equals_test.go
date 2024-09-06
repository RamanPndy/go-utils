package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestEquals(t *testing.T) {
	// Test for integers
	if !goutils.Equals(1, 1) {
		t.Errorf("Equals(1, 1) = false; want true")
	}
	if goutils.Equals(1, 2) {
		t.Errorf("Equals(1, 2) = true; want false")
	}

	// Test for strings
	if !goutils.Equals("hello", "hello") {
		t.Errorf("Equals('hello', 'hello') = false; want true")
	}
	if goutils.Equals("hello", "world") {
		t.Errorf("Equals('hello', 'world') = true; want false")
	}

	// Test for structs
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Bob", Age: 25}

	if !goutils.Equals(p1, p2) {
		t.Errorf("Equals(p1, p2) = false; want true")
	}
	if goutils.Equals(p1, p3) {
		t.Errorf("Equals(p1, p3) = true; want false")
	}
}

func TestEqualsSlice(t *testing.T) {
	// Comparison function for integers
	intEquals := func(a, b int) bool { return a == b }
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}

	if !goutils.EqualsSlice(s1, s2, intEquals) {
		t.Errorf("EqualsSlice(s1, s2, intEquals) = false; want true")
	}
	if goutils.EqualsSlice(s1, s3, intEquals) {
		t.Errorf("EqualsSlice(s1, s3, intEquals) = true; want false")
	}

	// Comparison function for strings
	stringEquals := func(a, b string) bool { return a == b }
	s4 := []string{"a", "b"}
	s5 := []string{"a", "b"}
	s6 := []string{"a", "c"}

	if !goutils.EqualsSlice(s4, s5, stringEquals) {
		t.Errorf("EqualsSlice(s4, s5, stringEquals) = false; want true")
	}
	if goutils.EqualsSlice(s4, s6, stringEquals) {
		t.Errorf("EqualsSlice(s4, s6, stringEquals) = true; want false")
	}
}

func TestEqualsMap(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}

	if !goutils.EqualsMap(m1, m2) {
		t.Errorf("EqualsMap(m1, m2) = false; want true")
	}
	if goutils.EqualsMap(m1, m3) {
		t.Errorf("EqualsMap(m1, m3) = true; want false")
	}
}

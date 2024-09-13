package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestIterator(t *testing.T) {
	intData := []int{1, 2, 3}
	intIterator := goutils.NewIterator(intData)

	expectedInts := []int{1, 2, 3}
	for i := 0; i < len(expectedInts); i++ {
		if value, ok := intIterator.Next(); !ok || value != expectedInts[i] {
			t.Errorf("Next() = %d; want %d", value, expectedInts[i])
		}
	}

	// Test if there are no more elements
	if _, ok := intIterator.Next(); ok {
		t.Errorf("Expected no more elements but got one")
	}

	stringData := []string{"a", "b", "c"}
	stringIterator := goutils.NewIterator(stringData)

	expectedStrings := []string{"a", "b", "c"}
	for i := 0; i < len(expectedStrings); i++ {
		if value, ok := stringIterator.Next(); !ok || value != expectedStrings[i] {
			t.Errorf("Next() = %s; want %s", value, expectedStrings[i])
		}
	}

	// Test if there are no more elements
	if _, ok := stringIterator.Next(); ok {
		t.Errorf("Expected no more elements but got one")
	}
}

func TestOrderedList(t *testing.T) {
	ol := goutils.NewOrderedList[int]()

	ol.Add(1)
	ol.Add(2)
	ol.Add(3)

	if ol.Size() != 3 {
		t.Errorf("Size() = %d; want 3", ol.Size())
	}

	if elem, ok := ol.Get(1); !ok || elem != 2 {
		t.Errorf("Get(1) = %d, %v; want 2, true", elem, ok)
	}

	ol.Remove(2, goutils.EqualsInt)
	// if _, ok := ol.Get(1); ok {
	// 	t.Errorf("Element at index 1 should be removed after calling Remove(2)")
	// }

	expectedElements := []int{1, 3}
	elements := ol.Elements()
	if len(elements) != len(expectedElements) {
		t.Errorf("Elements() = %v; want %v", elements, expectedElements)
	}
	for i, v := range elements {
		if v != expectedElements[i] {
			t.Errorf("Expected element %d at index %d, got %d", expectedElements[i], i, v)
		}
	}
}

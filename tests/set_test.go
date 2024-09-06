package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestSet(t *testing.T) {
	set := goutils.NewSet[int]()
	set.Add(1)
	set.Add(2)

	if !set.Contains(1) {
		t.Errorf("Set should contain 1")
	}

	if !set.Contains(2) {
		t.Errorf("Set should contain 2")
	}

	if set.Contains(3) {
		t.Errorf("Set should not contain 3")
	}

	set.Remove(1)
	if set.Contains(1) {
		t.Errorf("Set should not contain 1 after removal")
	}

	if set.Size() != 1 {
		t.Errorf("Set size should be 1 after removing an element")
	}

	expectedElements := []int{2}
	elements := set.Elements()
	if len(elements) != len(expectedElements) {
		t.Errorf("Expected %v elements, got %v", expectedElements, elements)
	}
	for i, v := range elements {
		if v != expectedElements[i] {
			t.Errorf("Expected element %d at index %d, got %d", expectedElements[i], i, v)
		}
	}
}

func TestFrozenSet(t *testing.T) {
	fs := goutils.NewFrozenSet(1, 2, 3)

	if !fs.Contains(1) {
		t.Errorf("FrozenSet should contain 1")
	}

	if !fs.Contains(2) {
		t.Errorf("FrozenSet should contain 2")
	}

	if fs.Contains(4) {
		t.Errorf("FrozenSet should not contain 4")
	}

	if fs.Size() != 3 {
		t.Errorf("FrozenSet size should be 3")
	}

	expectedElements := []int{1, 2, 3}
	elements := fs.Elements()
	if len(elements) != len(expectedElements) {
		t.Errorf("Expected %v elements, got %v", expectedElements, elements)
	}
	for i, v := range elements {
		if v != expectedElements[i] {
			t.Errorf("Expected element %d at index %d, got %d", expectedElements[i], i, v)
		}
	}
}

func TestFrozenSetImmutable(t *testing.T) {
	fs := goutils.NewFrozenSet(1, 2, 3)
	elements := fs.Elements()
	elements[0] = 99 // Attempt to modify the elements slice

	// Check if the internal set is not modified
	if fs.Contains(99) {
		t.Errorf("FrozenSet should not contain 99 after attempting to modify elements")
	}
}

func TestOrderedSet(t *testing.T) {
	os := goutils.NewOrderedSet[int]()

	os.Add(1)
	os.Add(2)
	os.Add(3)
	os.Add(2) // Duplicate, should not be added

	if !os.Contains(1) {
		t.Errorf("OrderedSet should contain 1")
	}

	if !os.Contains(2) {
		t.Errorf("OrderedSet should contain 2")
	}

	if os.Contains(4) {
		t.Errorf("OrderedSet should not contain 4")
	}

	if os.Size() != 3 {
		t.Errorf("OrderedSet size should be 3")
	}

	os.Remove(2)
	if os.Contains(2) {
		t.Errorf("OrderedSet should not contain 2 after removal")
	}

	expectedElements := []int{1, 3}
	elements := os.Elements()
	if len(elements) != len(expectedElements) {
		t.Errorf("Expected %v elements, got %v", expectedElements, elements)
	}
	for i, v := range elements {
		if v != expectedElements[i] {
			t.Errorf("Expected element %d at index %d, got %d", expectedElements[i], i, v)
		}
	}
}

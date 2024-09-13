package goutils_test

import (
	"sort"
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestImmutableDict(t *testing.T) {
	entries := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	dict := goutils.NewImmutableDict(entries)

	if value, exists := dict.Get("b"); !exists || value != 2 {
		t.Errorf("Get('b') = %d, %v; want 2, true", value, exists)
	}

	if value, exists := dict.Get("d"); exists {
		t.Errorf("Get('d') = %d, %v; want 0, false", value, exists)
	}

	if dict.Size() != 3 {
		t.Errorf("Size() = %d; want 3", dict.Size())
	}

	expectedKeys := []string{"a", "b", "c"}
	keys := dict.Keys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() = %v; want %v", keys, expectedKeys)
	}
	sort.Strings(expectedKeys)
	sort.Strings(keys)
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %s at index %d, got %s", expectedKeys[i], i, k)
		}
	}

	expectedValues := []int{1, 2, 3}
	values := dict.Values()
	sort.Ints(expectedValues)
	sort.Ints(values)
	if len(values) != len(expectedValues) {
		t.Errorf("Values() = %v; want %v", values, expectedValues)
	}
	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %d at index %d, got %d", expectedValues[i], i, v)
		}
	}
}

func TestImmutableDictImmutable(t *testing.T) {
	entries := map[string]int{
		"a": 1,
		"b": 2,
	}
	dict := goutils.NewImmutableDict(entries)
	entries["c"] = 3 // Modify the original map

	// Check if the immutable dictionary is not affected
	if _, exists := dict.Get("c"); exists {
		t.Errorf("ImmutableDict should not contain key 'c' after modifying the original map")
	}
}

func TestOrderedDict(t *testing.T) {
	od := goutils.NewOrderedDict[string, int]()

	od.Set("a", 1)
	od.Set("b", 2)
	od.Set("c", 3)
	od.Set("b", 4) // Update value for existing key

	if value, exists := od.Get("b"); !exists || value != 4 {
		t.Errorf("Get('b') = %d, %v; want 4, true", value, exists)
	}

	if _, exists := od.Get("d"); exists {
		t.Errorf("Get('d') should not exist")
	}

	if od.Size() != 3 {
		t.Errorf("Size() = %d; want 3", od.Size())
	}

	od.Remove("b")
	if _, exists := od.Get("b"); exists {
		t.Errorf("Get('b') should not exist after removal")
	}

	expectedKeys := []string{"a", "c"}
	keys := od.Keys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() = %v; want %v", keys, expectedKeys)
	}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %s at index %d, got %s", expectedKeys[i], i, k)
		}
	}

	expectedValues := []int{1, 3}
	values := od.Values()
	if len(values) != len(expectedValues) {
		t.Errorf("Values() = %v; want %v", values, expectedValues)
	}
	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %d at index %d, got %d", expectedValues[i], i, v)
		}
	}

	expectedItems := []struct {
		Key   string
		Value int
	}{
		{"a", 1},
		{"c", 3},
	}
	items := od.Items()
	if len(items) != len(expectedItems) {
		t.Errorf("Items() = %v; want %v", items, expectedItems)
	}
	for i, item := range items {
		if item.Key != expectedItems[i].Key || item.Value != expectedItems[i].Value {
			t.Errorf("Expected item %+v at index %d, got %+v", expectedItems[i], i, item)
		}
	}
}

package goutils

// Iterator is a generic iterator for slices of any type T
type Iterator[T any] struct {
	data  []T
	index int
}

// NewIterator creates a new Iterator for a slice of any type T
func NewIterator[T any](data []T) *Iterator[T] {
	return &Iterator[T]{data: data, index: 0}
}

// HasNext checks if there are more elements in the iterator
func (it *Iterator[T]) HasNext() bool {
	return it.index < len(it.data)
}

// Next returns the next element and advances the iterator
func (it *Iterator[T]) Next() (T, bool) {
	if it.HasNext() {
		value := it.data[it.index]
		it.index++
		return value, true
	}
	var zeroValue T
	return zeroValue, false
}

// OrderedList represents an ordered list
type OrderedList[T any] struct {
	elements []T
}

// NewOrderedList creates a new OrderedList
func NewOrderedList[T any]() *OrderedList[T] {
	return &OrderedList[T]{}
}

// Add adds an element to the end of the OrderedList
func (ol *OrderedList[T]) Add(element T) {
	ol.elements = append(ol.elements, element)
}

// Remove removes an element from the OrderedList
func (ol *OrderedList[T]) Remove(element T, equals func(a, b T) bool) {
	for i, e := range ol.elements {
		if equals(e, element) {
			ol.elements = append(ol.elements[:i], ol.elements[i+1:]...)
			return
		}
	}
}

// Get retrieves an element by index
func (ol *OrderedList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(ol.elements) {
		var zero T
		return zero, false
	}
	return ol.elements[index], true
}

// Size returns the number of elements in the OrderedList
func (ol *OrderedList[T]) Size() int {
	return len(ol.elements)
}

// Elements returns a slice of all elements in the OrderedList
func (ol *OrderedList[T]) Elements() []T {
	return append([]T(nil), ol.elements...)
}

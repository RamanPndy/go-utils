package goutils

// Set is a simple implementation of a set using a map
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates a new Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

// Add adds an element to the set
func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

// Remove removes an element from the set
func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Contains checks if the set contains the element
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Elements returns a slice of elements in the set
func (s *Set[T]) Elements() []T {
	elems := make([]T, 0, len(s.elements))
	for elem := range s.elements {
		elems = append(elems, elem)
	}
	return elems
}

// FrozenSet is an immutable set implementation
type FrozenSet[T comparable] struct {
	elements map[T]struct{}
	list     []T
}

// NewFrozenSet creates a new FrozenSet with the given elements
func NewFrozenSet[T comparable](elements ...T) *FrozenSet[T] {
	set := &FrozenSet[T]{
		elements: make(map[T]struct{}),
		list:     make([]T, 0, len(elements)),
	}
	for _, elem := range elements {
		if _, exists := set.elements[elem]; !exists {
			set.elements[elem] = struct{}{}
			set.list = append(set.list, elem)
		}
	}
	return set
}

// Contains checks if the element is in the FrozenSet
func (fs *FrozenSet[T]) Contains(element T) bool {
	_, exists := fs.elements[element]
	return exists
}

// Size returns the number of elements in the FrozenSet
func (fs *FrozenSet[T]) Size() int {
	return len(fs.elements)
}

// Elements returns a slice of elements in the FrozenSet
func (fs *FrozenSet[T]) Elements() []T {
	// Return a copy of the list to avoid modification
	return append([]T(nil), fs.list...)
}

// OrderedSet is an ordered set implementation
type OrderedSet[T comparable] struct {
	elements map[T]struct{}
	order    []T
}

// NewOrderedSet creates a new OrderedSet
func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		elements: make(map[T]struct{}),
		order:    make([]T, 0),
	}
}

// Add adds an element to the OrderedSet
func (os *OrderedSet[T]) Add(element T) {
	if _, exists := os.elements[element]; !exists {
		os.elements[element] = struct{}{}
		os.order = append(os.order, element)
	}
}

// Remove removes an element from the OrderedSet
func (os *OrderedSet[T]) Remove(element T) {
	if _, exists := os.elements[element]; exists {
		delete(os.elements, element)
		for i, e := range os.order {
			if e == element {
				os.order = append(os.order[:i], os.order[i+1:]...)
				break
			}
		}
	}
}

// Contains checks if an element is in the OrderedSet
func (os *OrderedSet[T]) Contains(element T) bool {
	_, exists := os.elements[element]
	return exists
}

// Size returns the number of elements in the OrderedSet
func (os *OrderedSet[T]) Size() int {
	return len(os.elements)
}

// Elements returns a slice of elements in insertion order
func (os *OrderedSet[T]) Elements() []T {
	return append([]T(nil), os.order...)
}

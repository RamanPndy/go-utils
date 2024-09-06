package goutils

// ImmutableDict represents an immutable dictionary
type ImmutableDict[K comparable, V any] struct {
	data map[K]V
}

// NewImmutableDict creates a new ImmutableDict with the given key-value pairs
func NewImmutableDict[K comparable, V any](entries map[K]V) *ImmutableDict[K, V] {
	// Create a copy of the entries map to ensure immutability
	dict := make(map[K]V, len(entries))
	for k, v := range entries {
		dict[k] = v
	}
	return &ImmutableDict[K, V]{data: dict}
}

// Get retrieves the value for the given key
func (d *ImmutableDict[K, V]) Get(key K) (V, bool) {
	value, exists := d.data[key]
	return value, exists
}

// Size returns the number of entries in the ImmutableDict
func (d *ImmutableDict[K, V]) Size() int {
	return len(d.data)
}

// Keys returns a slice of all keys in the ImmutableDict
func (d *ImmutableDict[K, V]) Keys() []K {
	keys := make([]K, 0, len(d.data))
	for k := range d.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the ImmutableDict
func (d *ImmutableDict[K, V]) Values() []V {
	values := make([]V, 0, len(d.data))
	for _, v := range d.data {
		values = append(values, v)
	}
	return values
}

// OrderedDict represents an ordered dictionary
type OrderedDict[K comparable, V any] struct {
	data  map[K]V
	order []K
}

// NewOrderedDict creates a new OrderedDict
func NewOrderedDict[K comparable, V any]() *OrderedDict[K, V] {
	return &OrderedDict[K, V]{
		data:  make(map[K]V),
		order: make([]K, 0),
	}
}

// Set adds or updates the value for a given key
func (od *OrderedDict[K, V]) Set(key K, value V) {
	if _, exists := od.data[key]; !exists {
		od.order = append(od.order, key)
	}
	od.data[key] = value
}

// Get retrieves the value for the given key
func (od *OrderedDict[K, V]) Get(key K) (V, bool) {
	value, exists := od.data[key]
	return value, exists
}

// Remove removes a key-value pair from the OrderedDict
func (od *OrderedDict[K, V]) Remove(key K) {
	if _, exists := od.data[key]; exists {
		delete(od.data, key)
		for i, k := range od.order {
			if k == key {
				od.order = append(od.order[:i], od.order[i+1:]...)
				break
			}
		}
	}
}

// Size returns the number of items in the OrderedDict
func (od *OrderedDict[K, V]) Size() int {
	return len(od.data)
}

// Keys returns a slice of keys in the order they were inserted
func (od *OrderedDict[K, V]) Keys() []K {
	return append([]K(nil), od.order...)
}

// Values returns a slice of values in the order of the keys
func (od *OrderedDict[K, V]) Values() []V {
	values := make([]V, 0, len(od.data))
	for _, k := range od.order {
		values = append(values, od.data[k])
	}
	return values
}

// Items returns a slice of key-value pairs in the order of insertion
func (od *OrderedDict[K, V]) Items() []struct {
	Key   K
	Value V
} {
	items := make([]struct {
		Key   K
		Value V
	}, 0, len(od.data))
	for _, k := range od.order {
		items = append(items, struct {
			Key   K
			Value V
		}{Key: k, Value: od.data[k]})
	}
	return items
}

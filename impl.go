package main

import (
	"fmt"

	goutils "github.com/RamanPndy/go-utils/utils"
)

// Person struct with a deep copy method
type Person struct {
	ID        int
	FirstName string
	LastName  string
	Name      string
	Age       int
	Address   *Address
	Email     string
	City      string
}

// Address struct
type Address struct {
	City  string
	State string
}

func EqualsImpl() {
	// Working with integers
	fmt.Println(goutils.Equals(1, 1)) // true
	fmt.Println(goutils.Equals(1, 2)) // false

	// Working with strings
	fmt.Println(goutils.Equals("hello", "hello")) // true
	fmt.Println(goutils.Equals("hello", "world")) // false

	// Working with structs
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Bob", Age: 25}

	fmt.Println(goutils.Equals(p1, p2)) // true
	fmt.Println(goutils.Equals(p1, p3)) // false
}

func EqualsSliceImpl() {
	// Comparison function for integers
	intEquals := func(a, b int) bool { return a == b }
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}

	fmt.Println(goutils.EqualsSlice(s1, s2, intEquals)) // true
	fmt.Println(goutils.EqualsSlice(s1, s3, intEquals)) // false

	// Comparison function for strings
	stringEquals := func(a, b string) bool { return a == b }
	s4 := []string{"a", "b"}
	s5 := []string{"a", "b"}
	s6 := []string{"a", "c"}

	fmt.Println(goutils.EqualsSlice(s4, s5, stringEquals)) // true
	fmt.Println(goutils.EqualsSlice(s4, s6, stringEquals)) // false
}

func EqualsMapImpl() {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}

	fmt.Println(goutils.EqualsMap(m1, m2)) // true
	fmt.Println(goutils.EqualsMap(m1, m3)) // false
}

func ZipImpl() {
	slice1 := []int{1, 2, 3}
	slice2 := []string{"a", "b", "c"}

	result := goutils.Zip(slice1, slice2)
	fmt.Println(result) // Output: [[1 a] [2 b] [3 c]]
}

func ContainsImpl() {
	// Example usage with integers
	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Println(goutils.Contains(intSlice, 3)) // Output: true
	fmt.Println(goutils.Contains(intSlice, 6)) // Output: false

	// Example usage with strings
	stringSlice := []string{"apple", "banana", "cherry"}
	fmt.Println(goutils.Contains(stringSlice, "banana")) // Output: true
	fmt.Println(goutils.Contains(stringSlice, "grape"))  // Output: false
}

func CombineSlicesToMapImpl() {
	// Example usage with int keys and string values
	keys := []int{1, 2, 3}
	values := []string{"apple", "banana", "cherry"}
	combinedMap := goutils.CombineSlicesToMap(keys, values)
	fmt.Println(combinedMap) // map[1:apple 2:banana 3:cherry]

	// Example usage with string keys and float values
	keys2 := []string{"a", "b", "c"}
	values2 := []float64{1.1, 2.2, 3.3}
	combinedMap2 := goutils.CombineSlicesToMap(keys2, values2)
	fmt.Println(combinedMap2) // map[a:1.1 b:2.2 c:3.3]

	// Example with Unequal Lengths
	keys3 := []int{1, 2, 3, 4}
	values3 := []string{"apple", "banana"}
	combinedMap3 := goutils.CombineSlicesToMap(keys3, values3)
	fmt.Println(combinedMap3) // map[1:apple 2:banana]
}

func MapImpl() {
	// Example: Mapping integers to their squares
	input := []int{1, 2, 3, 4, 5}
	squareFn := func(x int) int {
		return x * x
	}
	result := goutils.Map(input, squareFn)
	fmt.Println("Squares:", result) // Squares: [1 4 9 16 25]

	// Example: Mapping strings to their lengths
	inputStrings := []string{"apple", "banana", "cherry"}
	lengthFn := func(s string) int {
		return len(s)
	}
	resultLengths := goutils.Map(inputStrings, lengthFn)
	fmt.Println("Lengths:", resultLengths) // Lengths: [5 6 6]
}

func FilterImpl() {
	// Example 1: Filter even numbers from a list of integers
	inputInts := []int{1, 2, 3, 4, 5, 6}
	isEven := func(x int) bool {
		return x%2 == 0
	}
	resultInts := goutils.Filter(inputInts, isEven)
	fmt.Println("Filtered even numbers:", resultInts) // Filtered even numbers: [2 4 6]

	// Example 2: Filter strings with length greater than 4
	inputStrings := []string{"apple", "banana", "fig", "grape"}
	longerThanFour := func(s string) bool {
		return len(s) > 4
	}
	resultStrings := goutils.Filter(inputStrings, longerThanFour)
	fmt.Println("Filtered strings longer than 4 characters:", resultStrings) // Filtered strings longer than 4 characters: [apple banana grape]

	// Example 3: Filter float numbers greater than 2.0
	inputFloats := []float64{1.5, 2.5, 3.0, 1.1, 5.5}
	greaterThanTwo := func(f float64) bool {
		return f > 2.0
	}
	resultFloats := goutils.Filter(inputFloats, greaterThanTwo)
	fmt.Println("Filtered floats greater than 2.0:", resultFloats) // Filtered floats greater than 2.0: [2.5 3 5.5]
}

func DeepCopyJSONImpl() {
	// Example usage with a struct
	type Person struct {
		Name    string
		Age     int
		Address struct {
			City  string
			State string
		}
	}

	p1 := Person{
		Name: "Alice",
		Age:  30,
	}
	p1.Address.City = "Wonderland"
	p1.Address.State = "Imaginary"

	p2, err := goutils.DeepCopyJSON(p1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Modify p2 to see if p1 remains unchanged
	p2.Name = "Bob"
	p2.Address.City = "Newland"

	fmt.Println("Original:", p1)
	fmt.Println("Copy:", p2)
}

// DeepCopy performs a manual deep copy of Person
func (p *Person) DeepCopy() *Person {
	if p == nil {
		return nil
	}
	// Create a new instance of Person
	copy := &Person{
		Name: p.Name,
		Age:  p.Age,
	}
	// Deep copy Address
	if p.Address != nil {
		copy.Address = &Address{
			City:  p.Address.City,
			State: p.Address.State,
		}
	}
	return copy
}

func ManualDeepCopyJSONImpl() {
	// Example usage
	p1 := &Person{
		Name: "Alice",
		Age:  30,
		Address: &Address{
			City:  "Wonderland",
			State: "Imaginary",
		},
	}

	p2 := p1.DeepCopy()

	// Modify p2 to see if p1 remains unchanged
	p2.Name = "Bob"
	p2.Address.City = "Newland"

	fmt.Println("Original:", *p1)
	fmt.Println("Copy:", *p2)
}

func AnyImpl() {
	// Example usage with integers
	numbers := []int{1, 2, 3, 4, 5}
	hasEven := goutils.Any(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("Has even number:", hasEven) // Output: Has even number: true

	// Example usage with strings
	words := []string{"apple", "banana", "cherry"}
	hasLongWord := goutils.Any(words, func(s string) bool {
		return len(s) > 5
	})
	fmt.Println("Has long word:", hasLongWord) // Output: Has long word: true

	// Example usage with floats
	floats := []float64{1.1, 2.2, 3.3}
	hasGreaterThan2 := goutils.Any(floats, func(f float64) bool {
		return f > 2.0
	})
	fmt.Println("Has number greater than 2:", hasGreaterThan2) // Output: Has number greater than 2: true
}

func AllImpl() {
	// Example usage with integers
	numbers := []int{2, 4, 6, 8}
	allEven := goutils.All(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("All numbers are even:", allEven) // Output: All numbers are even: true

	// Example usage with strings
	words := []string{"apple", "banana", "cherry"}
	allLongWords := goutils.All(words, func(s string) bool {
		return len(s) > 3
	})
	fmt.Println("All words are longer than 3 characters:", allLongWords) // Output: All words are longer than 3 characters: true

	// Example usage with floats
	floats := []float64{1.1, 2.2, 3.3}
	allGreaterThanZero := goutils.All(floats, func(f float64) bool {
		return f > 0
	})
	fmt.Println("All numbers are greater than zero:", allGreaterThanZero) // Output: All numbers are greater than zero: true
}

func IsSubclassImpl() {
	fmt.Println("Derived is a subclass of Base:", goutils.IsSubclass(&goutils.Derived{}, (*goutils.Base)(nil))) // true
	fmt.Println("int is a subclass of Base:", goutils.IsSubclass(1, (*goutils.Base)(nil)))                      // false
}

func IteratorImpl() {
	// Example with integers
	intData := []int{1, 2, 3, 4, 5}
	intIterator := goutils.NewIterator(intData)

	fmt.Println("Integer Iterator:")
	for intIterator.HasNext() {
		value, _ := intIterator.Next()
		fmt.Println(value)
	}

	// Example with strings
	stringData := []string{"a", "b", "c"}
	stringIterator := goutils.NewIterator(stringData)

	fmt.Println("String Iterator:")
	for stringIterator.HasNext() {
		value, _ := stringIterator.Next()
		fmt.Println(value)
	}
}

func SetImpl() {
	set := goutils.NewSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	fmt.Println("Set contains 2:", set.Contains(2)) // true
	fmt.Println("Set size:", set.Size())            // 3

	set.Remove(2)
	fmt.Println("Set contains 2:", set.Contains(2)) // false
	fmt.Println("Set size:", set.Size())            // 2

	elements := set.Elements()
	fmt.Println("Set elements:", elements) // [1 3] (order may vary)
}

func FrozenSetImpl() {
	fs := goutils.NewFrozenSet(1, 2, 3, 4, 5)

	fmt.Println("Contains 3:", fs.Contains(3)) // true
	fmt.Println("Contains 6:", fs.Contains(6)) // false
	fmt.Println("Size:", fs.Size())            // 5

	elements := fs.Elements()
	fmt.Println("Elements:", elements) // [1 2 3 4 5]
}

func OrderedSetImpl() {
	os := goutils.NewOrderedSet[int]()

	os.Add(1)
	os.Add(2)
	os.Add(3)
	os.Add(2) // Duplicate, should not be added

	fmt.Println("Contains 2:", os.Contains(2)) // true
	fmt.Println("Contains 4:", os.Contains(4)) // false
	fmt.Println("Size:", os.Size())            // 3

	os.Remove(2)
	fmt.Println("Contains 2 after removal:", os.Contains(2)) // false

	elements := os.Elements()
	fmt.Println("Elements:", elements) // [1 3]
}

func ImmutableDictImpl() {
	entries := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	dict := goutils.NewImmutableDict(entries)
	bval, _ := dict.Get("b")
	fmt.Println("Value for key 'b':", bval) // 2
	fmt.Println("Size:", dict.Size())       // 3

	keys := dict.Keys()
	fmt.Println("Keys:", keys) // [a b c] (order may vary)

	values := dict.Values()
	fmt.Println("Values:", values) // [1 2 3] (order may vary)
}

func OrderedDictImpl() {
	od := goutils.NewOrderedDict[string, int]()

	od.Set("a", 1)
	od.Set("b", 2)
	od.Set("c", 3)
	od.Set("b", 4) // Update value for existing key

	bval, _ := od.Get("b")
	fmt.Println("Value for key 'b':", bval) // 4
	fmt.Println("Size:", od.Size())         // 3

	od.Remove("b")
	bval, _ = od.Get("b")
	fmt.Println("Contains 'b' after removal:", bval) // false

	keys := od.Keys()
	fmt.Println("Keys:", keys) // [a c]

	values := od.Values()
	fmt.Println("Values:", values) // [1 3]

	items := od.Items()
	fmt.Println("Items:")
	for _, item := range items {
		fmt.Printf("Key: %s, Value: %d\n", item.Key, item.Value)
	}
}

func OrderedListImpl() {
	ol := goutils.NewOrderedList[int]()

	ol.Add(1)
	ol.Add(2)
	ol.Add(3)

	fmt.Println("Size:", ol.Size()) // 3

	if elem, ok := ol.Get(1); ok {
		fmt.Println("Element at index 1:", elem) // 2
	}

	ol.Remove(2, goutils.EqualsInt)
	if _, ok := ol.Get(1); !ok {
		fmt.Println("Element at index 1 has been removed")
	}

	elements := ol.Elements()
	fmt.Println("Elements:", elements) // [1 3]
}

func HasAttrImpl() {
	person := Person{Name: "Alice", Age: 30}

	// Test HasAttr with existing and non-existing fields
	fmt.Println(goutils.HasAttr(person, "Name"))  // true
	fmt.Println(goutils.HasAttr(person, "Age"))   // true
	fmt.Println(goutils.HasAttr(person, "Email")) // false
}

func SetAttrImpl() {
	person := &Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

	// Set the Name field
	if err := goutils.SetAttr(person, "Name", "Bob"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Updated Name:", person.Name)
	}

	// Set the Age field
	if err := goutils.SetAttr(person, "Age", 35); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Updated Age:", person.Age)
	}

	// Try setting an invalid field
	if err := goutils.SetAttr(person, "Address", "New York"); err != nil {
		fmt.Println(err)
	}
}

func IsInstanceImpl() {
	// Test examples
	var myInt int = 42
	var myString string = "hello"
	var anotherInt int = 100

	fmt.Println(goutils.IsInstance(myInt, 0))          // true (int)
	fmt.Println(goutils.IsInstance(myString, ""))      // true (string)
	fmt.Println(goutils.IsInstance(myInt, anotherInt)) // true (both int)
	fmt.Println(goutils.IsInstance(myInt, ""))         // false (int vs string)
	fmt.Println(goutils.IsInstance(myString, 0))       // false (string vs int)
}

func VarsImpl() {
	p := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	// Use Vars to get field names and values
	fields := goutils.Vars(p)

	// Print the map of field names and values
	for name, value := range fields {
		fmt.Printf("%s: %v\n", name, value)
	}
}

func MergeUniqueFieldsImpl() {
	// Create two objects of Person with different values
	person1 := &Person{Name: "John", Age: 30, City: ""}
	person2 := &Person{Name: "", Age: 25, City: "123 Street"}

	// Merge the objects
	mergedPerson := goutils.MergeUniqueFields(person1, person2).(Person)

	// Print the merged result
	fmt.Printf("Merged Person: %+v\n", mergedPerson) //Merged Person: {Name: John Age: 30 Address: 123 Street}
}

func SkipMergeUniqueFieldsImpl() {
	person1 := &Person{
		ID:   1,
		Name: "John",
		Age:  25,
		City: "New York",
	}

	person2 := &Person{
		ID:   2,
		Name: "Alice",
		Age:  30,
		City: "New York",
	}

	mergedPerson := goutils.SkipMergeUniqueFields(person1, person2, []string{"ID"}).(Person)
	fmt.Printf("Merged Person: %+v\n", mergedPerson) //Merged Person: {ID:1 Name:John Age:25 City:New York}
}

func GetStructFieldNamesImpl() {
	// Create an object of Person
	person := Person{Name: "John", Age: 30, City: "123 Street"}

	// Get field names using reflection
	fieldNames := goutils.GetStructFieldNames(person)

	// Print the field names
	fmt.Println("Field Names:", fieldNames) //Field Names: [Name Age Address]
}

func GetStructFieldValueImpl() {
	// Create an object of Person
	person := Person{Name: "John", Age: 30, City: "123 Street"}

	// Get field values using reflection
	name := goutils.GetStructFieldValue(person, "Name")
	age := goutils.GetStructFieldValue(person, "Age")
	address := goutils.GetStructFieldValue(person, "Address")

	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Address:", address)

	// Try to get a non-existent field
	nonExistent := goutils.GetStructFieldValue(person, "NonExistent")
	fmt.Println("NonExistent:", nonExistent) // Should print: Field 'NonExistent' not found in struct
}

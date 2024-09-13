package goutils

import (
	"fmt"
	"reflect"
)

type Base interface{}
type Derived struct{}

// IsSubclass checks if the type of 'sub' is a subtype of the type of 'base'.
func IsSubclass(sub, base interface{}) bool {
	subType := reflect.TypeOf(sub)
	baseType := reflect.TypeOf(base)

	// Check if baseType is an interface
	if baseType.Kind() == reflect.Interface {
		return subType.Implements(baseType)
	}

	// If baseType is a struct, check if subType is a pointer to a struct and check if it implements baseType
	if baseType.Kind() == reflect.Struct && subType.Kind() == reflect.Ptr {
		return reflect.TypeOf(sub).Elem().Implements(baseType)
	}

	return false
}

// IsInstance checks if the object is of a certain type
func IsInstance(object interface{}, targetType interface{}) bool {
	// Use reflect.TypeOf to get the dynamic type of the object and the targetType
	objectType := reflect.TypeOf(object)
	targetTypeType := reflect.TypeOf(targetType)

	// Compare the dynamic types
	return objectType == targetTypeType
}

// HasAttr checks if a struct has a specific field name
func HasAttr(obj interface{}, fieldName string) bool {
	// Get the reflection value of the object
	value := reflect.ValueOf(obj)

	// Check if the object is a pointer and dereference it
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// The object must be a struct to check for fields
	if value.Kind() != reflect.Struct {
		return false
	}

	// Check if the struct has the given field name
	if value.FieldByName(fieldName).IsValid() {
		return true
	}
	return false
}

// SetAttr sets the value of a field in a struct by its name
func SetAttr(obj interface{}, fieldName string, value interface{}) error {
	// Get the reflection value of the object (pointer to a struct)
	v := reflect.ValueOf(obj)

	// Ensure the object is a pointer
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("must pass a pointer to a struct")
	}

	// Dereference the pointer to get the actual struct
	v = v.Elem()

	// Ensure we have a struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct")
	}

	// Get the field by its name
	field := v.FieldByName(fieldName)

	// Ensure the field is valid and settable
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	// Get the reflection value of the new value
	val := reflect.ValueOf(value)

	// Ensure the type matches
	if field.Type() != val.Type() {
		return fmt.Errorf("provided value type (%s) does not match field type (%s)", val.Type(), field.Type())
	}

	// Set the field with the new value
	field.Set(val)
	return nil
}

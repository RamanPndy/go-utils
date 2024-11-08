package goutils

import (
	"fmt"
	"reflect"
)

// Function to merge two structs if they have unique values
func MergeUniqueFields(obj1, obj2 interface{}) interface{} {
	// Ensure that both objects are the same type
	v1 := reflect.ValueOf(obj1).Elem()
	v2 := reflect.ValueOf(obj2).Elem()

	// Create a new object of the same type
	merged := reflect.New(v1.Type()).Elem()

	// Iterate through the fields
	for i := 0; i < v1.NumField(); i++ {
		field1 := v1.Field(i)
		field2 := v2.Field(i)

		// If one field is zero value and the other is not, copy the non-zero field to the merged object
		if !field1.IsZero() && field2.IsZero() {
			merged.Field(i).Set(field1)
		} else if field1.IsZero() && !field2.IsZero() {
			merged.Field(i).Set(field2)
		}
	}
	return merged.Interface()
}

func SkipMergeUniqueFields(obj1, obj2 interface{}, skipFields []string) interface{} {
	// Ensure that both objects are the same type
	v1 := reflect.ValueOf(obj1).Elem()
	v2 := reflect.ValueOf(obj2).Elem()

	t1 := v1.Type()

	// Create a new instance of the struct
	result := reflect.New(t1).Elem()

	// Iterate through the fields
	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i)

		// Skip fields marked to be ignored
		if Contains(skipFields, field.Name) {
			result.Field(i).Set(v2.Field(i))
			continue
		}

		val1 := v1.Field(i).Interface()
		val2 := v2.Field(i).Interface()

		// If the values are different, prefer the unique value
		if val1 != val2 {
			result.Field(i).Set(v1.Field(i))
		}
	}
	return result.Interface()
}

func GetStructFieldNames(obj interface{}) []string {
	var fieldNames []string

	// Use reflect.ValueOf to get the reflection value of the object
	v := reflect.ValueOf(obj)

	// Ensure the object is passed as a pointer or struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get the type of the object
	t := v.Type()

	// Iterate through the struct fields
	for i := 0; i < v.NumField(); i++ {
		// Get the field name
		fieldName := t.Field(i).Name
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames
}

func GetStructFieldValue(obj interface{}, fieldName string) interface{} {
	if IsNilInterface(obj) {
		return nil
	}
	// Use reflect.ValueOf to get the reflection value of the object
	v := reflect.ValueOf(obj)

	// Ensure the object is passed as a pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get the field by name
	fieldVal := v.FieldByName(fieldName)

	// Check if the field exists and is valid
	if !fieldVal.IsValid() {
		fmt.Printf("Field '%s' not found in struct\n", fieldName)
		return nil
	}

	// Return the interface value of the field
	return fieldVal.Interface()
}

package goutils

import (
	"encoding/json"
)

// zip utility is similar to zip in python. it combines two slices of any type
func Zip[T, U any](slice1 []T, slice2 []U) [][2]interface{} {
	length := min(len(slice1), len(slice2))
	zipped := make([][2]interface{}, length)

	for i := 0; i < length; i++ {
		zipped[i] = [2]interface{}{slice1[i], slice2[i]}
	}
	return zipped
}

// Generic function to check if an item is present in a slice
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// CombineSlicesToMap combines two slices into a map where keys are from the first slice and values are from the second slice
func CombineSlicesToMap[K comparable, V any](keys []K, values []V) map[K]V {
	result := make(map[K]V)
	minLen := len(keys)

	if len(values) < len(keys) {
		minLen = len(values) // Ensure we don't exceed the length of the shorter slice
	}

	for i := 0; i < minLen; i++ {
		result[keys[i]] = values[i]
	}
	return result
}

// Map applies a given function to each element of the input slice and returns a new slice with the transformed values.
func Map[T any, R any](input []T, fn func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

// Filter applies a filtering function `fn` to each element of the input slice `input` and returns a new slice containing only the elements for which `fn` returns true.
func Filter[T any](input []T, fn func(T) bool) []T {
	var result []T
	for _, v := range input {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// DeepCopyJSON performs a deep copy of the input object using JSON serialization
func DeepCopyJSON[T any](src T) (T, error) {
	var copy T
	data, err := json.Marshal(src)
	if err != nil {
		return copy, err
	}
	err = json.Unmarshal(data, &copy)
	if err != nil {
		return copy, err
	}
	return copy, nil
}

// Any checks if any element in the slice satisfies the predicate function.
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All checks if all elements in the slice satisfy the predicate function.
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

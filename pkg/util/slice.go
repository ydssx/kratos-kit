package util

import (
	"slices"

	"github.com/samber/lo"
)

// AppendIfMissing appends an element to a slice if it is not already present.
//
// It takes a slice of type T and an element of type T as parameters.
// It returns a new slice of type T.
//
// Example:
//
//	slice := []int{1, 2, 3}
//	slice = AppendIfMissing(slice, 2) // [1, 2, 3]
//	slice = AppendIfMissing(slice, 4) // [1, 2, 3, 4]
func AppendIfMissing[T comparable](slice []T, elem T) []T {
	for _, s := range slice {
		if s == elem {
			return slice
		}
	}
	return append(slice, elem)
}

// Filter filters the elements of a slice based on a given function.
//
// Parameters:
//   - slice: The slice to filter.
//   - f: The function used to filter the elements.
//
// Returns:
//   - The filtered slice.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5}
//	filteredSlice := Filter(slice, func(x int) bool { return x%2 == 0 })
//	fmt.Println(filteredSlice) // Output: [2, 4]
func Filter[T comparable](slice []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, s := range slice {
		if f(s) {
			result = append(result, s)
		}
	}
	return result
}

// SliceUnion finds the union of multiple slices in Go.
//
// It takes in multiple slices of any comparable type as input and returns a single slice that contains all the unique elements from the input slices.
//
// The function uses a map to keep track of the unique elements and then creates a slice from the keys of the map.
//
// Parameters:
//   - slices: variadic parameter representing multiple slices of comparable type.
//
// Return type:
//   - []T: a single slice containing the union of all input slices.
//
// Example:
//
//	slice1 := []int{1, 2, 3, 4, 5}
//	slice2 := []int{5, 6, 7, 8, 9}
//	slice3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	union := SliceUnion(slice1, slice2, slice3)
//	fmt.Println(union) // Output: [1 2 3 4 5 6 7 8 9]
func SliceUnion[T comparable, Slice ~[]T](slices ...Slice) Slice {
	return lo.Union(slices...)
}

// SliceRemove removes elements from the given slice s.
func SliceRemove[T comparable, Slice ~[]T](s Slice, elems ...T) Slice {
	return lo.Without(s, elems...)
}

// Unique removes duplicate elements from the given slice s.
func Unique[T comparable, Slice ~[]T](s Slice) Slice {
	return lo.Uniq(s)
}

// SliceIntersect returns the intersection of two slices.
// It takes two slices a and b as parameters.
// It returns a new slice containing the elements that are present in both a and b.
// The two slices must be of comparable type T.
// It uses a map to keep track of the elements in slice a.
// Then it iterates through slice b and appends elements to the result if they are present in the map.
func SliceIntersect[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	for _, v := range a {
		m[v] = true
	}
	result := make([]T, 0, len(a))
	for _, v := range b {
		if m[v] {
			result = append(result, v)
		}
	}
	return result
}

// SliceEqualAny checks if all elements of slice a are present in slice b.
func SliceEqualAny[T comparable, Slice ~[]T](a, b Slice) bool {
	return slices.Equal(a, b)
}

// UniqueRandomSample returns a random sample of unique elements from the given slice.
func UniqueRandomSample[T comparable, Slice ~[]T](slice Slice, count int) Slice {
	return lo.Samples(slice, count)
}

// FlattenSlice 将多维切片平铺为一维切片
func FlattenSlice[T any, Slice ~[]T](slice []Slice) Slice {
	return lo.Flatten(slice)
}

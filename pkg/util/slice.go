package util

import (
	"math/rand"
	"reflect"
	"slices"
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
func SliceUnion[T comparable](slices ...[]T) []T {
	unionMap := make(map[T]bool)
	for _, slice := range slices {
		for _, elem := range slice {
			unionMap[elem] = true
		}
	}
	union := make([]T, 0, len(unionMap))
	for elem := range unionMap {
		union = append(union, elem)
	}
	return union
}

// SliceRemove 函数从给定的切片 s 中移除与参数 elems 中出现的元素相同的元素，并返回移除后的切片。
// 参数：
//   - s: 要进行元素移除操作的切片
//   - elems: 要移除的元素列表
//
// 返回值：
//   - 移除元素后的切片
//
// 示例：
//
//	s := []int{1, 2, 3, 4, 5}
//	result := SliceRemove(s, 2, 4)
//	fmt.Println(result) // Output: [1 3 5]
func SliceRemove[T comparable](s []T, elems ...T) []T {
	m := make(map[T]struct{})
	for _, elem := range elems {
		m[elem] = struct{}{}
	}
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, exists := m[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}

// Unique removes duplicate elements from the given slice s.
// It returns a new slice with all the unique elements from s in the original order.
// The element type T must be comparable.
// It uses a map to keep track of seen elements.
//
// Parameters:
//   - s: The slice to remove duplicates from.
//
// Returns:
//   - A new slice with all the unique elements from s in the original order.
//
// Example:
//
//	s := []int{1, 2, 3, 2, 1, 4}
//	result := Unique(s)
//	fmt.Println(result) // Output: [1 2 3 4]
func Unique[T comparable](s []T) []T {
	if len(s) <= 1 {
		return s
	}
	m := make(map[T]bool)
	result := make([]T, 0, len(s))
	for _, v := range s {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
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

func SliceEqualAny[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !slices.Contains(b, v) {
			return false
		}
	}
	return true
}

// UniqueRandomSample returns a random sample of unique elements from the given slice.
func UniqueRandomSample(slice []int, count int) []int {
	sample := []int{}
	indices := rand.Perm(len(slice))
	for i := 0; i < count && i < len(slice); i++ {
		sample = append(sample, slice[indices[i]])
	}
	return sample
}

// FlattenSlice 将多维切片平铺为一维切片
func FlattenSlice[T any](slice interface{}) []T {
	var result []T

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i).Interface()
			if reflect.TypeOf(item).Kind() == reflect.Slice {
				result = append(result, FlattenSlice[T](item)...)
			} else if v, ok := item.(T); ok {
				result = append(result, v)
			}
		}
	}

	return result
}

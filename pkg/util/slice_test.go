package util

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	// Test with an empty slice
	result := Unique([]int{})
	if len(result) != 0 {
		t.Errorf("Expected an empty slice, but got %v", result)
	}

	// Test with a slice containing duplicate elements
	result = Unique([]int{1, 2, 3, 2, 4, 1, 5})
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test with a slice containing no duplicate elements
	result = Unique([]int{1, 2, 3, 4, 5})
	expected = []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test with a slice of strings
	resultStr := Unique([]string{"apple", "banana", "orange", "banana", "kiwi"})
	expectedStr := []string{"apple", "banana", "orange", "kiwi"}
	if !reflect.DeepEqual(resultStr, expectedStr) {
		t.Errorf("Expected %v, but got %v", expectedStr, resultStr)
	}
}

func TestFlattenSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []int
	}{
		{
			name:     "一维整数切片",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "二维整数切片",
			input:    [][]int{{1, 2}, {3, 4}, {5}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "三维整数切片",
			input:    [][][]int{{{1, 2}, {3}}, {{4}, {5}}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "混合维度整数切片",
			input:    []interface{}{1, []int{2, 3}, [][]int{{4}, {5}}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "空切片",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "非切片类型",
			input:    42,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenSlice[int](tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlattenSlice() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}
package dtx

import (
	"reflect"
	"testing"
)

func TestUniqueSlice(t *testing.T) {
	type testCase[T comparable] struct {
		name    string
		slice   []T
		wantOut []T
	}

	t.Run("WithStrings", func(t *testing.T) {
		tests := []testCase[string]{
			{
				name:    "EmptySlice",
				slice:   []string{},
				wantOut: []string{},
			},
			{
				name:    "SingleElement",
				slice:   []string{"a"},
				wantOut: []string{"a"},
			},
			{
				name:    "AllUnique",
				slice:   []string{"a", "b", "c", "d"},
				wantOut: []string{"a", "b", "c", "d"},
			},
			{
				name:    "AllDuplicate",
				slice:   []string{"a", "a", "a", "a"},
				wantOut: []string{"a"},
			},
			{
				name:    "DuplicatesAtStart",
				slice:   []string{"a", "a", "b", "c"},
				wantOut: []string{"a", "b", "c"},
			},
			{
				name:    "DuplicatesAtEnd",
				slice:   []string{"a", "b", "c", "c"},
				wantOut: []string{"a", "b", "c"},
			},
			{
				name:    "DuplicatesInMiddle",
				slice:   []string{"a", "b", "b", "c"},
				wantOut: []string{"a", "b", "c"},
			},
			{
				name:    "DuplicatesNotAdjacent",
				slice:   []string{"a", "b", "a", "c", "b", "d"},
				wantOut: []string{"a", "b", "c", "d"},
			},
			{
				name:    "MixedDuplicates",
				slice:   []string{"x", "y", "x", "z", "y", "x"},
				wantOut: []string{"x", "y", "z"},
			},
			{
				name:    "TwoElements_BothSame",
				slice:   []string{"a", "a"},
				wantOut: []string{"a"},
			},
			{
				name:    "TwoElements_Different",
				slice:   []string{"a", "b"},
				wantOut: []string{"a", "b"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotOut := UniqueSlice(tt.slice); !reflect.DeepEqual(gotOut, tt.wantOut) {
					t.Errorf("UniqueSlice() = %v, want %v", gotOut, tt.wantOut)
				}
			})
		}
	})

	t.Run("WithIntegers", func(t *testing.T) {
		tests := []testCase[int]{
			{
				name:    "EmptySlice",
				slice:   []int{},
				wantOut: []int{},
			},
			{
				name:    "SingleElement",
				slice:   []int{42},
				wantOut: []int{42},
			},
			{
				name:    "AllUnique",
				slice:   []int{1, 2, 3, 4, 5},
				wantOut: []int{1, 2, 3, 4, 5},
			},
			{
				name:    "AllDuplicate",
				slice:   []int{7, 7, 7, 7},
				wantOut: []int{7},
			},
			{
				name:    "DuplicatesPreservesOrder",
				slice:   []int{5, 1, 5, 2, 1, 3},
				wantOut: []int{5, 1, 2, 3},
			},
			{
				name:    "NegativeNumbers",
				slice:   []int{-1, 2, -1, 3, 2},
				wantOut: []int{-1, 2, 3},
			},
			{
				name:    "ZeroValues",
				slice:   []int{0, 1, 0, 2, 0},
				wantOut: []int{0, 1, 2},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotOut := UniqueSlice(tt.slice); !reflect.DeepEqual(gotOut, tt.wantOut) {
					t.Errorf("UniqueSlice() = %v, want %v", gotOut, tt.wantOut)
				}
			})
		}
	})

	t.Run("WithBools", func(t *testing.T) {
		tests := []testCase[bool]{
			{
				name:    "EmptySlice",
				slice:   []bool{},
				wantOut: []bool{},
			},
			{
				name:    "AllTrue",
				slice:   []bool{true, true, true},
				wantOut: []bool{true},
			},
			{
				name:    "AllFalse",
				slice:   []bool{false, false, false},
				wantOut: []bool{false},
			},
			{
				name:    "MixedBools",
				slice:   []bool{true, false, true, false},
				wantOut: []bool{true, false},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotOut := UniqueSlice(tt.slice); !reflect.DeepEqual(gotOut, tt.wantOut) {
					t.Errorf("UniqueSlice() = %v, want %v", gotOut, tt.wantOut)
				}
			})
		}
	})
}

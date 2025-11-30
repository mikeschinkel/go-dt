package test

import (
	"reflect"
	"testing"

	"github.com/mikeschinkel/go-dt/dtx"
)

func TestIsNilable(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{
			name:  "pointer is nilable",
			value: (*int)(nil),
			want:  true,
		},
		{
			name:  "slice is nilable",
			value: []int(nil),
			want:  true,
		},
		{
			name:  "map is nilable",
			value: map[string]int(nil),
			want:  true,
		},
		{
			name:  "channel is nilable",
			value: (chan int)(nil),
			want:  true,
		},
		{
			name:  "function is nilable",
			value: (func())(nil),
			want:  true,
		},
		{
			name:  "nil interface value",
			value: (interface{})(nil),
			want:  false, // A nil interface{} itself is not nilable type when passed as any
		},
		{
			name:  "int is not nilable",
			value: 42,
			want:  false,
		},
		{
			name:  "string is not nilable",
			value: "test",
			want:  false,
		},
		{
			name:  "struct is not nilable",
			value: struct{}{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtx.IsNilable(tt.value)
			if got != tt.want {
				t.Errorf("IsNilable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{
			name:  "nil pointer",
			value: (*int)(nil),
			want:  true,
		},
		{
			name:  "non-nil pointer",
			value: new(int),
			want:  false,
		},
		{
			name:  "nil slice",
			value: []int(nil),
			want:  true,
		},
		{
			name:  "empty slice",
			value: []int{},
			want:  false,
		},
		{
			name:  "nil interface",
			value: (interface{})(nil),
			want:  true,
		},
		{
			name:  "non-nilable value",
			value: 42,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtx.IsNil(tt.value)
			if got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNilableKind(t *testing.T) {
	tests := []struct {
		name string
		kind reflect.Kind
		want bool
	}{
		{name: "Ptr", kind: reflect.Ptr, want: true},
		{name: "Slice", kind: reflect.Slice, want: true},
		{name: "Map", kind: reflect.Map, want: true},
		{name: "Chan", kind: reflect.Chan, want: true},
		{name: "Func", kind: reflect.Func, want: true},
		{name: "Interface", kind: reflect.Interface, want: true},
		{name: "Int", kind: reflect.Int, want: false},
		{name: "String", kind: reflect.String, want: false},
		{name: "Struct", kind: reflect.Struct, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtx.IsNilableKind(tt.kind)
			if got != tt.want {
				t.Errorf("IsNilableKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

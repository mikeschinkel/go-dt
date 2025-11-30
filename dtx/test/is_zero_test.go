package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt/dtx"
)

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{
			name:  "zero int",
			value: 0,
			want:  true,
		},
		{
			name:  "non-zero int",
			value: 42,
			want:  false,
		},
		{
			name:  "empty string",
			value: "",
			want:  true,
		},
		{
			name:  "non-empty string",
			value: "hello",
			want:  false,
		},
		{
			name:  "nil pointer",
			value: (*int)(nil),
			want:  true,
		},
		{
			name:  "zero struct",
			value: struct{}{},
			want:  true,
		},
		{
			name:  "nil slice",
			value: []int(nil),
			want:  true,
		},
		{
			name:  "empty non-nil slice",
			value: []int{},
			want:  false, // Empty non-nil slice is not zero; only nil slice is zero
		},
		{
			name:  "non-empty slice",
			value: []int{1, 2, 3},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtx.IsZero(tt.value)
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

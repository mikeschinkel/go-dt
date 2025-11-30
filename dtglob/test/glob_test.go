package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt/dtglob"
)

func TestGlobContains(t *testing.T) {
	tests := []struct {
		name   string
		glob   dtglob.Glob
		substr string
		want   bool
	}{
		{
			name:   "contains wildcard",
			glob:   dtglob.Glob("*.go"),
			substr: "*",
			want:   true,
		},
		{
			name:   "contains extension",
			glob:   dtglob.Glob("**/*.txt"),
			substr: ".txt",
			want:   true,
		},
		{
			name:   "does not contain",
			glob:   dtglob.Glob("*.go"),
			substr: ".py",
			want:   false,
		},
		{
			name:   "contains path separator",
			glob:   dtglob.Glob("src/**/*.go"),
			substr: "/",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.glob.Contains(tt.substr)
			if got != tt.want {
				t.Errorf("Glob.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobSplit(t *testing.T) {
	tests := []struct {
		name string
		glob dtglob.Glob
		sep  string
		want []string
	}{
		{
			name: "split by slash",
			glob: dtglob.Glob("src/test/*.go"),
			sep:  "/",
			want: []string{"src", "test", "*.go"},
		},
		{
			name: "split by double star",
			glob: dtglob.Glob("**/*.txt"),
			sep:  "/",
			want: []string{"**", "*.txt"},
		},
		{
			name: "no separator found",
			glob: dtglob.Glob("*.go"),
			sep:  "/",
			want: []string{"*.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.glob.Split(tt.sep)
			if len(got) != len(tt.want) {
				t.Errorf("Glob.Split() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Glob.Split()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

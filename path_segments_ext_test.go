package dt_test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestPathSegments_IsRoot(t *testing.T) {
	tests := []struct {
		name string
		pss  dt.PathSegments
		want bool
	}{
		{name: "empty is normalized root", pss: "", want: true},
		{name: "dot is filepath.Rel same-dir", pss: ".", want: true},
		{name: "single segment", pss: "a", want: false},
		{name: "nested segments", pss: "a/b", want: false},
		{name: "dot-dot is not root", pss: "..", want: false},
		{name: "dot-slash is not normalized", pss: "./", want: false},
		{name: "slash is not root segments", pss: "/", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pss.IsRoot(); got != tt.want {
				t.Errorf("IsRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

package dt_test

import (
	"errors"
	"runtime"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

type tildeDirPathTests []struct {
	name    string
	input   string
	want    dt.TildeDirPath
	wantErr error
}

func TestParse_TildeDirPath_Windows(t *testing.T) {
	tests := tildeDirPathTests{
		{name: "root tilde only", input: "~", want: "~"},
		{name: "tilde separator", input: "~\\", want: "~\\"},
		{name: "tilde separator alt", input: "~/", want: "~\\"},
		{name: "tilde nested", input: "~\\sub\\dir", want: "~\\sub\\dir"},
		{name: "tilde double separators", input: "~\\\\deep\\\\path", want: "~\\\\deep\\\\path"},
		{name: "tilde nested alt", input: "~/sub/dir", want: "~\\sub\\dir"},
		{name: "tilde alt separator", input: "~/sub", want: "~\\sub"},
		{name: "tilde missing separator", input: "~noslash", wantErr: dt.ErrNotTildePath},
		{name: "no tilde prefix", input: "C:\\tmp", wantErr: dt.ErrNotTildePath},
		{name: "empty", input: "", wantErr: dt.ErrEmpty},
	}
	if runtime.GOOS != "windows" {
		t.Skipf("Skipping non-Windows tests.")
	}
	runTildeDirPathTests(t, tests)
}

func TestParse_TildeDirPath_Nix(t *testing.T) {
	tests := tildeDirPathTests{
		{name: "root tilde only", input: "~", want: "~"},
		{name: "tilde separator", input: "~/", want: "~/"},
		{name: "tilde nested", input: "~/sub/dir", want: "~/sub/dir"},
		{name: "tilde double separators", input: "~//deep//path", want: "~//deep//path"},
		{name: "wrong separator", input: "~\\sub", wantErr: dt.ErrNotTildePath},
		{name: "tilde missing separator", input: "~noslash", wantErr: dt.ErrNotTildePath},
		{name: "no tilde prefix", input: "/tmp", wantErr: dt.ErrNotTildePath},
		{name: "empty", input: "", wantErr: dt.ErrEmpty},
	}

	if runtime.GOOS == "windows" {
		t.Skipf("Skipping non-Windows tests.")
	}
	runTildeDirPathTests(t, tests)
}

func runTildeDirPathTests(t *testing.T, tests tildeDirPathTests) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseTildeDirPath(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse dt.TildeDirPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Parse dt.TildeDirPath() gotTdp = %v, want %v", got, tt.want)
			}
		})
	}
}

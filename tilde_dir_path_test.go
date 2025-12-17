package dt_test

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParse_TildeDirPath(t *testing.T) {
	var tests []struct {
		name    string
		input   string
		wantTdp dt.TildeDirPath
		wantErr error
	}

	switch runtime.GOOS {
	case "windows":
		tests = []struct {
			name    string
			input   string
			wantTdp dt.TildeDirPath
			wantErr error
		}{
			{name: "root tilde only", input: "~", wantTdp: "~"},
			{name: "tilde separator", input: "~\\", wantTdp: "~\\"},
			{name: "tilde separator alt", input: "~/", wantTdp: "~\\"},
			{name: "tilde nested", input: "~\\sub\\dir", wantTdp: "~\\sub\\dir"},
			{name: "tilde double separators", input: "~\\\\deep\\\\path", wantTdp: "~\\\\deep\\\\path"},
			{name: "tilde nested alt", input: "~/sub/dir", wantTdp: "~\\sub\\dir"},
			{name: "tilde alt separator", input: "~/sub", wantTdp: "~\\sub"},
			{name: "tilde missing separator", input: "~noslash", wantErr: dt.ErrNotTildePath},
			{name: "no tilde prefix", input: "C:\\tmp", wantErr: dt.ErrNotTildePath},
			{name: "empty", input: "", wantErr: dt.ErrEmpty},
		}
	default:
		tests = []struct {
			name    string
			input   string
			wantTdp dt.TildeDirPath
			wantErr error
		}{
			{name: "root tilde only", input: "~", wantTdp: "~"},
			{name: "tilde separator", input: "~/", wantTdp: "~/"},
			{name: "tilde nested", input: "~/sub/dir", wantTdp: "~/sub/dir"},
			{name: "tilde double separators", input: "~//deep//path", wantTdp: "~//deep//path"},
			{name: "wrong separator", input: "~\\sub", wantErr: dt.ErrNotTildePath},
			{name: "tilde missing separator", input: "~noslash", wantErr: dt.ErrNotTildePath},
			{name: "no tilde prefix", input: "/tmp", wantErr: dt.ErrNotTildePath},
			{name: "empty", input: "", wantErr: dt.ErrEmpty},
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTdp, err := dt.ParseTildeDirPath(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse dt.TildeDirPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotTdp != tt.wantTdp {
				t.Fatalf("Parse dt.TildeDirPath() gotTdp = %v, want %v", gotTdp, tt.wantTdp)
			}
		})
	}
}

func Test_TildeDirPath_Expand(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	var tests []struct {
		name    string
		path    dt.TildeDirPath
		wantDp  dt.DirPath
		wantErr error
	}
	wd, _ := dt.Getwd()

	switch runtime.GOOS {
	case "windows":
		tests = []struct {
			name    string
			path    dt.TildeDirPath
			wantDp  dt.DirPath
			wantErr error
		}{
			{name: "current dir", path: ".", wantDp: wd},
			{name: "tilde only", path: "~", wantDp: dt.DirPath(home)},
			{name: "tilde separator", path: "~\\", wantDp: dt.DirPath(home)},
			{name: "tilde separator alt", path: "~/", wantDp: dt.DirPath(home)},
			{name: "tilde nested", path: "~\\sub\\dir", wantDp: dirPathJoin3(home, "sub", "dir")},
			{name: "tilde nested alt", path: "~/sub/dir", wantDp: dirPathJoin3(home, "sub", "dir")},
			{name: "tilde double separators", path: "~\\\\deep\\\\path", wantDp: dirPathJoin3(home, "deep", "path")},
			{name: "tilde alt separator", path: "~/sub", wantDp: dirPathJoin(home, "sub")},
			{name: "missing separator treated literal", path: "~noslash", wantDp: dirPathJoin(wd, "~noslash")},
			{name: "no tilde", path: "C:\\tmp", wantDp: "C:\\tmp"},
			{name: "empty", path: "", wantErr: dt.ErrEmpty},
		}
	default:
		tests = []struct {
			name    string
			path    dt.TildeDirPath
			wantDp  dt.DirPath
			wantErr error
		}{
			{name: "current dir", path: ".", wantDp: wd},
			{name: "tilde only", path: "~", wantDp: dt.DirPath(home)},
			{name: "tilde separator", path: "~/", wantDp: dt.DirPath(home)},
			{name: "tilde nested", path: "~/sub/dir", wantDp: dirPathJoin(home, "sub/dir")},
			{name: "tilde double separators", path: "~//deep//path", wantDp: dirPathJoin(home, "deep//path")},
			{name: "wrong separator treated literal", path: "~\\sub", wantDp: dirPathJoin(wd, "~\\sub")},
			{name: "missing separator treated literal", path: "~noslash", wantDp: dirPathJoin(wd, "~noslash")},
			{name: "no tilde", path: "/tmp", wantDp: dt.DirPath("/tmp")},
			{name: "empty", path: "", wantErr: dt.ErrEmpty},
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDp, err := tt.path.Expand()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Expand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotDp != tt.wantDp {
				t.Fatalf("Expand() gotDp = %v, want %v", gotDp, tt.wantDp)
			}
		})
	}
}

// These are so we don't rely on the funcs from the package. Why? We should not
// use the system under test to test the system under test.

// dirPathJoin joins two ~string into a Ddt.DirPath
func dirPathJoin[P ~string](part1, part2 P) dt.DirPath {
	return dt.DirPath(filepath.Clean(filepath.Join(string(part1), string(part2))))
}

// dirPathJoins joins three ~string into a Ddt.DirPath
func dirPathJoin3[P ~string](part1, part2, part3 P) dt.DirPath {
	return dt.DirPath(filepath.Clean(filepath.Join(string(part1), string(part2), string(part3))))
}

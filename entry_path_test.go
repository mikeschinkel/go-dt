package dt_test

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

type entryPathTests []struct {
	name    string
	path    dt.EntryPath
	want    dt.EntryPath
	wantErr error
}

func curDir() dt.EntryPath {
	wd, _ := dt.Getwd()
	return dt.EntryPath(wd)
}

func Test_EntryPath_Expand_Windows(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	tests := entryPathTests{
		{name: "current dir", path: ".", want: curDir()},
		{name: "tilde only", path: "~", want: dt.EntryPath(home)},
		{name: "tilde separator", path: "~\\", want: dt.EntryPath(home)},
		{name: "tilde separator alt", path: "~/", want: dt.EntryPath(home)},
		{name: "tilde nested", path: "~\\sub\\dir", want: entryPathJoin3(home, "sub", "dir")},
		{name: "tilde nested alt", path: "~/sub/dir", want: entryPathJoin3(home, "sub", "dir")},
		{name: "tilde double separators", path: "~\\\\deep\\\\path", want: entryPathJoin3(home, "deep", "path")},
		{name: "tilde alt separator", path: "~/sub", want: entryPathJoin(home, "sub")},
		{name: "missing separator treated literal", path: "~noslash", want: entryPathJoin(curDir(), "~noslash")},
		{name: "no tilde", path: "C:\\tmp", want: "C:\\tmp"},
		{name: "empty", path: "", wantErr: dt.ErrEmpty},
	}

	if runtime.GOOS != "windows" {
		t.Skipf("Skipping Windows tests.")
	}
	runEntryPathTests(t, tests)
}

func Test_EntryPath_Expand_Nix(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	tests := entryPathTests{
		{name: "current dir", path: ".", want: curDir()},
		{name: "tilde only", path: "~", want: dt.EntryPath(home)},
		{name: "tilde separator", path: "~/", want: dt.EntryPath(home)},
		{name: "tilde nested", path: "~/sub/dir", want: entryPathJoin(home, "sub/dir")},
		{name: "tilde double separators", path: "~//deep//path", want: entryPathJoin(home, "deep//path")},
		{name: "wrong separator treated literal", path: "~\\sub", want: entryPathJoin(curDir(), "~\\sub")},
		{name: "missing separator treated literal", path: "~noslash", want: entryPathJoin(curDir(), "~noslash")},
		{name: "no tilde", path: "/tmp", want: dt.EntryPath("/tmp")},
		{name: "empty", path: "", wantErr: dt.ErrEmpty},
	}

	if runtime.GOOS == "windows" {
		t.Skipf("Skipping Windows tests.")
	}
	runEntryPathTests(t, tests)
}

func runEntryPathTests(t *testing.T, tests entryPathTests) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.path.Expand()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Expand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Expand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// These are so we don't rely on the funcs from the package. Why? We should not
// use the system under test to test the system under test.

// entryPathJoin joins two ~string into an dt.EntryPath
func entryPathJoin[P ~string](part1, part2 P) dt.EntryPath {
	return dt.EntryPath(filepath.Clean(filepath.Join(string(part1), string(part2))))
}

// entryPathJoins joins three ~string into an dt.EntryPath
func entryPathJoin3[P ~string](part1, part2, part3 P) dt.EntryPath {
	return dt.EntryPath(filepath.Clean(filepath.Join(string(part1), string(part2), string(part3))))
}

package dt

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

type testDirTree struct {
	root      DirPath
	sub       DirPath
	nested    DirPath
	file1Path Filepath
	file2Path Filepath
}

func makeTestDirTree(t *testing.T) testDirTree {
	t.Helper()

	root := DirPath(t.TempDir())
	sub := DirPathJoin(root, "sub")
	nested := DirPathJoin(sub, "nested")

	if err := os.MkdirAll(string(nested), 0o755); err != nil {
		t.Fatalf("MkdirAll(nested) error = %v", err)
	}

	file1Path := FilepathJoin(root, "file1.txt")
	if err := os.WriteFile(string(file1Path), []byte("file1"), 0o644); err != nil {
		t.Fatalf("WriteFile(file1) error = %v", err)
	}

	file2Path := FilepathJoin(nested, "file2.txt")
	if err := os.WriteFile(string(file2Path), []byte("file2"), 0o644); err != nil {
		t.Fatalf("WriteFile(file2) error = %v", err)
	}

	return testDirTree{
		root:      root,
		sub:       sub,
		nested:    nested,
		file1Path: file1Path,
		file2Path: file2Path,
	}
}

func sortedStrings(ss []string) []string {
	cp := append([]string(nil), ss...)
	sort.Strings(cp)
	return cp
}

func collectWalkRels(t *testing.T, seq iter.Seq2[DirEntry, error]) (rels []string) {
	t.Helper()
	for de, err := range seq {
		if err != nil {
			t.Fatalf("walk err = %v", err)
		}
		rels = append(rels, filepath.ToSlash(string(de.Rel)))
	}
	return rels
}

func countSeq2[K any, V any](seq iter.Seq2[K, V]) (count int) {
	seq(func(_ K, _ V) bool {
		count++
		return true
	})
	return count
}

func TestDirPath_String(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want string
	}{
		{name: "empty", dp: "", want: ""},
		{name: "simple", dp: "a/b", want: "a/b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_Clean(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "cleans dot dot", dp: DirPath(filepath.Join("a", "b", "..", "c")), want: DirPath(filepath.Join("a", "c"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.Clean(); got != tt.want {
				t.Fatalf("Clean() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_Abs(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	tests := []struct {
		name string
		dp   DirPath
		want func() (DirPath, error)
	}{
		{name: "dot", dp: ".", want: func() (DirPath, error) { return DirPath(wd), nil }},
		{name: "relative", dp: DirPath(filepath.Join("a", "b")), want: func() (DirPath, error) {
			abs, err := filepath.Abs(filepath.Join("a", "b"))
			return DirPath(abs), err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := tt.want()
			got, err := tt.dp.Abs()
			if (err != nil) != (wantErr != nil) {
				t.Fatalf("Abs() error = %v, wantErr %v", err, wantErr)
			}
			if got != want {
				t.Fatalf("Abs() = %q, want %q", got, want)
			}
		})
	}
}

func TestDirPath_IsAbs(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	tests := []struct {
		name string
		dp   DirPath
		want bool
	}{
		{name: "relative", dp: "a/b", want: false},
		{name: "absolute", dp: DirPath(wd), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.IsAbs(); got != tt.want {
				t.Fatalf("IsAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirPath_Dir(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "simple", dp: DirPath(filepath.Join("a", "b", "c")), want: DirPath(filepath.Join("a", "b"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.Dir(); got != tt.want {
				t.Fatalf("Dir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_Base(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want PathSegment
	}{
		{name: "simple", dp: DirPath(filepath.Join("a", "b", "c")), want: "c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.Base(); got != tt.want {
				t.Fatalf("Base() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_Join(t *testing.T) {
	tests := []struct {
		name  string
		dp    DirPath
		elems []any
		want  DirPath
	}{
		{name: "joins elems", dp: DirPath(filepath.Join("a", "b")), elems: []any{"c", PathSegment("d")}, want: DirPath(filepath.Join("a", "b", "c", "d"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.Join(tt.elems...); got != tt.want {
				t.Fatalf("Join() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_Contains(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		part any
		want bool
	}{
		{name: "contains string", dp: DirPath(filepath.Join("a", "b", "c")), part: "b", want: true},
		{name: "not contains", dp: DirPath(filepath.Join("a", "b", "c")), part: "z", want: false},
		{name: "contains dirpath", dp: DirPath(filepath.Join("a", "b", "c")), part: DirPath("b"), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.Contains(tt.part); got != tt.want {
				t.Fatalf("Contains(%v) = %v, want %v", tt.part, got, tt.want)
			}
		})
	}
}

func TestDirPath_HasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		dp     DirPath
		prefix DirPath
		want   bool
	}{
		{name: "has prefix", dp: DirPath(filepath.Join("a", "b", "c")), prefix: DirPath(filepath.Join("a", "b")), want: true},
		{name: "missing prefix", dp: DirPath(filepath.Join("a", "b", "c")), prefix: "z", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.HasPrefix(tt.prefix); got != tt.want {
				t.Fatalf("HasPrefix(%q) = %v, want %v", tt.prefix, got, tt.want)
			}
		})
	}
}

func TestDirPath_HasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		dp     DirPath
		suffix string
		want   bool
	}{
		{name: "has suffix", dp: DirPath(filepath.Join("a", "b", "c")), suffix: "c", want: true},
		{name: "missing suffix", dp: DirPath(filepath.Join("a", "b", "c")), suffix: "z", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.HasSuffix(tt.suffix); got != tt.want {
				t.Fatalf("HasSuffix(%q) = %v, want %v", tt.suffix, got, tt.want)
			}
		})
	}
}

func TestDirPath_TrimPrefix(t *testing.T) {
	tests := []struct {
		name   string
		dp     DirPath
		prefix DirPath
		want   DirPath
	}{
		{
			name:   "trims",
			dp:     DirPath(filepath.Join("a", "b", "c")),
			prefix: DirPath(filepath.Join("a", "b")),
			want:   DirPath(strings.TrimPrefix(filepath.Join("a", "b", "c"), filepath.Join("a", "b"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.TrimPrefix(tt.prefix); got != tt.want {
				t.Fatalf("TrimPrefix(%q) = %q, want %q", tt.prefix, got, tt.want)
			}
		})
	}
}

func TestDirPath_TrimSuffix(t *testing.T) {
	tests := []struct {
		name       string
		dp         DirPath
		trimSuffix string
		want       DirPath
	}{
		{
			name:       "trims",
			dp:         DirPath(filepath.Join("a", "b", "c")),
			trimSuffix: "c",
			want:       DirPath(strings.TrimSuffix(filepath.Join("a", "b", "c"), "c")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.TrimSuffix(tt.trimSuffix); got != tt.want {
				t.Fatalf("TrimSuffix(%q) = %q, want %q", tt.trimSuffix, got, tt.want)
			}
		})
	}
}

func TestDirPath_EnsureTrailSep(t *testing.T) {
	sep := string(os.PathSeparator)
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "empty unchanged", dp: "", want: ""},
		{name: "adds sep", dp: DirPath(filepath.Join("a", "b")), want: DirPath(filepath.Join("a", "b") + sep)},
		{name: "already has sep", dp: DirPath("a" + sep), want: DirPath("a" + sep)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.EnsureTrailSep(); got != tt.want {
				t.Fatalf("EnsureTrailSep() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_HasDotDotPrefix(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want bool
	}{
		{name: "dotdot", dp: "..", want: true},
		{name: "dotdot segment", dp: DirPath(".."+string(os.PathSeparator)) + "a", want: true},
		{name: "not segment", dp: "..foo", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.HasDotDotPrefix(); got != tt.want {
				t.Fatalf("HasDotDotPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirPath_ToLower(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "lowercases", dp: "AbC", want: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.ToLower(); got != tt.want {
				t.Fatalf("ToLower() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_ToUpper(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "uppercases", dp: "AbC", want: "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.ToUpper(); got != tt.want {
				t.Fatalf("ToUpper() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_ToSlash(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want DirPath
	}{
		{name: "converts", dp: DirPath(filepath.Join("a", "b")), want: DirPath(filepath.ToSlash(filepath.Join("a", "b")))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.ToSlash(); got != tt.want {
				t.Fatalf("ToSlash() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_VolumeName(t *testing.T) {
	tests := []struct {
		name string
		dp   DirPath
		want VolumeName
	}{
		{name: "wrapper", dp: DirPath(filepath.Join("a", "b")), want: VolumeName(filepath.VolumeName(filepath.Join("a", "b")))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dp.VolumeName(); got != tt.want {
				t.Fatalf("VolumeName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_ParseDirPathAndParseDirPaths(t *testing.T) {
	tree := makeTestDirTree(t)

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	tests := []struct {
		name    string
		input   string
		wantDp  DirPath
		wantErr error
	}{
		{name: "empty", input: "", wantErr: ErrEmpty},
		{name: "literal tilde name", input: "~noslash", wantDp: "~noslash"},
		{name: "absolute path passthrough", input: string(tree.root), wantDp: tree.root},
		{name: "relative path passthrough", input: "relative/path", wantDp: "relative/path"},
		{name: "tilde only", input: "~", wantDp: DirPath(home)},
	}

	switch runtime.GOOS {
	case "windows":
		tests = append(tests, struct {
			name    string
			input   string
			wantDp  DirPath
			wantErr error
		}{name: "tilde nested slash", input: "~/sub/dir", wantDp: DirPathJoin(DirPath(home), filepath.Join("sub", "dir")).Clean()})
	default:
		tests = append(tests, struct {
			name    string
			input   string
			wantDp  DirPath
			wantErr error
		}{name: "tilde nested", input: "~/sub/dir", wantDp: DirPathJoin(DirPath(home), "sub/dir").Clean()})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDp, err := ParseDirPath(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseDirPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotDp != tt.wantDp {
				t.Fatalf("ParseDirPath() gotDp = %q, want %q", gotDp, tt.wantDp)
			}
		})
	}

	t.Run("ParseDirPaths collects errors", func(t *testing.T) {
		got, err := ParseDirPaths([]string{string(tree.root), ""})
		if err == nil {
			t.Fatalf("ParseDirPaths() error = nil, want non-nil")
		}
		if len(got) != 1 || got[0] != tree.root {
			t.Fatalf("ParseDirPaths() got = %v, want [%v]", got, tree.root)
		}
	})
}

func TestDirPath_Normalize(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	wd, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	abs := func(s string) DirPath {
		p, err := filepath.Abs(s)
		if err != nil {
			t.Fatalf("filepath.Abs(%q) error = %v", s, err)
		}
		return DirPath(p)
	}

	type testCase struct {
		name  string
		input DirPath
		want  DirPath
	}

	tests := []testCase{
		{name: "dot", input: ".", want: wd},
		{name: "tilde only", input: "~", want: DirPath(home)},
		{name: "tilde missing separator treated literal", input: "~noslash", want: abs("~noslash")},
		{name: "relative becomes absolute", input: "relative/path", want: abs("relative/path")},
	}

	switch runtime.GOOS {
	case "windows":
		tests = append(tests,
			testCase{name: "tilde nested backslash", input: "~\\sub\\dir", want: DirPathJoin(DirPath(home), filepath.Join("sub", "dir")).Clean()},
			testCase{name: "tilde nested slash", input: "~/sub/dir", want: DirPathJoin(DirPath(home), filepath.Join("sub", "dir")).Clean()},
			testCase{name: "absolute passthrough", input: "C:\\tmp", want: abs("C:\\tmp")},
		)
	default:
		tests = append(tests,
			testCase{name: "tilde nested", input: "~/sub/dir", want: DirPathJoin(DirPath(home), "sub/dir").Clean()},
			testCase{name: "wrong separator treated literal", input: "~\\sub", want: abs("~\\sub")},
			testCase{name: "absolute passthrough", input: "/tmp", want: abs("/tmp")},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Normalize()
			if err != nil {
				t.Fatalf("Normalize() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Normalize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirPath_EnsureExists(t *testing.T) {
	tests := []struct {
		name    string
		dp      func(t *testing.T, tree testDirTree) DirPath
		wantErr error
		wantDir bool
	}{
		{
			name: "creates when missing",
			dp: func(_ *testing.T, tree testDirTree) DirPath {
				return DirPathJoin(tree.root, "created")
			},
			wantDir: true,
		},
		{
			name: "errors on file",
			dp: func(_ *testing.T, tree testDirTree) DirPath {
				return DirPath(tree.file1Path)
			},
			wantErr: ErrPathIsFile,
		},
		{
			name: "no-op on existing dir",
			dp: func(_ *testing.T, tree testDirTree) DirPath {
				return tree.sub
			},
			wantDir: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			err := dp.EnsureExists()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("EnsureExists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			info, statErr := os.Stat(string(dp))
			if statErr != nil {
				t.Fatalf("Stat(%q) error = %v", dp, statErr)
			}
			if info.IsDir() != tt.wantDir {
				t.Fatalf("Stat(%q).IsDir = %v, want %v", dp, info.IsDir(), tt.wantDir)
			}
		})
	}
}

func TestDirPath_Exists(t *testing.T) {
	tests := []struct {
		name       string
		dp         func(t *testing.T, tree testDirTree) DirPath
		wantExists bool
		wantErr    bool
	}{
		{name: "existing dir", dp: func(_ *testing.T, tree testDirTree) DirPath { return tree.root }, wantExists: true},
		{name: "missing dir", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPathJoin(tree.root, "nope") }, wantExists: false},
		{name: "existing file", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPath(tree.file1Path) }, wantExists: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			got, err := dp.Exists()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.wantExists {
				t.Fatalf("Exists() = %v, want %v", got, tt.wantExists)
			}
		})
	}
}

func TestDirPath_Status(t *testing.T) {
	tests := []struct {
		name       string
		dp         func(t *testing.T, tree testDirTree) DirPath
		wantStatus EntryStatus
		wantErr    bool
	}{
		{name: "dir", dp: func(_ *testing.T, tree testDirTree) DirPath { return tree.root }, wantStatus: IsDirEntry},
		{name: "file", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPath(tree.file1Path) }, wantStatus: IsFileEntry},
		{name: "missing", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPathJoin(tree.root, "missing") }, wantStatus: IsMissingEntry},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			got, err := dp.Status()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.wantStatus {
				t.Fatalf("Status() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestDirPath_ReadDir(t *testing.T) {
	tests := []struct {
		name      string
		dp        func(t *testing.T, tree testDirTree) DirPath
		wantNames []string
		wantErr   bool
	}{
		{name: "root has file and sub", dp: func(_ *testing.T, tree testDirTree) DirPath { return tree.root }, wantNames: []string{"file1.txt", "sub"}},
		{name: "missing errors", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPathJoin(tree.root, "missing") }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			entries, err := dp.ReadDir()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			var names []string
			for _, e := range entries {
				names = append(names, e.Name())
			}
			if got, want := sortedStrings(names), sortedStrings(tt.wantNames); strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("ReadDir() names = %v, want %v", got, want)
			}
		})
	}
}

func TestDirPathRead(t *testing.T) {
	tests := []struct {
		name      string
		dp        func(t *testing.T, tree testDirTree) DirPath
		wantNames []string
		wantErr   bool
	}{
		{name: "root has file and sub", dp: func(_ *testing.T, tree testDirTree) DirPath { return tree.root }, wantNames: []string{"file1.txt", "sub"}},
		{name: "missing errors", dp: func(_ *testing.T, tree testDirTree) DirPath { return DirPathJoin(tree.root, "missing") }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			des, err := DirPathRead(dp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DirPathRead() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			var names []string
			for _, de := range des {
				names = append(names, string(de.Rel))
			}
			if got, want := sortedStrings(names), sortedStrings(tt.wantNames); strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("DirPathRead() names = %v, want %v", got, want)
			}
		})
	}
}

func TestDirPath_DirFS(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, tree testDirTree)
	}{
		{
			name: "can read file relative to fs root",
			run: func(t *testing.T, tree testDirTree) {
				fsys := tree.root.DirFS()
				b, err := fs.ReadFile(fsys, "file1.txt")
				if err != nil {
					t.Fatalf("ReadFile(DirFS, file1) error = %v", err)
				}
				if string(b) != "file1" {
					t.Fatalf("ReadFile(DirFS, file1) = %q, want %q", string(b), "file1")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			tt.run(t, tree)
		})
	}
}

func TestDirPath_Stat(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, tree testDirTree)
	}{
		{
			name: "Stat on OS path works",
			run: func(t *testing.T, tree testDirTree) {
				info, err := tree.root.Stat()
				if err != nil {
					t.Fatalf("Stat(os, root) error = %v", err)
				}
				if !info.IsDir() {
					t.Fatalf("Stat(os, root).IsDir = false, want true")
				}
			},
		},
		{
			name: "Stat with DirFS uses relative path",
			run: func(t *testing.T, tree testDirTree) {
				fsys := tree.root.DirFS()
				info, err := DirPath(".").Stat(fsys)
				if err != nil {
					t.Fatalf("Stat(DirFS, .) error = %v", err)
				}
				if !info.IsDir() {
					t.Fatalf("Stat(DirFS, .).IsDir = false, want true")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			tt.run(t, tree)
		})
	}
}

func TestDirPath_Rel(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, tree testDirTree)
	}{
		{
			name: "nested rel to root",
			run: func(t *testing.T, tree testDirTree) {
				got, err := tree.nested.Rel(tree.root)
				if err != nil {
					t.Fatalf("Rel(nested, root) error = %v", err)
				}
				want := PathSegments(filepath.Join("sub", "nested"))
				if got != want {
					t.Fatalf("Rel(nested, root) = %q, want %q", got, want)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			tt.run(t, tree)
		})
	}
}

func TestDirPath_CanWrite(t *testing.T) {
	tests := []struct {
		name string
		dp   func(t *testing.T, tree testDirTree) DirPath
		want bool
	}{
		{name: "temp dir writable", dp: func(_ *testing.T, tree testDirTree) DirPath { return tree.root }, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			dp := tt.dp(t, tree)
			got, err := dp.CanWrite()
			if err != nil {
				t.Fatalf("CanWrite() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CanWrite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirPath_Walk(t *testing.T) {
	tests := []struct {
		name     string
		wantRels []string
	}{
		{
			name: "yields all entries",
			wantRels: []string{
				"file1.txt",
				"sub",
				"sub/nested",
				"sub/nested/file2.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			rels := collectWalkRels(t, tree.root.Walk())
			if got, want := sortedStrings(rels), sortedStrings(tt.wantRels); strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("Walk() rels = %v, want %v", got, want)
			}
		})
	}
}

func TestDirPath_WalkFiles(t *testing.T) {
	tests := []struct {
		name     string
		wantRels []string
		wantCnt  int
	}{
		{name: "yields only files", wantRels: []string{"file1.txt", "sub/nested/file2.txt"}, wantCnt: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			rels := collectWalkRels(t, tree.root.WalkFiles())
			if got, want := sortedStrings(rels), sortedStrings(tt.wantRels); strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("WalkFiles() rels = %v, want %v", got, want)
			}
			if cnt := countSeq2(tree.root.WalkFiles()); cnt != tt.wantCnt {
				t.Fatalf("WalkFiles() count = %d, want %d", cnt, tt.wantCnt)
			}
		})
	}
}

func TestDirPath_WalkDirs(t *testing.T) {
	tests := []struct {
		name     string
		wantRels []string
		wantCnt  int
	}{
		{name: "yields only dirs", wantRels: []string{"sub", "sub/nested"}, wantCnt: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			rels := collectWalkRels(t, tree.root.WalkDirs())
			if got, want := sortedStrings(rels), sortedStrings(tt.wantRels); strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("WalkDirs() rels = %v, want %v", got, want)
			}
			if cnt := countSeq2(tree.root.WalkDirs()); cnt != tt.wantCnt {
				t.Fatalf("WalkDirs() count = %d, want %d", cnt, tt.wantCnt)
			}
		})
	}
}

func TestDirPath_WalkFSStopsAfterBreak(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "stops after first yield", want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := makeTestDirTree(t)
			fsys := tree.root.DirFS()
			count := 0
			tree.root.WalkFS(fsys)(func(_ DirEntry, _ error) bool {
				count++
				return false
			})
			if count != tt.want {
				t.Fatalf("WalkFS stop count = %d, want %d", count, tt.want)
			}
		})
	}
}

func TestDirPath_ToTildeRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	tests := []struct {
		name string
		dp   DirPath
	}{
		{name: "home", dp: DirPath(home)},
		{name: "nested", dp: DirPathJoin(DirPath(home), filepath.Join("sub", "dir")).Clean()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdp, err := tt.dp.ToTilde()
			if err != nil {
				t.Fatalf("ToTilde() error = %v", err)
			}
			got, err := tdp.Expand()
			if err != nil {
				t.Fatalf("Expand(ToTilde()) error = %v", err)
			}
			want, err := tt.dp.Clean().Abs()
			if err != nil {
				t.Fatalf("Abs() error = %v", err)
			}
			if got != want {
				t.Fatalf("Expand(ToTilde()) = %q, want %q", got, want)
			}
		})
	}
}

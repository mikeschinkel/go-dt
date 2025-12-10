package dt

import (
	"reflect"
	"testing"
)

// Dummy types for testing generic functions
type TestSegment string
type TestSegments []TestSegment

func TestSplitSegments(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantOut []TestSegment
	}{
		{
			name:    "empty string",
			s:       "",
			wantOut: []TestSegment{},
		},
		{
			name:    "single segment",
			s:       "a",
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "single longer segment",
			s:       "alpha",
			wantOut: []TestSegment{"alpha"},
		},
		{
			name:    "two segments",
			s:       "a/b",
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "two longer segments",
			s:       "alpha/beta",
			wantOut: []TestSegment{"alpha", "beta"},
		},
		{
			name:    "three segments",
			s:       "a/b/c",
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "three longer segments",
			s:       "alpha/beta/gamma",
			wantOut: []TestSegment{"alpha", "beta", "gamma"},
		},
		{
			name:    "four segments",
			s:       "a/b/c/d",
			wantOut: []TestSegment{"a", "b", "c", "d"},
		},
		{
			name:    "six segments",
			s:       "a/b/c/d/e/f",
			wantOut: []TestSegment{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:    "deep path ten segments",
			s:       "a/b/c/d/e/f/g/h/i/j",
			wantOut: []TestSegment{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			name:    "segments with hyphens",
			s:       "foo-bar/baz-qux/quux-corge",
			wantOut: []TestSegment{"foo-bar", "baz-qux", "quux-corge"},
		},
		{
			name:    "segments with underscores",
			s:       "foo_bar/baz_qux",
			wantOut: []TestSegment{"foo_bar", "baz_qux"},
		},
		{
			name:    "segments with numbers",
			s:       "path1/path2/path3",
			wantOut: []TestSegment{"path1", "path2", "path3"},
		},
		{
			name:    "segments with dots",
			s:       "file.txt/dir.name/file.ext",
			wantOut: []TestSegment{"file.txt", "dir.name", "file.ext"},
		},
		{
			name:    "consecutive slashes two segments",
			s:       "a//b",
			wantOut: []TestSegment{"a", "", "b"},
		},
		{
			name:    "consecutive slashes three segments",
			s:       "a///b",
			wantOut: []TestSegment{"a", "", "", "b"},
		},
		{
			name:    "leading slash",
			s:       "/a/b",
			wantOut: []TestSegment{"", "a", "b"},
		},
		{
			name:    "domain path",
			s:       "api/v1/users",
			wantOut: []TestSegment{"api", "v1", "users"},
		},
		{
			name:    "api endpoint",
			s:       "api/v2/users/123",
			wantOut: []TestSegment{"api", "v2", "users", "123"},
		},
		{
			name:    "mixed case",
			s:       "MyPath/MyFile/MyExt",
			wantOut: []TestSegment{"MyPath", "MyFile", "MyExt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := SplitSegments[TestSegment](tt.s, "/")
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("SplitSegments() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestIndexSegments(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		index   int
		wantSeg TestSegment
	}{
		{
			name:    "empty segments with index 0",
			s:       "",
			index:   0,
			wantSeg: TestSegment(""),
		},
		{
			name:    "empty segments with index 1",
			s:       "",
			index:   1,
			wantSeg: TestSegment(""),
		},
		{
			name:    "empty segments with negative index",
			s:       "",
			index:   -1,
			wantSeg: TestSegment(""),
		},
		{
			name:    "single segment with index 0",
			s:       "a",
			index:   0,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "single segment with index 1",
			s:       "a",
			index:   1,
			wantSeg: TestSegment(""),
		},
		{
			name:    "two segments with index 0",
			s:       "a/b",
			index:   0,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "two segments with index 1",
			s:       "a/b",
			index:   1,
			wantSeg: TestSegment("b"),
		},
		{
			name:    "two segments with index 2",
			s:       "a/b",
			index:   2,
			wantSeg: TestSegment(""),
		},
		{
			name:    "two segments with negative index",
			s:       "a/b",
			index:   -1,
			wantSeg: TestSegment(""),
		},
		{
			name:    "three segments with index 0",
			s:       "a/b/c",
			index:   0,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "three segments with index 1",
			s:       "a/b/c",
			index:   1,
			wantSeg: TestSegment("b"),
		},
		{
			name:    "three segments with index 2",
			s:       "a/b/c",
			index:   2,
			wantSeg: TestSegment("c"),
		},
		{
			name:    "three segments with index 3",
			s:       "a/b/c",
			index:   3,
			wantSeg: TestSegment(""),
		},
		{
			name:    "longer segment names index 0",
			s:       "alpha/beta/gamma",
			index:   0,
			wantSeg: TestSegment("alpha"),
		},
		{
			name:    "longer segment names index 1",
			s:       "alpha/beta/gamma",
			index:   1,
			wantSeg: TestSegment("beta"),
		},
		{
			name:    "longer segment names index 2",
			s:       "alpha/beta/gamma",
			index:   2,
			wantSeg: TestSegment("gamma"),
		},
		{
			name:    "deep path index 0",
			s:       "a/b/c/d/e/f",
			index:   0,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "deep path index 3",
			s:       "a/b/c/d/e/f",
			index:   3,
			wantSeg: TestSegment("d"),
		},
		{
			name:    "deep path index 5",
			s:       "a/b/c/d/e/f",
			index:   5,
			wantSeg: TestSegment("f"),
		},
		{
			name:    "segments with hyphens index 0",
			s:       "foo-bar/baz-qux",
			index:   0,
			wantSeg: TestSegment("foo-bar"),
		},
		{
			name:    "segments with hyphens index 1",
			s:       "foo-bar/baz-qux",
			index:   1,
			wantSeg: TestSegment("baz-qux"),
		},
		{
			name:    "segments with underscores index 0",
			s:       "foo_bar/baz_qux",
			index:   0,
			wantSeg: TestSegment("foo_bar"),
		},
		{
			name:    "segments with underscores index 1",
			s:       "foo_bar/baz_qux",
			index:   1,
			wantSeg: TestSegment("baz_qux"),
		},
		{
			name:    "segments with numbers index 0",
			s:       "path1/path2",
			index:   0,
			wantSeg: TestSegment("path1"),
		},
		{
			name:    "segments with numbers index 1",
			s:       "path1/path2",
			index:   1,
			wantSeg: TestSegment("path2"),
		},
		{
			name:    "consecutive slashes index 0",
			s:       "a//b",
			index:   0,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "consecutive slashes index 1",
			s:       "a//b",
			index:   1,
			wantSeg: TestSegment(""),
		},
		{
			name:    "consecutive slashes index 2",
			s:       "a//b",
			index:   2,
			wantSeg: TestSegment("b"),
		},
		{
			name:    "leading slash index 0",
			s:       "/a/b",
			index:   0,
			wantSeg: TestSegment(""),
		},
		{
			name:    "leading slash index 1",
			s:       "/a/b",
			index:   1,
			wantSeg: TestSegment("a"),
		},
		{
			name:    "leading slash index 2",
			s:       "/a/b",
			index:   2,
			wantSeg: TestSegment("b"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeg := IndexSegments[TestSegment](tt.s, "/", tt.index)
			if gotSeg != tt.wantSeg {
				t.Errorf("IndexSegments() = %v, want %v", gotSeg, tt.wantSeg)
			}
		})
	}
}

func TestSliceSegments(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		start   int
		end     int
		wantOut []TestSegment
	}{
		{
			name:    "empty segments slice 0:1",
			s:       "",
			start:   0,
			end:     1,
			wantOut: []TestSegment{},
		},
		{
			name:    "empty segments slice 0:0",
			s:       "",
			start:   0,
			end:     0,
			wantOut: []TestSegment{},
		},
		{
			name:    "single segment slice 0:1",
			s:       "a",
			start:   0,
			end:     1,
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "single segment slice 0:2",
			s:       "a",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "single segment slice 1:2",
			s:       "a",
			start:   1,
			end:     2,
			wantOut: []TestSegment{},
		},
		{
			name:    "two segments slice 0:1",
			s:       "a/b",
			start:   0,
			end:     1,
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "two segments slice 0:2",
			s:       "a/b",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "two segments slice 1:2",
			s:       "a/b",
			start:   1,
			end:     2,
			wantOut: []TestSegment{"b"},
		},
		{
			name:    "two segments slice 0:3",
			s:       "a/b",
			start:   0,
			end:     3,
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "two segments slice 2:3",
			s:       "a/b",
			start:   2,
			end:     3,
			wantOut: []TestSegment{},
		},
		{
			name:    "three segments slice 0:1",
			s:       "a/b/c",
			start:   0,
			end:     1,
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "three segments slice 0:2",
			s:       "a/b/c",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "three segments slice 0:3",
			s:       "a/b/c",
			start:   0,
			end:     3,
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "three segments slice 1:2",
			s:       "a/b/c",
			start:   1,
			end:     2,
			wantOut: []TestSegment{"b"},
		},
		{
			name:    "three segments slice 1:3",
			s:       "a/b/c",
			start:   1,
			end:     3,
			wantOut: []TestSegment{"b", "c"},
		},
		{
			name:    "three segments slice 2:3",
			s:       "a/b/c",
			start:   2,
			end:     3,
			wantOut: []TestSegment{"c"},
		},
		{
			name:    "six segments slice 0:3",
			s:       "a/b/c/d/e/f",
			start:   0,
			end:     3,
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "six segments slice 2:5",
			s:       "a/b/c/d/e/f",
			start:   2,
			end:     5,
			wantOut: []TestSegment{"c", "d", "e"},
		},
		{
			name:    "six segments slice 3:6",
			s:       "a/b/c/d/e/f",
			start:   3,
			end:     6,
			wantOut: []TestSegment{"d", "e", "f"},
		},
		{
			name:    "six segments slice 4:6",
			s:       "a/b/c/d/e/f",
			start:   4,
			end:     6,
			wantOut: []TestSegment{"e", "f"},
		},
		{
			name:    "out of bounds slice 0:10",
			s:       "a/b/c",
			start:   0,
			end:     10,
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "out of bounds slice 2:10",
			s:       "a/b/c",
			start:   2,
			end:     10,
			wantOut: []TestSegment{"c"},
		},
		{
			name:    "out of bounds slice 10:20",
			s:       "a/b/c",
			start:   10,
			end:     20,
			wantOut: []TestSegment{},
		},
		{
			name:    "start equals end slice 1:1",
			s:       "a/b/c",
			start:   1,
			end:     1,
			wantOut: []TestSegment{},
		},
		{
			name:    "start greater than end slice 2:1",
			s:       "a/b/c",
			start:   2,
			end:     1,
			wantOut: []TestSegment{},
		},
		{
			name:    "start greater than end slice 3:0",
			s:       "a/b/c",
			start:   3,
			end:     0,
			wantOut: []TestSegment{},
		},
		{
			name:    "negative start slice -1:2",
			s:       "a/b/c",
			start:   -1,
			end:     2,
			wantOut: []TestSegment{},
		},
		{
			name:    "negative end slice 0:-1",
			s:       "a/b/c",
			start:   0,
			end:     -1,
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "both negative slice -2:-1",
			s:       "a/b/c",
			start:   -2,
			end:     -1,
			wantOut: []TestSegment{},
		},
		{
			name:    "longer names slice 0:2",
			s:       "alpha/beta/gamma",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"alpha", "beta"},
		},
		{
			name:    "longer names slice 1:3",
			s:       "alpha/beta/gamma",
			start:   1,
			end:     3,
			wantOut: []TestSegment{"beta", "gamma"},
		},
		{
			name:    "hyphens slice 0:2",
			s:       "foo-bar/baz-qux/quux-corge",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"foo-bar", "baz-qux"},
		},
		{
			name:    "underscores slice 1:2",
			s:       "foo_bar/baz_qux/test_file",
			start:   1,
			end:     2,
			wantOut: []TestSegment{"baz_qux"},
		},
		{
			name:    "consecutive slashes slice 0:3",
			s:       "a//b/c",
			start:   0,
			end:     3,
			wantOut: []TestSegment{"a", "", "b"},
		},
		{
			name:    "consecutive slashes slice 1:3",
			s:       "a//b/c",
			start:   1,
			end:     3,
			wantOut: []TestSegment{"", "b"},
		},
		{
			name:    "leading slash slice 0:2",
			s:       "/a/b/c",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"", "a"},
		},
		{
			name:    "leading slash slice 1:3",
			s:       "/a/b/c",
			start:   1,
			end:     3,
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "api path slice 0:2",
			s:       "api/v1/users/123",
			start:   0,
			end:     2,
			wantOut: []TestSegment{"api", "v1"},
		},
		{
			name:    "api path slice 2:4",
			s:       "api/v1/users/123",
			start:   2,
			end:     4,
			wantOut: []TestSegment{"users", "123"},
		},
		{
			name:    "single segment slice 0:-1",
			s:       "a",
			start:   0,
			end:     -1,
			wantOut: []TestSegment{"a"},
		},
		{
			name:    "two segments slice 0:-1",
			s:       "a/b",
			start:   0,
			end:     -1,
			wantOut: []TestSegment{"a", "b"},
		},
		{
			name:    "three segments slice 0:-1",
			s:       "a/b/c",
			start:   0,
			end:     -1,
			wantOut: []TestSegment{"a", "b", "c"},
		},
		{
			name:    "three segments slice 1:-1",
			s:       "a/b/c",
			start:   1,
			end:     -1,
			wantOut: []TestSegment{"b", "c"},
		},
		{
			name:    "three segments slice 2:-1",
			s:       "a/b/c",
			start:   2,
			end:     -1,
			wantOut: []TestSegment{"c"},
		},
		{
			name:    "six segments slice 2:-1",
			s:       "a/b/c/d/e/f",
			start:   2,
			end:     -1,
			wantOut: []TestSegment{"c", "d", "e", "f"},
		},
		{
			name:    "api path slice 2:-1",
			s:       "api/v1/users/123",
			start:   2,
			end:     -1,
			wantOut: []TestSegment{"users", "123"},
		},
		{
			name:    "api path slice 3:-1",
			s:       "api/v1/users/123",
			start:   3,
			end:     -1,
			wantOut: []TestSegment{"123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := SliceSegments[TestSegment](tt.s, "/", tt.start, tt.end)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("SliceSegments() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestJoinSegments(t *testing.T) {
	tests := []struct {
		name    string
		ss      []TestSegment
		sep     string
		wantOut TestSegment
	}{
		{
			name:    "empty slice with slash",
			ss:      []TestSegment{},
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "single segment with slash",
			ss:      []TestSegment{"a"},
			sep:     "/",
			wantOut: TestSegment("a"),
		},
		{
			name:    "two segments with slash",
			ss:      []TestSegment{"a", "b"},
			sep:     "/",
			wantOut: TestSegment("a/b"),
		},
		{
			name:    "three segments with slash",
			ss:      []TestSegment{"a", "b", "c"},
			sep:     "/",
			wantOut: TestSegment("a/b/c"),
		},
		{
			name:    "longer segment names with slash",
			ss:      []TestSegment{"alpha", "beta", "gamma"},
			sep:     "/",
			wantOut: TestSegment("alpha/beta/gamma"),
		},
		{
			name:    "api path with slash",
			ss:      []TestSegment{"api", "v1", "users", "123"},
			sep:     "/",
			wantOut: TestSegment("api/v1/users/123"),
		},
		{
			name:    "segments with hyphens and slash",
			ss:      []TestSegment{"foo-bar", "baz-qux"},
			sep:     "/",
			wantOut: TestSegment("foo-bar/baz-qux"),
		},
		{
			name:    "segments with underscores and slash",
			ss:      []TestSegment{"foo_bar", "baz_qux"},
			sep:     "/",
			wantOut: TestSegment("foo_bar/baz_qux"),
		},
		{
			name:    "two segments with backslash",
			ss:      []TestSegment{"a", "b"},
			sep:     "\\",
			wantOut: TestSegment("a\\b"),
		},
		{
			name:    "three segments with backslash",
			ss:      []TestSegment{"a", "b", "c"},
			sep:     "\\",
			wantOut: TestSegment("a\\b\\c"),
		},
		{
			name:    "windows path with backslash",
			ss:      []TestSegment{"Users", "john", "Documents"},
			sep:     "\\",
			wantOut: TestSegment("Users\\john\\Documents"),
		},
		{
			name:    "deep path with slash",
			ss:      []TestSegment{"a", "b", "c", "d", "e", "f"},
			sep:     "/",
			wantOut: TestSegment("a/b/c/d/e/f"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := JoinSegments[TestSegment](tt.ss, tt.sep)
			if gotOut != tt.wantOut {
				t.Errorf("JoinSegments() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestSliceSegmentsScalar(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		start   int
		end     int
		sep     string
		wantOut TestSegment
	}{
		{
			name:    "empty string slice 0:1 with slash",
			s:       "",
			start:   0,
			end:     1,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "single segment slice 0:1 with slash",
			s:       "a",
			start:   0,
			end:     1,
			sep:     "/",
			wantOut: TestSegment("a"),
		},
		{
			name:    "two segments slice 0:2 with slash",
			s:       "a/b",
			start:   0,
			end:     2,
			sep:     "/",
			wantOut: TestSegment("a/b"),
		},
		{
			name:    "three segments slice 0:3 with slash",
			s:       "a/b/c",
			start:   0,
			end:     3,
			sep:     "/",
			wantOut: TestSegment("a/b/c"),
		},
		{
			name:    "three segments slice 1:2 with slash",
			s:       "a/b/c",
			start:   1,
			end:     2,
			sep:     "/",
			wantOut: TestSegment("b"),
		},
		{
			name:    "three segments slice 1:3 with slash",
			s:       "a/b/c",
			start:   1,
			end:     3,
			sep:     "/",
			wantOut: TestSegment("b/c"),
		},
		{
			name:    "longer names slice 0:2 with slash",
			s:       "alpha/beta/gamma",
			start:   0,
			end:     2,
			sep:     "/",
			wantOut: TestSegment("alpha/beta"),
		},
		{
			name:    "api path slice 2:4 with slash",
			s:       "api/v1/users/123",
			start:   2,
			end:     4,
			sep:     "/",
			wantOut: TestSegment("users/123"),
		},
		{
			name:    "two segments slice 0:2 with backslash",
			s:       "a\\b",
			start:   0,
			end:     2,
			sep:     "\\",
			wantOut: TestSegment("a\\b"),
		},
		{
			name:    "windows path slice 1:3 with backslash",
			s:       "Users\\john\\Documents",
			start:   1,
			end:     3,
			sep:     "\\",
			wantOut: TestSegment("john\\Documents"),
		},
		{
			name:    "single segment slice 0:-1 with slash",
			s:       "a",
			start:   0,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("a"),
		},
		{
			name:    "two segments slice 0:-1 with slash",
			s:       "a/b",
			start:   0,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("a/b"),
		},
		{
			name:    "three segments slice 0:-1 with slash",
			s:       "a/b/c",
			start:   0,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("a/b/c"),
		},
		{
			name:    "three segments slice 1:-1 with slash",
			s:       "a/b/c",
			start:   1,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("b/c"),
		},
		{
			name:    "three segments slice 2:-1 with slash",
			s:       "a/b/c",
			start:   2,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("c"),
		},
		{
			name:    "api path slice 2:-1 with slash",
			s:       "api/v1/users/123",
			start:   2,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("users/123"),
		},
		{
			name:    "api path slice 3:-1 with slash",
			s:       "api/v1/users/123",
			start:   3,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("123"),
		},
		{
			name:    "windows path slice 1:-1 with backslash",
			s:       "Users\\john\\Documents",
			start:   1,
			end:     -1,
			sep:     "\\",
			wantOut: TestSegment("john\\Documents"),
		},
		// Edge cases for beyond-segment-count start position (regression test for makeslice panic)
		{
			name:    "two segments slice 10:20 out of bounds with slash",
			s:       "a/b",
			start:   10,
			end:     20,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "three segments slice 10:-1 start beyond with slash",
			s:       "a/b/c",
			start:   10,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "api path slice 10:15 start beyond segments with slash",
			s:       "api/v1/users/123",
			start:   10,
			end:     15,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		// GitHub path test case - extracting repo name from "org/repo"
		// This is the exact use case from demo_source.go:113
		// slug := ds.Repo.SliceScalar(ds.Repo.LastIndex("/")+1, -1, "/")
		{
			name:    "github repo name extraction from org/repo",
			s:       "xmlui-org/hello",
			start:   1,
			end:     -1,
			sep:     "/",
			wantOut: TestSegment("hello"),
		},
		{
			name:    "github repo last index extraction (regression for makeslice panic)",
			s:       "xmlui-org/demo-app",
			start:   10, // LastIndex("/") + 1 when / is at position 9
			end:     -1,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		// Test with different separators
		{
			name:    "dot separator slice 0:2",
			s:       "com.example.domain",
			start:   0,
			end:     2,
			sep:     ".",
			wantOut: TestSegment("com.example"),
		},
		{
			name:    "dot separator slice 1:3",
			s:       "com.example.domain",
			start:   1,
			end:     3,
			sep:     ".",
			wantOut: TestSegment("example.domain"),
		},
		// Empty segment handling
		{
			name:    "consecutive separators slice 0:3",
			s:       "a///b",
			start:   0,
			end:     3,
			sep:     "/",
			wantOut: TestSegment("a//"),
		},
		{
			name:    "consecutive separators slice 1:4",
			s:       "a///b",
			start:   1,
			end:     4,
			sep:     "/",
			wantOut: TestSegment("//b"),
		},
		// Negative/invalid indices
		{
			name:    "negative start returns empty",
			s:       "a/b/c",
			start:   -1,
			end:     2,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "invalid range start > end",
			s:       "a/b/c",
			start:   3,
			end:     1,
			sep:     "/",
			wantOut: TestSegment(""),
		},
		{
			name:    "invalid range start == end",
			s:       "a/b/c",
			start:   1,
			end:     1,
			sep:     "/",
			wantOut: TestSegment(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := SliceSegmentsScalar[TestSegment](tt.s, tt.sep, tt.start, tt.end)
			if gotOut != tt.wantOut {
				t.Errorf("SliceSegmentsScalar() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

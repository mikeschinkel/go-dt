package dt

import (
	"reflect"
	"testing"
)

func TestURLSegments_Segment(t *testing.T) {
	tests := []struct {
		name    string
		pss     URLSegments
		index   int
		wantSeg URLSegment
	}{
		// Empty URLSegments
		{
			name:    "empty segments with index 0",
			pss:     "",
			index:   0,
			wantSeg: "",
		},
		{
			name:    "empty segments with index 1",
			pss:     "",
			index:   1,
			wantSeg: "",
		},
		{
			name:    "empty segments with negative index",
			pss:     "",
			index:   -1,
			wantSeg: "",
		},
		// Single segment (no slashes)
		{
			name:    "single segment with index 0",
			pss:     "a",
			index:   0,
			wantSeg: "a",
		},
		{
			name:    "single segment with index 1",
			pss:     "a",
			index:   1,
			wantSeg: "",
		},
		// Two segments
		{
			name:    "two segments with index 0",
			pss:     "a/b",
			index:   0,
			wantSeg: "a",
		},
		{
			name:    "two segments with index 1",
			pss:     "a/b",
			index:   1,
			wantSeg: "b",
		},
		{
			name:    "two segments with index 2",
			pss:     "a/b",
			index:   2,
			wantSeg: "",
		},
		{
			name:    "two segments with negative index",
			pss:     "a/b",
			index:   -1,
			wantSeg: "",
		},
		// Three segments
		{
			name:    "three segments with index 0",
			pss:     "a/b/c",
			index:   0,
			wantSeg: "a",
		},
		{
			name:    "three segments with index 1",
			pss:     "a/b/c",
			index:   1,
			wantSeg: "b",
		},
		{
			name:    "three segments with index 2",
			pss:     "a/b/c",
			index:   2,
			wantSeg: "c",
		},
		{
			name:    "three segments with index 3",
			pss:     "a/b/c",
			index:   3,
			wantSeg: "",
		},
		// Multiple segments with longer names
		{
			name:    "longer segment names index 0",
			pss:     "alpha/beta/gamma",
			index:   0,
			wantSeg: "alpha",
		},
		{
			name:    "longer segment names index 1",
			pss:     "alpha/beta/gamma",
			index:   1,
			wantSeg: "beta",
		},
		{
			name:    "longer segment names index 2",
			pss:     "alpha/beta/gamma",
			index:   2,
			wantSeg: "gamma",
		},
		// Deep path
		{
			name:    "deep path index 0",
			pss:     "a/b/c/d/e/f",
			index:   0,
			wantSeg: "a",
		},
		{
			name:    "deep path index 3",
			pss:     "a/b/c/d/e/f",
			index:   3,
			wantSeg: "d",
		},
		{
			name:    "deep path index 5",
			pss:     "a/b/c/d/e/f",
			index:   5,
			wantSeg: "f",
		},
		// Segments with special characters
		{
			name:    "segments with hyphens index 0",
			pss:     "foo-bar/baz-qux",
			index:   0,
			wantSeg: "foo-bar",
		},
		{
			name:    "segments with hyphens index 1",
			pss:     "foo-bar/baz-qux",
			index:   1,
			wantSeg: "baz-qux",
		},
		{
			name:    "segments with underscores index 0",
			pss:     "foo_bar/baz_qux",
			index:   0,
			wantSeg: "foo_bar",
		},
		{
			name:    "segments with underscores index 1",
			pss:     "foo_bar/baz_qux",
			index:   1,
			wantSeg: "baz_qux",
		},
		{
			name:    "segments with numbers index 0",
			pss:     "path1/path2",
			index:   0,
			wantSeg: "path1",
		},
		{
			name:    "segments with numbers index 1",
			pss:     "path1/path2",
			index:   1,
			wantSeg: "path2",
		},
		// Edge cases with consecutive slashes
		{
			name:    "consecutive slashes index 0",
			pss:     "a//b",
			index:   0,
			wantSeg: "a",
		},
		{
			name:    "consecutive slashes index 1",
			pss:     "a//b",
			index:   1,
			wantSeg: "",
		},
		{
			name:    "consecutive slashes index 2",
			pss:     "a//b",
			index:   2,
			wantSeg: "b",
		},
		// Leading slash
		{
			name:    "leading slash index 0",
			pss:     "/a/b",
			index:   0,
			wantSeg: "",
		},
		{
			name:    "leading slash index 1",
			pss:     "/a/b",
			index:   1,
			wantSeg: "a",
		},
		{
			name:    "leading slash index 2",
			pss:     "/a/b",
			index:   2,
			wantSeg: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSeg := tt.pss.Segment(tt.index); gotSeg != tt.wantSeg {
				t.Errorf("Segment() = %v, want %v", gotSeg, tt.wantSeg)
			}
		})
	}
}

func TestURLSegments_Split(t *testing.T) {
	tests := []struct {
		name    string
		pss     URLSegments
		wantOut []URLSegment
	}{
		// Empty segments
		{
			name:    "empty string",
			pss:     "",
			wantOut: []URLSegment{},
		},
		// Single segment (no slashes)
		{
			name:    "single segment",
			pss:     "a",
			wantOut: []URLSegment{"a"},
		},
		{
			name:    "single longer segment",
			pss:     "alpha",
			wantOut: []URLSegment{"alpha"},
		},
		// Two segments
		{
			name:    "two segments",
			pss:     "a/b",
			wantOut: []URLSegment{"a", "b"},
		},
		{
			name:    "two longer segments",
			pss:     "alpha/beta",
			wantOut: []URLSegment{"alpha", "beta"},
		},
		// Three segments
		{
			name:    "three segments",
			pss:     "a/b/c",
			wantOut: []URLSegment{"a", "b", "c"},
		},
		{
			name:    "three longer segments",
			pss:     "alpha/beta/gamma",
			wantOut: []URLSegment{"alpha", "beta", "gamma"},
		},
		// Four segments
		{
			name:    "four segments",
			pss:     "a/b/c/d",
			wantOut: []URLSegment{"a", "b", "c", "d"},
		},
		// Deep path with many segments
		{
			name:    "six segments",
			pss:     "a/b/c/d/e/f",
			wantOut: []URLSegment{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:    "deep path ten segments",
			pss:     "a/b/c/d/e/f/g/h/i/j",
			wantOut: []URLSegment{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		// Segments with special characters
		{
			name:    "segments with hyphens",
			pss:     "foo-bar/baz-qux/quux-corge",
			wantOut: []URLSegment{"foo-bar", "baz-qux", "quux-corge"},
		},
		{
			name:    "segments with underscores",
			pss:     "foo_bar/baz_qux",
			wantOut: []URLSegment{"foo_bar", "baz_qux"},
		},
		{
			name:    "segments with numbers",
			pss:     "path1/path2/path3",
			wantOut: []URLSegment{"path1", "path2", "path3"},
		},
		{
			name:    "segments with dots",
			pss:     "file.txt/dir.name/file.ext",
			wantOut: []URLSegment{"file.txt", "dir.name", "file.ext"},
		},
		// Edge cases with consecutive slashes
		{
			name:    "consecutive slashes two segments",
			pss:     "a//b",
			wantOut: []URLSegment{"a", "", "b"},
		},
		{
			name:    "consecutive slashes three segments",
			pss:     "a///b",
			wantOut: []URLSegment{"a", "", "", "b"},
		},
		{
			name: "leading slash",
			pss:  "/a/b",
			//wantOut: []URLSegment{""), URLSegment("a", "b"},
			wantOut: []URLSegment{"", "a", "b"},
		},
		// URL-like paths
		{
			name:    "domain path",
			pss:     "api/v1/users",
			wantOut: []URLSegment{"api", "v1", "users"},
		},
		{
			name:    "api endpoint",
			pss:     "api/v2/users/123",
			wantOut: []URLSegment{"api", "v2", "users", "123"},
		},
		// Mixed case and character segments
		{
			name:    "mixed case",
			pss:     "MyPath/MyFile/MyExt",
			wantOut: []URLSegment{"MyPath", "MyFile", "MyExt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := tt.pss.Split()
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("Split() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestURLSegments_Slice(t *testing.T) {
	tests := []struct {
		name    string
		pss     URLSegments
		start   int
		end     int
		wantOut []URLSegment
	}{
		// Empty segments
		{
			name:    "empty segments slice 0:1",
			pss:     "",
			start:   0,
			end:     1,
			wantOut: []URLSegment{},
		},
		{
			name:    "empty segments slice 0:0",
			pss:     "",
			start:   0,
			end:     0,
			wantOut: []URLSegment{},
		},
		// Single segment
		{
			name:    "single segment slice 0:1",
			pss:     "a",
			start:   0,
			end:     1,
			wantOut: []URLSegment{"a"},
		},
		{
			name:    "single segment slice 0:2",
			pss:     "a",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"a"},
		},
		{
			name:    "single segment slice 1:2",
			pss:     "a",
			start:   1,
			end:     2,
			wantOut: []URLSegment{},
		},
		// Two segments
		{
			name:    "two segments slice 0:1",
			pss:     "a/b",
			start:   0,
			end:     1,
			wantOut: []URLSegment{"a"},
		},
		{
			name:    "two segments slice 0:2",
			pss:     "a/b",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"a", "b"},
		},
		{
			name:    "two segments slice 1:2",
			pss:     "a/b",
			start:   1,
			end:     2,
			wantOut: []URLSegment{"b"},
		},
		{
			name:    "two segments slice 0:3",
			pss:     "a/b",
			start:   0,
			end:     3,
			wantOut: []URLSegment{"a", "b"},
		},
		{
			name:    "two segments slice 2:3",
			pss:     "a/b",
			start:   2,
			end:     3,
			wantOut: []URLSegment{},
		},
		// Three segments
		{
			name:    "three segments slice 0:1",
			pss:     "a/b/c",
			start:   0,
			end:     1,
			wantOut: []URLSegment{"a"},
		},
		{
			name:    "three segments slice 0:2",
			pss:     "a/b/c",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"a", "b"},
		},
		{
			name:    "three segments slice 0:3",
			pss:     "a/b/c",
			start:   0,
			end:     3,
			wantOut: []URLSegment{"a", "b", "c"},
		},
		{
			name:    "three segments slice 1:2",
			pss:     "a/b/c",
			start:   1,
			end:     2,
			wantOut: []URLSegment{"b"},
		},
		{
			name:    "three segments slice 1:3",
			pss:     "a/b/c",
			start:   1,
			end:     3,
			wantOut: []URLSegment{"b", "c"},
		},
		{
			name:    "three segments slice 2:3",
			pss:     "a/b/c",
			start:   2,
			end:     3,
			wantOut: []URLSegment{"c"},
		},
		// Deep path (six segments)
		{
			name:    "six segments slice 0:3",
			pss:     "a/b/c/d/e/f",
			start:   0,
			end:     3,
			wantOut: []URLSegment{"a", "b", "c"},
		},
		{
			name:    "six segments slice 2:5",
			pss:     "a/b/c/d/e/f",
			start:   2,
			end:     5,
			wantOut: []URLSegment{"c", "d", "e"},
		},
		{
			name:    "six segments slice 3:6",
			pss:     "a/b/c/d/e/f",
			start:   3,
			end:     6,
			wantOut: []URLSegment{"d", "e", "f"},
		},
		{
			name:    "six segments slice 4:6",
			pss:     "a/b/c/d/e/f",
			start:   4,
			end:     6,
			wantOut: []URLSegment{"e", "f"},
		},
		// Out of bounds
		{
			name:    "out of bounds slice 0:10",
			pss:     "a/b/c",
			start:   0,
			end:     10,
			wantOut: []URLSegment{"a", "b", "c"},
		},
		{
			name:    "out of bounds slice 2:10",
			pss:     "a/b/c",
			start:   2,
			end:     10,
			wantOut: []URLSegment{"c"},
		},
		{
			name:    "out of bounds slice 10:20",
			pss:     "a/b/c",
			start:   10,
			end:     20,
			wantOut: []URLSegment{},
		},
		// Invalid ranges
		{
			name:    "start equals end slice 1:1",
			pss:     "a/b/c",
			start:   1,
			end:     1,
			wantOut: []URLSegment{},
		},
		{
			name:    "start greater than end slice 2:1",
			pss:     "a/b/c",
			start:   2,
			end:     1,
			wantOut: []URLSegment{},
		},
		{
			name:    "start greater than end slice 3:0",
			pss:     "a/b/c",
			start:   3,
			end:     0,
			wantOut: []URLSegment{},
		},
		// Negative indices
		{
			name:    "negative start slice -1:2",
			pss:     "a/b/c",
			start:   -1,
			end:     2,
			wantOut: []URLSegment{},
		},
		{
			name:    "negative end slice 0:-1",
			pss:     "a/b/c",
			start:   0,
			end:     -1,
			wantOut: []URLSegment{"a", "b", "c"},
		},
		{
			name:    "both negative slice -2:-1",
			pss:     "a/b/c",
			start:   -2,
			end:     -1,
			wantOut: []URLSegment{},
		},
		// Longer segment names
		{
			name:    "longer names slice 0:2",
			pss:     "alpha/beta/gamma",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"alpha", "beta"},
		},
		{
			name:    "longer names slice 1:3",
			pss:     "alpha/beta/gamma",
			start:   1,
			end:     3,
			wantOut: []URLSegment{"beta", "gamma"},
		},
		// Segments with special characters
		{
			name:    "hyphens slice 0:2",
			pss:     "foo-bar/baz-qux/quux-corge",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"foo-bar", "baz-qux"},
		},
		{
			name:    "underscores slice 1:2",
			pss:     "foo_bar/baz_qux/test_file",
			start:   1,
			end:     2,
			wantOut: []URLSegment{"baz_qux"},
		},
		// Edge case with consecutive slashes
		{
			name:    "consecutive slashes slice 0:3",
			pss:     "a//b/c",
			start:   0,
			end:     3,
			wantOut: []URLSegment{"a", "", "b"},
		},
		{
			name:    "consecutive slashes slice 1:3",
			pss:     "a//b/c",
			start:   1,
			end:     3,
			wantOut: []URLSegment{"", "b"},
		},
		// Leading slash
		{
			name:    "leading slash slice 0:2",
			pss:     "/a/b/c",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"", "a"},
		},
		{
			name:    "leading slash slice 1:3",
			pss:     "/a/b/c",
			start:   1,
			end:     3,
			wantOut: []URLSegment{"a", "b"},
		},
		// URL-like paths
		{
			name:    "api path slice 0:2",
			pss:     "api/v1/users/123",
			start:   0,
			end:     2,
			wantOut: []URLSegment{"api", "v1"},
		},
		{
			name:    "api path slice 2:4",
			pss:     "api/v1/users/123",
			start:   2,
			end:     4,
			wantOut: []URLSegment{"users", "123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := tt.pss.Slice(tt.start, tt.end)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("Slice(%d, %d) = %v, want %v", tt.start, tt.end, gotOut, tt.wantOut)
			}
		})
	}
}

func TestURLSegments_SliceScalar(t *testing.T) {
	tests := []struct {
		name    string
		uss     URLSegments
		start   int
		end     int
		sep     string
		wantOut string
	}{
		{
			name:    "empty segments slice 0:1",
			uss:     "",
			start:   0,
			end:     1,
			sep:     "/",
			wantOut: "",
		},
		{
			name:    "single segment slice 0:1",
			uss:     "a",
			start:   0,
			end:     1,
			sep:     "/",
			wantOut: "a",
		},
		{
			name:    "two segments slice 0:2",
			uss:     "a/b",
			start:   0,
			end:     2,
			sep:     "/",
			wantOut: "a/b",
		},
		{
			name:    "three segments slice 1:2",
			uss:     "a/b/c",
			start:   1,
			end:     2,
			sep:     "/",
			wantOut: "b",
		},
		{
			name:    "three segments slice 1:3",
			uss:     "a/b/c",
			start:   1,
			end:     3,
			sep:     "/",
			wantOut: "b/c",
		},
		{
			name:    "api path slice 0:2",
			uss:     "api/v1/users/123",
			start:   0,
			end:     2,
			sep:     "/",
			wantOut: "api/v1",
		},
		{
			name:    "api path slice 2:4",
			uss:     "api/v1/users/123",
			start:   2,
			end:     4,
			sep:     "/",
			wantOut: "users/123",
		},
		{
			name:    "single segment slice 0:-1",
			uss:     "a",
			start:   0,
			end:     -1,
			sep:     "/",
			wantOut: "a",
		},
		{
			name:    "two segments slice 0:-1",
			uss:     "a/b",
			start:   0,
			end:     -1,
			sep:     "/",
			wantOut: "a/b",
		},
		{
			name:    "three segments slice 1:-1",
			uss:     "a/b/c",
			start:   1,
			end:     -1,
			sep:     "/",
			wantOut: "b/c",
		},
		{
			name:    "three segments slice 2:-1",
			uss:     "a/b/c",
			start:   2,
			end:     -1,
			sep:     "/",
			wantOut: "c",
		},
		{
			name:    "api path slice 2:-1",
			uss:     "api/v1/users/123",
			start:   2,
			end:     -1,
			sep:     "/",
			wantOut: "users/123",
		},
		{
			name:    "api path slice 3:-1",
			uss:     "api/v1/users/123",
			start:   3,
			end:     -1,
			sep:     "/",
			wantOut: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := string(tt.uss.SliceScalar(tt.start, tt.end, tt.sep))
			if gotOut != tt.wantOut {
				t.Errorf("SliceScalar(%d, %d, %q) = %q, want %q", tt.start, tt.end, tt.sep, gotOut, tt.wantOut)
			}
		})
	}
}

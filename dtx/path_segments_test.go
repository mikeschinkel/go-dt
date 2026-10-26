package dtx

import (
	"testing"
)

// TestParsePathSegmentForOS covers each OS's branch on whatever OS runs the
// test; ParseOSPathSegment alone only ever exercises the host's.
func TestParsePathSegmentForOS(t *testing.T) {
	tests := []struct {
		goos    string
		input   string
		wantErr bool
	}{
		{goos: "windows", input: "file.txt"},
		{goos: "windows", input: "a:b", wantErr: true},
		{goos: "darwin", input: "file.txt"},
		{goos: "darwin", input: "a/b", wantErr: true},
		{goos: "linux", input: "file.txt"},
		{goos: "linux", input: "a/b", wantErr: true},
		{goos: "freebsd", input: "file.txt"},
		{goos: "freebsd", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.goos+"/"+tt.input, func(t *testing.T) {
			ps, err := parsePathSegmentForOS(tt.goos, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parsePathSegmentForOS(%q, %q) = %q, want an error", tt.goos, tt.input, ps)
				}
				return
			}
			if err != nil {
				t.Fatalf("parsePathSegmentForOS(%q, %q) error: %v", tt.goos, tt.input, err)
			}
			if string(ps) != tt.input {
				t.Errorf("parsePathSegmentForOS(%q, %q) = %q, want %q", tt.goos, tt.input, ps, tt.input)
			}
		})
	}
}

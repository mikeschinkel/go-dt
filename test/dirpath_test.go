package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParseDirPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    dt.DirPath
		wantErr bool
	}{
		{
			name:    "valid absolute path",
			input:   "/usr/local/bin",
			want:    dt.DirPath("/usr/local/bin"),
			wantErr: false,
		},
		{
			name:    "valid relative path",
			input:   "relative/path",
			want:    dt.DirPath("relative/path"),
			wantErr: false,
		},
		{
			name:    "root path",
			input:   "/",
			want:    dt.DirPath("/"),
			wantErr: false,
		},
		{
			name:    "empty path",
			input:   "",
			want:    dt.DirPath(""),
			wantErr: true,
		},
		{
			name:    "literal tilde name",
			input:   "~noslash",
			want:    dt.DirPath("~noslash"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseDirPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDirPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []dt.DirPath
		wantErr bool
	}{
		{
			name:  "multiple valid paths",
			input: []string{"/usr/local", "/var/log", "/tmp"},
			want: []dt.DirPath{
				dt.DirPath("/usr/local"),
				dt.DirPath("/var/log"),
				dt.DirPath("/tmp"),
			},
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []string{},
			want:    []dt.DirPath{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseDirPaths(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseDirPaths() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ParseDirPaths()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

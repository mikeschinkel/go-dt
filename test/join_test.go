package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestDirPathJoin(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want dt.DirPath
	}{
		{
			name: "join two paths",
			a:    "/usr/local",
			b:    "bin",
			want: dt.DirPath("/usr/local/bin"),
		},
		{
			name: "join with trailing slash",
			a:    "/etc/",
			b:    "config",
			want: dt.DirPath("/etc/config"),
		},
		{
			name: "join relative paths",
			a:    "home",
			b:    "user",
			want: dt.DirPath("home/user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dt.DirPathJoin(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("DirPathJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilepathJoin(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want dt.Filepath
	}{
		{
			name: "join dirpath and filename",
			a:    "/etc",
			b:    "config.json",
			want: dt.Filepath("/etc/config.json"),
		},
		{
			name: "join relative paths",
			a:    "home/user",
			b:    "document.txt",
			want: dt.Filepath("home/user/document.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dt.FilepathJoin(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("FilepathJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirPathMethodJoin(t *testing.T) {
	tests := []struct {
		name  string
		dp    dt.DirPath
		elems []any
		want  dt.DirPath
	}{
		{
			name:  "join with string",
			dp:    dt.DirPath("/usr/local"),
			elems: []any{"bin"},
			want:  dt.DirPath("/usr/local/bin"),
		},
		{
			name:  "join with multiple elements",
			dp:    dt.DirPath("/home"),
			elems: []any{"user", "documents"},
			want:  dt.DirPath("/home/user/documents"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dp.Join(tt.elems...)
			if got != tt.want {
				t.Errorf("DirPath.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

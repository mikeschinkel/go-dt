package dt

import (
	"path/filepath"
	"testing"
)

func TestToTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := EnsureUserHomeDir(); err != nil {
		t.Fatalf("EnsureUserHomeDir() error = %v", err)
	}

	sep := string(filepath.Separator)
	nested := filepath.Join(home, "sub", "dir")
	outside := filepath.Join(filepath.Dir(home), "outside")
	relative := filepath.Join("rel", "path")

	tests := []struct {
		name      string
		path      EntryPath
		option    TildeOption
		want      string
		wantPanic bool
	}{
		{name: "empty", path: "", option: OrFullPath, want: ""},
		{name: "home", path: EntryPath(home), option: OrPanic, want: "~"},
		{name: "nested", path: EntryPath(nested), option: OrPanic, want: "~" + sep + filepath.Join("sub", "dir")},
		{name: "already tilde", path: EntryPath("~" + sep + "sub"), option: OrFullPath, want: "~" + sep + "sub"},
		{name: "tilde only", path: EntryPath("~"), option: OrFullPath, want: "~"},
		{name: "outside empty", path: EntryPath(outside), option: OrEmptyString, want: ""},
		{name: "outside full", path: EntryPath(outside), option: OrFullPath, want: outside},
		{name: "outside panic", path: EntryPath(outside), option: OrPanic, wantPanic: true},
		{name: "relative empty", path: EntryPath(relative), option: OrEmptyString, want: ""},
		{name: "relative full", path: EntryPath(relative), option: OrFullPath, want: relative},
		{name: "relative panic", path: EntryPath(relative), option: OrPanic, wantPanic: true},
		{name: "invalid option", path: EntryPath(home), option: UnspecifiedTildeOption, wantPanic: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("ToTilde() expected panic")
					}
				}()
			}
			got := ToTilde[EntryPath, TildeEntryPath](tt.path, tt.option)
			if !tt.wantPanic && string(got) != tt.want {
				t.Fatalf("ToTilde() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimTilde(t *testing.T) {
	sep := string(filepath.Separator)
	tests := []struct {
		name string
		path EntryPath
		want PathSegments
	}{
		{name: "empty", path: "", want: ""},
		{name: "tilde only", path: "~", want: "~"},
		{name: "tilde sep", path: EntryPath("~" + sep), want: ""},
		{name: "tilde nested", path: EntryPath("~" + sep + filepath.Join("sub", "dir")), want: PathSegments(filepath.Join("sub", "dir"))},
		{name: "non-tilde", path: EntryPath(filepath.Join("sub", "dir")), want: PathSegments(filepath.Join("sub", "dir"))},
		{name: "tilde non-sep", path: "~x", want: "~x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimTilde(tt.path)
			if got != tt.want {
				t.Fatalf("TrimTilde() = %q, want %q", got, tt.want)
			}
		})
	}
}

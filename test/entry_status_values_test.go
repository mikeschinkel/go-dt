package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// TestEntryStatusValues pins each EntryStatus constant to its number, so a
// reordered or mistyped constant fails here rather than silently changing
// what a stored value means.
func TestEntryStatusValues(t *testing.T) {
	tests := []struct {
		name   string
		status dt.EntryStatus
		want   uint8
	}{
		{name: "IsInvalidEntryStatus", status: dt.IsInvalidEntryStatus, want: 0},
		{name: "IsEntryError", status: dt.IsEntryError, want: 1},
		{name: "IsMissingEntry", status: dt.IsMissingEntry, want: 2},
		{name: "IsFileEntry", status: dt.IsFileEntry, want: 3},
		{name: "IsDirEntry", status: dt.IsDirEntry, want: 4},
		{name: "IsSymlinkEntry", status: dt.IsSymlinkEntry, want: 5},
		{name: "IsSocketEntry", status: dt.IsSocketEntry, want: 6},
		{name: "IsPipeEntry", status: dt.IsPipeEntry, want: 7},
		{name: "IsDeviceEntry", status: dt.IsDeviceEntry, want: 8},
		{name: "IsUnclassifiedEntryStatus", status: dt.IsUnclassifiedEntryStatus, want: 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uint8(tt.status); got != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

// TestEntryStatusFlagsValues pins each EntryStatusFlags constant to its bit.
func TestEntryStatusFlagsValues(t *testing.T) {
	if got := uint32(dt.DontFollowSymlinks); got != 1 {
		t.Errorf("DontFollowSymlinks = %d, want 1", got)
	}
}

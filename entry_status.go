package dt

// EntryStatus represents the classification of a filesystem entry.
//
// It describes whether a given path refers to a file, directory, missing entry,
// or another type of object such as a symlink, socket, or device.
//
// The zero value (IsInvalidEntryStatus) indicates an uninitialized or invalid
// status. IsEntryError denotes that an underlying filesystem error occurred;
// in that case, the accompanying error value from Status or LstatStatus
// provides additional detail.
type EntryStatus uint8

// String returns a human-readable representation of the EntryStatus.
//
// It is primarily intended for diagnostic output and logging; program logic
// should use the EntryStatus constants directly.
func (s EntryStatus) String() string {
	switch s {
	case IsFileEntry:
		return "File entry"
	case IsDirEntry:
		return "Directory entry"
	case IsMissingEntry:
		return "Missing entry"
	case IsEntryError:
		return "Entry error"
	case IsInvalidEntryStatus:
		return "Invalid entry"
	case IsSymlinkEntry:
		return "Symlink entry"
	case IsSocketEntry:
		return "Socket entry"
	case IsPipeEntry:
		return "Pipe entry"
	case IsDeviceEntry:
		return "Device entry"
	case IsUnclassifiedEntryStatus:
		//fallthrough
	}
	return "Unclassified entry"
}

const (
	// IsInvalidEntryStatus is the zero value and indicates an uninitialized
	// or otherwise invalid EntryStatus.
	IsInvalidEntryStatus EntryStatus = iota

	// IsEntryError indicates a filesystem error occurred; the returned error
	// from Status or LstatStatus will contain the details.
	IsEntryError

	// IsMissingEntry indicates the entry does not exist (fs.ErrNotExist).
	IsMissingEntry

	// IsFileEntry indicates a regular file.
	IsFileEntry

	// IsDirEntry indicates a directory.
	IsDirEntry

	IsSymlinkEntry
	IsSocketEntry
	IsPipeEntry
	IsDeviceEntry

	// IsUnclassifiedEntryStatus indicates some other kind of filesystem entry such as
	// a symlink, socket, device, or named pipe.
	IsUnclassifiedEntryStatus
)

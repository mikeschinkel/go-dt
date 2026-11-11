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
	IsInvalidEntryStatus EntryStatus = 0

	// IsEntryError indicates a filesystem error occurred; the returned error
	// from Status or LstatStatus will contain the details.
	IsEntryError EntryStatus = 1

	// IsMissingEntry indicates the entry does not exist (fs.ErrNotExist).
	IsMissingEntry EntryStatus = 2

	// IsFileEntry indicates a regular file.
	IsFileEntry EntryStatus = 3

	// IsDirEntry indicates a directory.
	IsDirEntry EntryStatus = 4

	// IsSymlinkEntry indicates a symbolic link.
	IsSymlinkEntry EntryStatus = 5

	// IsSocketEntry indicates a Unix domain socket.
	IsSocketEntry EntryStatus = 6

	// IsPipeEntry indicates a named pipe (FIFO).
	IsPipeEntry EntryStatus = 7

	// IsDeviceEntry indicates a device file.
	IsDeviceEntry EntryStatus = 8

	// IsUnclassifiedEntryStatus indicates a filesystem entry whose mode no
	// other EntryStatus constant classifies.
	IsUnclassifiedEntryStatus EntryStatus = 9
)

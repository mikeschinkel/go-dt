package dt

import (
	"os"
)

// EntryModeFunc is a callback that returns the desired file mode for a destination path.
// Return 0 to preserve the source file's permissions.
type EntryModeFunc func(ep EntryPath) os.FileMode

// CopyOptions contains options for the copy operation
type CopyOptions struct {
	Force        bool          // Overwrite existing files
	DestModeFunc EntryModeFunc // Permission callback (nil = preserve source permissions)
}

// UnixModeFunc provides standard Unix permissions (0755 for directories, 0644 for files)
func UnixModeFunc(ep EntryPath) (mode os.FileMode) {
	var info os.FileInfo
	var err error

	// Check if path exists and is a directory
	info, err = ep.Stat()
	if err == nil && info.IsDir() {
		mode = 0755 // rwxr-xr-x for directories
		goto end
	}

	// Check if path ends with / (directory indicator)
	if ep.HasSuffix("/") {
		mode = 0755
		goto end
	}

	// Default file mode
	mode = 0644 // rw-r--r-- for files

end:
	return mode
}

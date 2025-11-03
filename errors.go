package dt

import (
	"errors"
)

var (
	ErrPathIsDir               = errors.New("path is a directory")
	ErrPathIsFile              = errors.New("path is a file (not a directory)")
	ErrFileDoesNotExist        = errors.New("file does not exist")
	ErrFileExists              = errors.New("file exists")
	ErrInvalidEntryStatus      = errors.New("invalid entry status")
	ErrUnclassifiedEntryStatus = errors.New("unclassified entry status")
	ErrUnsupportedEntryType    = errors.New("unsupported entry type")
)

var (
	ErrInvalidPathSegment = errors.New("invalid path segment")

	// ErrInvalidForOpen is used when ValidPath()==false
	ErrInvalidForOpen = errors.New("invalid for open")
)

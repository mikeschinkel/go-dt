package de

import (
	"errors"
)

var (
	ErrFailedTypeAssertion        = errors.New("failed type assertion")
	ErrFailedReadingSymlink       = errors.New("failed reading symlink")
	ErrFailedToLoadFile           = errors.New("failed to load file")
	ErrFailedToCopyFile           = errors.New("failed to copy file")
	ErrFailedCreatingDirectory    = errors.New("failed to create directory")
	ErrContainsBackslash          = errors.New("contains slash ('\\')")
	ErrContainsSlash              = errors.New("contains slash ('/')")
	ErrEmpty                      = errors.New("empty")
	ErrInvalidPercentEncoding     = errors.New("invalid percent encoding")
	ErrTooLong                    = errors.New("too long")
	ErrUnspecified                = errors.New("unspecified")
	ErrInvalid                    = errors.New("invalid")
	ErrInvalidfileSystemEntryType = errors.New("invalid file system entry type")
	ErrControlCharacter           = errors.New("control character")
	ErrInvalidCharacter           = errors.New("invalid charnacter")
	ErrTrailingSpace              = errors.New("trailing space")
	ErrTrailingPeriod             = errors.New("trailing period")
	ErrReservedDeviceName         = errors.New("reserved device name")
	ErrNotFileOrDirectory         = errors.New("not a file or directory")
	ErrNotADirectory              = errors.New("not a directory")
	ErrIsAFile                    = errors.New("is a file")
	ErrIsADirectory               = errors.New("is a directory")
)

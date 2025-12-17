package dt

import (
	"errors"
)

var (
	ErrPathIsDir               = errors.New("path is a directory")
	ErrPathIsFile              = errors.New("path is a file (not a directory)")
	ErrDoesNotExist            = errors.New("does not exist")
	ErrExists                  = errors.New("exists")
	ErrDirDoesNotExist         = errors.New("dir does not exist")
	ErrDirExists               = errors.New("dir exists")
	ErrFileNotExists           = errors.New("file does not exist")
	ErrFileExists              = errors.New("file exists")
	ErrInvalidEntryStatus      = errors.New("invalid entry status")
	ErrUnclassifiedEntryStatus = errors.New("unclassified entry status")
	ErrUnsupportedEntryType    = errors.New("unsupported entry type")
)

var (
	ErrInvalidPathSegment = errors.New("invalid path segment")
	ErrInvalidURLSegment  = errors.New("invalid URL segment")
	ErrInvalidURLSegments = errors.New("invalid URL segments")
	ErrInvalidIdentifier  = errors.New("invalid identifier")

	// ErrInvalidForOpen is used when ValidPath()==false
	ErrInvalidForOpen = errors.New("invalid for open")
)
var (
	ErrFileSystem                      = errors.New("file system error")
	ErrFailedReadingSymlink            = errors.New("failed reading symlink")
	ErrFailedToLoadFile                = errors.New("failed to load file")
	ErrFailedToCopyFile                = errors.New("failed to copy file")
	ErrFailedToMakeDirectory           = errors.New("failed to make directory")
	ErrFailedtoCreateTempFile          = errors.New("failed to create temp file")
	ErrFailedtoCreateFile              = errors.New("failed to create file")
	ErrContainsBackslash               = errors.New("contains slash ('\\')")
	ErrContainsSlash                   = errors.New("contains slash ('/')")
	ErrEmpty                           = errors.New("cannot be empty")
	ErrInvalidPercentEncoding          = errors.New("invalid percent encoding")
	ErrTooLong                         = errors.New("too long")
	ErrUnspecified                     = errors.New("unspecified")
	ErrInvalid                         = errors.New("invalid")
	ErrInvalidfileSystemEntryType      = errors.New("invalid file system entry type")
	ErrControlCharacter                = errors.New("control character")
	ErrInvalidCharacter                = errors.New("invalid charnacter")
	ErrTrailingSpace                   = errors.New("trailing space")
	ErrTrailingPeriod                  = errors.New("trailing period")
	ErrReservedDeviceName              = errors.New("reserved device name")
	ErrNotFileOrDirectory              = errors.New("not a file or directory")
	ErrNotDirectory                    = errors.New("not a directory")
	ErrIsAFile                         = errors.New("is a file")
	ErrIsADirectory                    = errors.New("is a directory")
	ErrCannotDetermineWorkingDirectory = errors.New("cannot determine working directory")
)
var (
	ErrValueIsNil          = errors.New("value is nil")
	ErrInterfaceValueIsNil = errors.New("interface value is nil")
)
var (
	ErrConnectFailed          = errors.New("failed to connect to database")
	ErrInvalidConnectString   = errors.New("invalid connection string")
	ErrFailedToPingDatabase   = errors.New("failed to ping database")
	ErrFailedToOpenDatabase   = errors.New("failed to open database")
	ErrFailedToExecuteQueries = errors.New("failed to execute query(s)")
)

var (
	ErrFlagIsRequired       = errors.New("flag is required")
	ErrInvalidDuplicateFlag = errors.New("invalid duplicate flag")
	ErrInvalidFlagName      = errors.New("invalid flag name")
	ErrFlagValidationFailed = errors.New("flag validation failed")
)
var (
	ErrUnexpectedError = errors.New("unexpected error")
	ErrInternalError   = errors.New("internal error")
)

var (
	ErrAccessingWorkingDir    = errors.New("error accessing working dir")
	ErrAccessingUserConfigDir = errors.New("error accessing user config dir")
	ErrAccessingCLIConfigDir  = errors.New("error accessing CLI config dir")
	ErrAccessingUserHomeDir   = errors.New("error accessing user home dir")
	ErrAccessingUserCacheDir  = errors.New("error accessing user cache dir")
)

var ErrNotTildePath = errors.New("not a tilde-prefixed path")
var ErrInvalidPathSeparator = errors.New("invalid path separator")

package dt

import (
	"errors"
)

var (
	ErrPathIsDir               = errors.New("path is a directory")
	ErrPathIsFile              = errors.New("path is a file (not a directory)")
	ErrNotExist                = errors.New("does not exist")
	ErrExists                  = errors.New("exists")
	ErrDirNotExist             = errors.New("dir does not exist")
	ErrDirExists               = errors.New("dir exists")
	ErrFileNotExist            = errors.New("file does not exist")
	ErrFileExists              = errors.New("file exists")
	ErrInvalidEntryStatus      = errors.New("invalid entry status")
	ErrUnclassifiedEntryStatus = errors.New("unclassified entry status")
	ErrUnsupportedEntryType    = errors.New("unsupported entry type")
)

var (
	ErrInvalidFilepath     = errors.New("invalid filepath")
	ErrInvalidRelFilepath  = errors.New("invalid relative filepath")
	ErrInvalidRelDirPath   = errors.New("invalid relative directory segment")
	ErrInvalidPathSegments = ErrInvalidRelDirPath
	ErrInvalidPathSegment  = errors.New("invalid path segment")
	ErrInvalidURLSegment   = errors.New("invalid URL segment")
	ErrInvalidURLSegments  = errors.New("invalid URL segments")
	ErrInvalidIdentifier   = errors.New("invalid identifier")
	ErrFailedTypeAssertion = errors.New("failed type assertion")

	// ErrInvalidForOpen is used when ValidPath()==false
	ErrInvalidForOpen = errors.New("invalid for open")
)
var (
	ErrFileSystem                      = errors.New("file system error")
	ErrFileStat                        = errors.New("file stat error")
	ErrFailedReadingSymlink            = errors.New("failed reading symlink")
	ErrFailedToLoadFile                = errors.New("failed to load file")
	ErrFailedToCopyFile                = errors.New("failed to copy file")
	ErrFailedToSaveFile                = errors.New("failed to save file")
	ErrFailedToRemoveFile              = errors.New("failed to remove file")
	ErrFailedToReadFile                = errors.New("failed to read file")
	ErrFailedToOpenFile                = errors.New("failed to open file")
	ErrFailedToWriteToFile             = errors.New("failed to write to file")
	ErrFailedToMakeDirectory           = errors.New("failed to make directory")
	ErrFailedtoCreateTempFile          = errors.New("failed to create temp file")
	ErrFailedtoCreateFile              = errors.New("failed to create file")
	ErrFailedtoRemoveFile              = errors.New("failed to remove file")
	ErrFailedtoCreateDir               = errors.New("failed to create directory")
	ErrContainsBackslash               = errors.New("contains slash ('\\')")
	ErrContainsSlash                   = errors.New("contains slash ('/')")
	ErrEmpty                           = errors.New("cannot be empty")
	ErrInvalidPercentEncoding          = errors.New("invalid percent encoding")
	ErrTooShort                        = errors.New("too short")
	ErrTooLong                         = errors.New("too long")
	ErrUnspecified                     = errors.New("unspecified")
	ErrInvalid                         = errors.New("invalid")
	ErrInvalidDirectory                = errors.New("invalid directory")
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
	ErrFailedToExpandPath              = errors.New("failed to expand path")
	ErrFailedToEnsureDir               = errors.New("failed to ensure directory")
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
	ErrNotImplemented  = errors.New("not implemented")
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
var ErrFailedToUnmarshalJSON = errors.New("failed to unmarshal JSON")
var ErrFailedToMarshalJSON = errors.New("failed to marshal JSON")

// Sentinels for the host, IP, DNS and email address types. Each type has an
// ErrInvalid<Type> sentinel; the detail sentinels below say which rule failed.
var (
	ErrInvalidIPAddress           = errors.New("invalid IP address")
	ErrInvalidIPv4Address         = errors.New("invalid IPv4 address")
	ErrInvalidIPv6Address         = errors.New("invalid IPv6 address")
	ErrInvalidScopedIPv6Address   = errors.New("invalid scoped IPv6 address")
	ErrInvalidHostDomainName      = errors.New("invalid host domain name")
	ErrInvalidSingleLabelHostname = errors.New("invalid single-label host name")
	ErrInvalidInternetDomain      = errors.New("invalid internet domain")
	ErrInvalidHostname            = errors.New("invalid hostname")
	ErrInvalidDNSName             = errors.New("invalid DNS name")
	ErrInvalidHostPort            = errors.New("invalid host and port")
	ErrInvalidPortNumber          = errors.New("invalid port number")
	ErrInvalidServiceName         = errors.New("invalid service name")
	ErrInvalidEmailLocalPart      = errors.New("invalid email local part")
	ErrInvalidEmailAddress        = errors.New("invalid email address")
	ErrInvalidEmailDisplayName    = errors.New("invalid email display name")
	ErrInvalidMailbox             = errors.New("invalid mailbox")

	// ErrUnsupportedSyntax marks input that some RFC allows but the type
	// deliberately rejects. The type's doc comment says what and why.
	ErrUnsupportedSyntax = errors.New("allowed by an RFC but not supported by this type; see the type's doc comment")

	// ErrNotYetSupported is returned by the Parse functions of types that
	// are declared but not yet implemented. The type's doc comment records
	// the open questions.
	ErrNotYetSupported = errors.New("not yet supported; see the type's doc comment")

	// ErrSMTPUTF8NotAllowed is returned when an email address has a UTF-8
	// local part (RFC 6531, RFC 6532) and the caller did not set
	// AllowSMTPUTF8.
	ErrSMTPUTF8NotAllowed = errors.New("UTF-8 local part requires SMTPUTF8, which AllowSMTPUTF8 did not allow")

	ErrIPv6Zone              = errors.New("IPv6 zone index")
	ErrNotIPv4Address        = errors.New("not an IPv4 address")
	ErrNotIPv6Address        = errors.New("not an IPv6 address")
	ErrEmptyLabel            = errors.New("empty label")
	ErrLabelTooLong          = errors.New("label longer than 63 octets")
	ErrNonLDHLabel           = errors.New("label is not letters, digits and interior hyphens")
	ErrNumericFinalLabel     = errors.New("final label is all digits")
	ErrRootDot               = errors.New("trailing root dot")
	ErrUnicodeHostname       = errors.New("non-ASCII host name needs an ASCIIToer to convert it to A-labels")
	ErrSingleLabel           = errors.New("single label")
	ErrSpecialUseDomain      = errors.New("under a special-use domain suffix")
	ErrRejectedSuffix        = errors.New("under a rejected suffix")
	ErrUnlistedPublicSuffix  = errors.New("under an unlisted public suffix")
	ErrIsPublicSuffix        = errors.New("is itself a public suffix")
	ErrIsIPAddress           = errors.New("is an IP address")
	ErrMissingAtSign         = errors.New("missing '@'")
	ErrGeneralAddressLiteral = errors.New("general address literal")
	ErrInvalidUTF8           = errors.New("invalid UTF-8")
	ErrAngleBracket          = errors.New("angle bracket")
)

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

// Sentinels for the git types. Each type has an ErrInvalidGit<Type> sentinel;
// the rule sentinels below say which rule failed. The ref rules come from
// git-check-ref-format(1); every Parse, Validate and To* error wraps one type
// sentinel and one rule sentinel.
var (
	// ErrInvalidGitRef marks a value that is not a full git refname under refs/.
	ErrInvalidGitRef = errors.New("invalid git ref")
	// ErrInvalidGitBranchName marks a value git-check-ref-format(1) --branch refuses.
	ErrInvalidGitBranchName = errors.New("invalid git branch name")
	// ErrInvalidGitTagName marks a value git refuses as a tag name.
	ErrInvalidGitTagName = errors.New("invalid git tag name")
	// ErrInvalidGitRemoteName marks a value git refuses as a remote name.
	ErrInvalidGitRemoteName = errors.New("invalid git remote name")
	// ErrInvalidGitRevision marks a value that fails GitRevision's basic checks.
	ErrInvalidGitRevision = errors.New("invalid git revision")
	// ErrInvalidGitObjectID marks a value that is not a full SHA-1 or SHA-256 object name.
	ErrInvalidGitObjectID = errors.New("invalid git object ID")
	// ErrInvalidGitShortObjectID marks a value that is not 4 to 64 lowercase hex digits.
	ErrInvalidGitShortObjectID = errors.New("invalid git short object ID")
	// ErrNotUnderRefs marks a GitRef that does not begin with refs/, such as HEAD.
	ErrNotUnderRefs = errors.New("not under refs/")
	// ErrNotBranchRef marks a ref GitRef.ToBranchName cannot convert: it is not under refs/heads/.
	ErrNotBranchRef = errors.New("not under refs/heads/")
	// ErrNotTagRef marks a ref GitRef.ToTagName cannot convert: it is not under refs/tags/.
	ErrNotTagRef = errors.New("not under refs/tags/")
	// ErrEmptyRefComponent breaks git-check-ref-format(1) rule 6: no leading or trailing '/', and no '//'.
	ErrEmptyRefComponent = errors.New("empty path component: leading, trailing or doubled '/'")
	// ErrRefComponentLeadingDot breaks git-check-ref-format(1) rule 1: no component begins with '.'.
	ErrRefComponentLeadingDot = errors.New("path component begins with '.'")
	// ErrRefComponentLockSuffix breaks git-check-ref-format(1) rule 1: no component ends with '.lock'.
	ErrRefComponentLockSuffix = errors.New("path component ends with '.lock'")
	// ErrDoubleDot breaks git-check-ref-format(1) rule 3: no '..' anywhere.
	ErrDoubleDot = errors.New("contains '..'")
	// ErrAtBrace breaks git-check-ref-format(1) rule 8: no '@{' anywhere.
	ErrAtBrace = errors.New("contains '@{'")
	// ErrLoneAt is commented out, with its check in checkGitRefFormat, because
	// no git type can reach it. git-check-ref-format(1) rule 9 rejects a refname
	// that is exactly '@', but every refname these types check begins with
	// refs/ (GitRef) or is built under refs/heads/, refs/tags/ or refs/remotes/,
	// so it is never '@'. Uncomment both if a type ever checks a refname that
	// need not begin with refs/.
	//
	// ErrLoneAt breaks git-check-ref-format(1) rule 9: the name is not the single character '@'.
	// ErrLoneAt = errors.New("is '@'")
	// ErrLeadingDash marks a value git would read as an option; check-ref-format(1) --branch refuses it.
	ErrLeadingDash = errors.New("begins with '-'")
	// ErrIsHEAD marks a branch or tag name of HEAD, which git-check-ref-format(1) --branch refuses.
	ErrIsHEAD = errors.New("is 'HEAD'")
	// ErrSurroundingWhitespace marks a revision with a leading or trailing space.
	ErrSurroundingWhitespace = errors.New("leading or trailing whitespace")
	// ErrNonHexDigit marks an object name with a character that is not a hex digit.
	ErrNonHexDigit = errors.New("contains a non-hex digit")
	// ErrUppercaseHex marks an object name with uppercase hex, which git accepts but never prints.
	ErrUppercaseHex = errors.New("uppercase hex digit; set AllowUppercaseHex to accept it")
	// ErrObjectIDLength marks an object name that is neither 40 (SHA-1) nor 64 (SHA-256) hex digits.
	ErrObjectIDLength = errors.New("not 40 or 64 hex digits")
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
	// ErrInvalidIPAddress is returned when input is not an IPAddress.
	ErrInvalidIPAddress = errors.New("invalid IP address")

	// ErrInvalidIPv4Address is returned when input is not an IPv4Address
	// (RFC 3986 IPv4address).
	ErrInvalidIPv4Address = errors.New("invalid IPv4 address")

	// ErrInvalidIPv6Address is returned when input is not an IPv6Address
	// (RFC 4291).
	ErrInvalidIPv6Address = errors.New("invalid IPv6 address")

	// ErrInvalidScopedIPv6Address is returned when input is not a
	// ScopedIPv6Address (RFC 4007, RFC 6874).
	ErrInvalidScopedIPv6Address = errors.New("invalid scoped IPv6 address")

	// ErrInvalidHostDomainName is returned when input is not a
	// HostDomainName (RFC 1123 section 2.1).
	ErrInvalidHostDomainName = errors.New("invalid host domain name")

	// ErrInvalidSingleLabelHostname is returned when input is not a
	// SingleLabelHostname (RFC 6762).
	ErrInvalidSingleLabelHostname = errors.New("invalid single-label host name")

	// ErrInvalidInternetDomain is returned when input is not an
	// InternetDomain.
	ErrInvalidInternetDomain = errors.New("invalid internet domain")

	// ErrInvalidHostname is returned when input is not a Hostname (RFC 3986
	// host without brackets).
	ErrInvalidHostname = errors.New("invalid hostname")

	// ErrInvalidDNSName is returned when input is not a DNSName (RFC 9499
	// domain name).
	ErrInvalidDNSName = errors.New("invalid DNS name")

	// ErrInvalidHostPort is returned when input is not a HostPort
	// (net.SplitHostPort's hostport).
	ErrInvalidHostPort = errors.New("invalid host and port")

	// ErrInvalidPortNumber is returned when input is not a PortNumber
	// (RFC 6335 port number).
	ErrInvalidPortNumber = errors.New("invalid port number")

	// ErrInvalidServiceName is returned when input is not a ServiceName
	// (RFC 6335 service name).
	ErrInvalidServiceName = errors.New("invalid service name")

	// ErrInvalidEmailLocalPart is returned when input is not an
	// EmailLocalPart (RFC 5322 local-part).
	ErrInvalidEmailLocalPart = errors.New("invalid email local part")

	// ErrInvalidEmailAddress is returned when input is not an EmailAddress
	// (RFC 5322 addr-spec).
	ErrInvalidEmailAddress = errors.New("invalid email address")

	// ErrInvalidEmailDisplayName is returned when input is not an
	// EmailDisplayName (RFC 5322 display-name).
	ErrInvalidEmailDisplayName = errors.New("invalid email display name")

	// ErrInvalidMailbox is returned when input is not a Mailbox (RFC 5322
	// mailbox).
	ErrInvalidMailbox = errors.New("invalid mailbox")

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

	// ErrIPv6Zone means the address carries a zone index (RFC 4007 section
	// 11), as in "fe80::1%en0".
	ErrIPv6Zone = errors.New("IPv6 zone index")

	// ErrNotIPv4Address means an IP address parsed but is IPv6.
	ErrNotIPv4Address = errors.New("not an IPv4 address")

	// ErrNotIPv6Address means an IP address parsed but is IPv4.
	ErrNotIPv6Address = errors.New("not an IPv6 address")

	// ErrEmptyLabel means a name has an empty label, as in "a..com" (RFC
	// 1034 section 3.1), or a dot-atom has an empty atom (RFC 5322 section
	// 3.2.3).
	ErrEmptyLabel = errors.New("empty label")

	// ErrLabelTooLong means a label exceeds 63 octets (RFC 1035 section
	// 2.3.4).
	ErrLabelTooLong = errors.New("label longer than 63 octets")

	// ErrNonLDHLabel means a label holds a character other than a letter,
	// digit or hyphen, or begins or ends with a hyphen (RFC 1034 section 3.5
	// as amended by RFC 1123 section 2.1).
	ErrNonLDHLabel = errors.New("label is not letters, digits and interior hyphens")

	// ErrNumericFinalLabel means a name's final label is all digits, which
	// RFC 1123 section 2.1 rules out so a name cannot pass for an IPv4
	// address.
	ErrNumericFinalLabel = errors.New("final label is all digits")

	// ErrRootDot means a name ends with the root label's dot (RFC 1034
	// section 3.1), which a host name does not carry.
	ErrRootDot = errors.New("trailing root dot")

	// ErrUnicodeHostname means a host name holds non-ASCII text (RFC 5890
	// U-labels) and no IDNAConverter was given to convert it to A-labels.
	ErrUnicodeHostname = errors.New("non-ASCII host name needs an IDNAConverter to convert it to A-labels")

	// ErrSingleLabel means a name has one label where InternetDomain needs
	// two or more.
	ErrSingleLabel = errors.New("single label")

	// ErrSpecialUseDomain means a name is under a suffix reserved from the
	// public DNS (RFC 6761, RFC 6762, RFC 8375).
	ErrSpecialUseDomain = errors.New("under a special-use domain suffix")

	// ErrRejectedSuffix means a name is under one of the caller's
	// ParseInternetDomainArgs.RejectedSuffixes.
	ErrRejectedSuffix = errors.New("under a rejected suffix")

	// ErrUnlistedPublicSuffix means the PublicSuffixClassifier does not list
	// the name's suffix, as with "foo.cromulent".
	ErrUnlistedPublicSuffix = errors.New("under an unlisted public suffix")

	// ErrIsPublicSuffix means the name is itself a public suffix, such as
	// "co.uk" or "github.io".
	ErrIsPublicSuffix = errors.New("is itself a public suffix")

	// ErrIsIPAddress means a Hostname holds an IP address, which has no
	// DNSName of its own.
	ErrIsIPAddress = errors.New("is an IP address")

	// ErrMissingAtSign means an email address has no "@" (RFC 5322
	// addr-spec).
	ErrMissingAtSign = errors.New("missing '@'")

	// ErrGeneralAddressLiteral means an email host is RFC 5321 section
	// 4.1.3's General-address-literal ("[tag:content]"), which EmailAddress
	// does not support.
	ErrGeneralAddressLiteral = errors.New("general address literal")

	// ErrInvalidUTF8 means text is not valid UTF-8 (RFC 3629), which RFC 6532
	// requires of internationalized email.
	ErrInvalidUTF8 = errors.New("invalid UTF-8")

	// ErrAngleBracket means a display name holds "<" or ">", or a mailbox's
	// angle brackets do not pair (RFC 5322 section 3.4 angle-addr).
	ErrAngleBracket = errors.New("angle bracket")
)

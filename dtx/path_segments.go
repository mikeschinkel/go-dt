package dtx

import (
	"errors"
	"net/url"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-dt/de"
)

// ParseOSPathSegment parses a string to returns an OS-specific Windows file system path segment
func ParseOSPathSegment(s string) (ps dt.PathSegment, err error) {
	var osPS any
	switch runtime.GOOS {
	case "windows":
		osPS, err = ParseWindowsPathSegment(s)
	case "darwin":
		osPS, err = ParseDarwinPathSegment(s)
	default:
		osPS, err = ParseOSPathSegment(s)
	}
	switch t := osPS.(type) {
	case string:
		ps = dt.PathSegment(t)
	case WindowsPathSegment:
		ps = dt.PathSegment(t)
	case DarwinPathSegment:
		ps = dt.PathSegment(t)
	case LinuxPathSegment:
		ps = dt.PathSegment(t)
	}
	return ps, err
}

type URLPathSegment dt.PathSegment

func (s URLPathSegment) String() string { return string(s) }

// ParseURLPathSegment parses a string to return s RFC 3986 URL path segment
// Requirements (simple & practical): - non-empty - no '/' (segment only) - must
// be valid percent-encoding (PathUnescape succeeds) - length ≤ 255 runes
// (conservative)
func ParseURLPathSegment(s string) (ps URLPathSegment, err error) {
	if s == "" {
		err = dt.NewErr(dt.ErrEmpty)
		goto end
	}
	if strings.Contains(s, "/") {
		err = dt.NewErr(dt.ErrContainsSlash)
		goto end
	}
	_, err = url.PathUnescape(s)
	if err != nil {
		err = dt.NewErr(dt.ErrInvalidPercentEncoding)
		goto end
	}
	{
		n := utf8.RuneCountInString(s)
		if n > 255 {
			err = dt.NewErr(
				dt.ErrTooLong,
				"length", n,
			)
			goto end
		}
	}
	ps = URLPathSegment(s)
end:
	if err != nil {
		err = dt.WithErr(err,
			dt.ErrInvalidPathSegment,
		)
	}
	return ps, err
}

type DarwinPathSegment dt.PathSegment

func (s DarwinPathSegment) String() string { return string(s) }

// ParseDarwinPathSegment parses a string to return s macOS file system path segment
func ParseDarwinPathSegment(s string) (ps DarwinPathSegment, err error) {
	err = validateNix(s)
	if err != nil {
		goto end
	}
	ps = DarwinPathSegment(s)
end:
	return ps, err
}

type LinuxPathSegment dt.PathSegment

func (s LinuxPathSegment) String() string { return string(s) }

// ParseLinuxPathSegment parses a string to return s Linux file system path segment
func ParseLinuxPathSegment(s string) (ps LinuxPathSegment, err error) {
	err = validateNix(s)
	if err != nil {
		goto end
	}
	ps = LinuxPathSegment(s)
end:
	return ps, err
}

type WindowsPathSegment dt.PathSegment

func (s WindowsPathSegment) String() string { return string(s) }

// ParseWindowsPathSegment parses a string to returns a Windows file system path segment
// Requirements (simplified but correct for common cases):
// - non-empty
// - no: < > : " / \ | ? * or control chars (0x00–0x1F)
// - no trailing space or dot
// - not a reserved device name (CON, PRN, AUX, NUL, COM1–9, LPT1–9) before first dot
// - length ≤ 255 runes (conservative)
func ParseWindowsPathSegment(s string) (ps WindowsPathSegment, err error) {
	if s == "" {
		err = dt.NewErr(dt.ErrEmpty)
		goto end
	}
	for _, r := range s {
		if r < 0x20 {
			err = dt.NewErr(
				dt.ErrControlCharacter,
				"ascii_value", int(r),
			)
			goto end
		}
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			err = dt.NewErr(
				dt.ErrInvalidCharacter,
				"character", string(r),
			)
			goto end
		}
	}
	if strings.HasSuffix(s, " ") {
		err = dt.NewErr(dt.ErrTrailingSpace)
	}
	if strings.HasSuffix(s, ".") {
		err = dt.NewErr(dt.ErrTrailingPeriod)
	}
	{
		base := s
		if i := strings.IndexRune(s, '.'); i >= 0 {
			base = s[:i]
		}
		switch strings.ToUpper(base) {
		case "CON", "PRN", "AUX", "NUL",
			"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
			"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
			err = dt.NewErr(
				dt.ErrReservedDeviceName,
				"device_name", base,
			)
			goto end
		}
	}
	{
		n := utf8.RuneCountInString(s)
		if n > 255 {
			err = dt.NewErr(
				dt.ErrTooLong,
				"length", n,
			)
			goto end
		}
	}

	ps = WindowsPathSegment(s)
end:
	if err != nil {
		err = dt.WithErr(err,
			dt.ErrInvalidPathSegment,
		)
	}
	return ps, err
}

// validateNix validates Unix-like path segments (Darwin & Linux share rules)
// Requirements:
// - non-empty
// - no '/' or NUL
// - optionally forbid "." and ".." if you don't want special dirs
func validateNix(s string) error {
	if s == "" {
		return errors.New("empty")
	}
	if strings.Contains(s, "/") {
		return errors.New("contains '/'")
	}
	if strings.IndexByte(s, 0x00) >= 0 {
		return errors.New("contains NUL")
	}
	// Uncomment if you want to forbid these:
	// if s == "." || s == ".." { return errors.New("reserved '.' or '..'") }
	if utf8.RuneCountInString(s) > 255 {
		return errors.New("too long")
	}
	return nil
}

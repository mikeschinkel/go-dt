package dt

import (
	"strconv"
	"time"
)

// ParseTimeDurationEx parses a string as EITHER a Go duration format
// (like "3s", "10m", "1h30m") OR as an integer representing seconds.
func ParseTimeDurationEx(s string) (td time.Duration, err error) {
	var seconds int
	var errs []error

	// First try parsing as integer seconds
	seconds, err = strconv.Atoi(s)
	if err == nil {
		td = time.Duration(seconds) * time.Second
		goto end
	}
	errs = append(errs, err)

	// If that fails, try parsing as a standard Go duration
	td, err = time.ParseDuration(s)
	if err == nil {
		goto end
	}
	errs = append(errs, err)

end:
	return td, CombineErrs(errs)
}

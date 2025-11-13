package dtglob

import (
	"strings"

	"github.com/mikeschinkel/go-dt"
)

// Glob represents a glob pattern string supporting wildcards
type Glob string

// Contains checks if the glob pattern contains the given substring
func (g Glob) Contains(substr string) bool {
	return strings.Contains(string(g), substr)
}

// Split splits the glob pattern on the given separator
func (g Glob) Split(sep string) []string {
	return strings.Split(string(g), sep)
}

// GlobRule represents a single copy operation
type GlobRule struct {
	From     Glob         // Pattern to match files
	To       dt.EntryPath // Destination (file or directory)
	Optional bool         // If true, no error when pattern matches nothing
}

// GlobRules is a container for multiple rules sharing a common base directory
type GlobRules struct {
	BaseDir dt.DirPath // Source directory containing files to match
	Rules   []GlobRule // Ordered list of copy rules
}

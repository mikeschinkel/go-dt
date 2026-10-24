package dt

import (
	"unicode"
	"unicode/utf8"
)

func (pss PathSegments) UpperFirst() PathSegments {
	if pss == "" {
		goto end
	}
	{
		r, n := utf8.DecodeRuneInString(string(pss))
		//goland:noinspection GoAssignmentToReceiver
		pss = PathSegments(string(unicode.ToTitle(r)) + string(pss[n:]))
	}
end:
	return pss
}

// IsRoot returns true if these path segments represent the root directory.
// Handles both "" (normalized root) and "." (filepath.Rel result for same-dir).
func (pss PathSegments) IsRoot() bool {
	return pss == "" || pss == "."
}

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

package dt

import (
	"strings"
)

const (
	// maxDomainNameLength is the longest name in text form without a root
	// dot: 255 octets on the wire (RFC 1035 section 2.3.4) less the length
	// octets and root label.
	maxDomainNameLength = 253

	// maxLabelLength is RFC 1035 section 2.3.4's limit on one label.
	maxLabelLength = 63
)

// HostDomainName holds a host name in the syntax RFC 1123 section 2.1 calls a
// "host domain name": RFC 1034 section 3.5's preferred name syntax as amended
// by RFC 1123, so a label may begin with a digit. Examples: "example.com",
// "mikes-mbp.local", "repo.test", "localhost".
//
// The name is RFC 1123's own term. go-dt uses it rather than HostName, which
// would sit confusingly beside Hostname, and it covers names that are neither
// single-label nor Internet domains, so no type is needed per subset.
// SingleLabelHostname and InternetDomain are narrower members, and convert to
// HostDomainName with a plain type conversion.
//
// HostDomainName keeps the text as given; it does not fold case.
//
// HostDomainName rejects:
//   - A label holding anything but ASCII letters, digits and hyphens, or
//     beginning or ending with a hyphen, such as "-bad-" or "_dmarc". DNS can
//     store any octet in a label (RFC 2181 section 11), but host names, email
//     and URLs cannot use them (RFC 1123, RFC 5321). DNSName accepts
//     underscored labels.
//   - A trailing root dot, as in "example.com.". The dot makes a name absolute
//     in DNS (RFC 1034 section 3.1) but is not part of a host name in email or
//     URLs. DNSName keeps it.
//   - Unicode, as in "bücher.example" (RFC 5890 U-labels). Write the name as
//     A-labels ("xn--bcher-kva.example"), or convert it with an IDNAConverter
//     through a Parse function that takes one. Without IDNA tables go-dt
//     cannot tell a real A-label from an RFC 5890 Fake A-label, so it accepts
//     any LDH label beginning "xn--".
//   - A final label made only of digits, such as "1.2.3.4", which RFC 1123
//     forbids so a name cannot be mistaken for an IPv4 address.
//   - An empty name, an empty label, a label over 63 octets or a name over
//     253 octets (RFC 1035 section 2.3.4).
//
// If you need a type that accepts what HostDomainName rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type HostDomainName string

// SingleLabelHostname holds a host domain name made of exactly one label,
// RFC 6762's "single-label host name": "localhost", "buildbox". git's default
// identity on a machine with no user.email is user@<single-label host name>.
//
// The name is RFC 6762's own phrase. go-dt does not call it LocalHostname:
// macOS's local host name ("Mikes-MBP.local") has two labels, and "local"
// suggests mDNS's .local domain.
//
// SingleLabelHostname rejects everything HostDomainName rejects, and any name
// with more than one label.
//
// If you need a type that accepts what SingleLabelHostname rejects, define
// your own type, or open an issue at
// https://github.com/mikeschinkel/go-dt/issues to propose one.
type SingleLabelHostname string

// ParseHostDomainName parses s as a host domain name. It fails with
// ErrInvalidHostDomainName, joined with ErrUnsupportedSyntax for a name some
// RFC allows but HostDomainName rejects; see HostDomainName.
func ParseHostDomainName(s string) (hdn HostDomainName, err error) {
	err = domainNameErr(s, hostDomainNameRules, ErrInvalidHostDomainName, "dt.HostDomainName")
	if err != nil {
		goto end
	}
	hdn = HostDomainName(s)
end:
	return hdn, err
}

// ParseSingleLabelHostname parses s as a host domain name of one label. It
// fails with ErrInvalidSingleLabelHostname; see SingleLabelHostname.
func ParseSingleLabelHostname(s string) (slh SingleLabelHostname, err error) {
	err = domainNameErr(s, hostDomainNameRules, ErrInvalidSingleLabelHostname, "dt.SingleLabelHostname")
	if err != nil {
		goto end
	}
	if strings.Contains(s, ".") {
		err = NewErr(ErrInvalidSingleLabelHostname, "input", s)
		goto end
	}
	slh = SingleLabelHostname(s)
end:
	return slh, err
}

// ErrKV returns hdn as error metadata under the key "host_domain_name".
func (hdn HostDomainName) ErrKV() ErrKV {
	return kv{k: "host_domain_name", v: hdn}
}

// ErrKV returns slh as error metadata under the key "single_label_hostname".
func (slh SingleLabelHostname) ErrKV() ErrKV {
	return kv{k: "single_label_hostname", v: slh}
}

// domainNameRules selects the looseness a name type allows beyond RFC 1123.
type domainNameRules struct {
	allowUnderscore bool
	allowRootDot    bool
}

var hostDomainNameRules = domainNameRules{}

// domainNameErr checks s against rules and wraps any fault in sentinel,
// joining ErrUnsupportedSyntax and the type's name when an RFC allows what
// the rules reject.
func domainNameErr(s string, rules domainNameRules, sentinel error, typeName string) (err error) {
	fault, unsupported := domainNameFault(s, rules)
	switch {
	case fault == nil:
	case unsupported:
		err = NewErr(sentinel, ErrUnsupportedSyntax, "type", typeName, "input", s, fault)
	default:
		err = NewErr(sentinel, "input", s, fault)
	}
	return err
}

// domainNameFault returns the first rule s breaks, or nil. unsupported is
// true when some RFC allows what rules reject.
func domainNameFault(s string, rules domainNameRules) (fault error, unsupported bool) {
	var labels []string
	switch {
	case s == "":
		fault = ErrEmpty
		goto end
	case !isASCII(s):
		fault, unsupported = ErrUnicodeHostname, true
		goto end
	case strings.HasSuffix(s, "."):
		if !rules.allowRootDot {
			fault, unsupported = ErrRootDot, true
			goto end
		}
		if s == "." {
			goto end
		}
		s = s[:len(s)-1]
	}
	if len(s) > maxDomainNameLength {
		fault = ErrTooLong
		goto end
	}
	labels = strings.Split(s, ".")
	for _, label := range labels {
		fault, unsupported = labelFault(label, rules)
		if fault != nil {
			goto end
		}
	}
	if isAllDigits(labels[len(labels)-1]) {
		fault = ErrNumericFinalLabel
	}
end:
	return fault, unsupported
}

// labelFault returns the rule label breaks, or nil.
func labelFault(label string, rules domainNameRules) (fault error, unsupported bool) {
	switch {
	case label == "":
		fault = ErrEmptyLabel
		goto end
	case len(label) > maxLabelLength:
		fault = ErrLabelTooLong
		goto end
	case label[0] == '-' || label[len(label)-1] == '-':
		fault, unsupported = ErrNonLDHLabel, true
		goto end
	}
	for i := 0; i < len(label); i++ {
		c := label[i]
		switch {
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z', '0' <= c && c <= '9', c == '-':
		case c == '_' && rules.allowUnderscore:
		default:
			fault, unsupported = ErrNonLDHLabel, true
			goto end
		}
	}
end:
	return fault, unsupported
}

// isASCII reports whether s holds only ASCII bytes.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// isAllDigits reports whether s is non-empty and holds only ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

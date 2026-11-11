package dt

import (
	"net/netip"
	"strings"
	"unicode/utf8"
)

const (
	// maxEmailLocalPartLength is RFC 5321 section 4.5.3.1.1's limit.
	maxEmailLocalPartLength = 64

	// ipv6LiteralTag is RFC 5321 section 4.1.3's IPv6 address literal tag.
	ipv6LiteralTag = "IPv6:"
)

// EmailLocalPart holds the part of an email address before the "@", RFC
// 5322's "local-part": a dot-atom such as "john.doe", or a quoted-string such
// as `"john doe"`, which keeps its quotes so the value stays a valid local
// part. Non-ASCII UTF-8 is allowed as RFC 6532 extends the syntax; an address
// holding it needs SMTPUTF8 (RFC 6531) to be sent.
//
// EmailLocalPart is case-sensitive, as RFC 5321 section 2.4 says a local part
// may be.
//
// EmailLocalPart rejects:
//   - RFC 5322's obsolete forms (obs-local-part), such as "a..b" or
//     "john . doe", and comments (CFWS) such as "john(comment)".
//   - Characters RFC 5321's quoted-string forbids, such as a raw control
//     character, inside quotes.
//   - Invalid UTF-8, and a local part over 64 octets (RFC 5321 section
//     4.5.3.1.1).
//
// If you need a type that accepts what EmailLocalPart rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type EmailLocalPart string

// EmailAddress holds an email address, RFC 5322's "addr-spec": an
// EmailLocalPart, "@", and a Hostname. x509.Certificate.EmailAddresses uses
// the same name.
//
// An IP host is written as RFC 5321 section 4.1.3's address literal:
// "user@[192.0.2.1]", "user@[IPv6:2001:db8::1]". Hostname() strips the
// brackets and the IPv6: tag, the way url.URL keeps "[::1]:80" in its text
// while url.URL.Hostname() returns "::1". A bare IP after the "@", as in
// "user@192.0.2.1", is accepted (Postfix accepts it with
// resolve_numeric_domain) and stored bracketed, so the stored text is
// canonical. An IP literal is stored as net/netip's canonical text; the rest
// of the text is kept as given.
//
// == compares the text exactly. Equals compares with the host folded to lower
// case, as RFC 5321 section 2.4 makes the domain case-insensitive while the
// local part may not be.
//
// The parser does not use net/mail: mail.ParseAddress accepts labels such as
// "-bad-" and a bare IP after the "@", and drops the quotes from a quoted
// local part, so `"john doe"@x.com` comes back as an address that is no longer
// valid.
//
// EmailAddress rejects everything EmailLocalPart and Hostname reject, and:
//   - RFC 5321's general address literal ("[tag:content]"). Its IANA
//     registry holds only IPv6, which EmailAddress accepts.
//   - A display name or angle brackets; Mailbox holds those.
//   - A UTF-8 local part, unless the caller sets AllowSMTPUTF8, since such an
//     address can only be sent over SMTPUTF8 (RFC 6531). It fails with
//     ErrSMTPUTF8NotAllowed.
//
// If you need a type that accepts what EmailAddress rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type EmailAddress string

// ParseEmailAddressArgs holds the options for ParseEmailAddress. A nil
// *ParseEmailAddressArgs, or a zero value, rejects a UTF-8 local part and a
// Unicode host name.
type ParseEmailAddressArgs struct {
	// ASCIIToer, when set, converts a Unicode host name to A-labels; the
	// parsed value holds the A-label form.
	ASCIIToer ASCIIToer

	// AllowSMTPUTF8 accepts a UTF-8 local part (RFC 6531, RFC 6532).
	AllowSMTPUTF8 bool
}

// ParseEmailLocalPart parses s as an email local part. It accepts UTF-8; the
// SMTPUTF8 opt-in belongs to ParseEmailAddress, which knows the address is
// meant to be sent. It fails with ErrInvalidEmailLocalPart; see
// EmailLocalPart.
func ParseEmailLocalPart(s string) (lp EmailLocalPart, err error) {
	fault := emailLocalPartFault(s)
	if fault != nil {
		err = NewErr(ErrInvalidEmailLocalPart, "input", s, fault)
		goto end
	}
	lp = EmailLocalPart(s)
end:
	return lp, err
}

// ParseEmailAddress parses s as an email address. args may be nil; see
// ParseEmailAddressArgs. ParseEmailAddress(`"john doe"@example.com`) returns
// the input unchanged. It fails with ErrInvalidEmailAddress, joined with
// ErrUnsupportedSyntax for input some RFC allows but EmailAddress rejects;
// see EmailAddress.
func ParseEmailAddress(s string, args *ParseEmailAddressArgs) (ea EmailAddress, err error) {
	var at int
	var lp EmailLocalPart
	var host string

	if args == nil {
		args = &ParseEmailAddressArgs{}
	}
	at = strings.LastIndexByte(s, '@')
	if at < 0 {
		err = NewErr(ErrInvalidEmailAddress, ErrMissingAtSign, "input", s)
		goto end
	}
	lp, err = ParseEmailLocalPart(s[:at])
	if err != nil {
		err = NewErr(ErrInvalidEmailAddress, "input", s, err)
		goto end
	}
	if !args.AllowSMTPUTF8 && !isASCII(string(lp)) {
		err = NewErr(ErrInvalidEmailAddress, ErrSMTPUTF8NotAllowed, "input", s)
		goto end
	}
	host, err = emailHostText(s[at+1:], args.ASCIIToer)
	if err != nil {
		err = NewErr(ErrInvalidEmailAddress, "input", s, err)
		goto end
	}
	ea = EmailAddress(string(lp) + "@" + host)
end:
	return ea, err
}

// LocalPart returns the part of ea before the "@".
func (ea EmailAddress) LocalPart() EmailLocalPart {
	lp, _ := ea.split()
	return EmailLocalPart(lp)
}

// Hostname returns the part of ea after the "@", with an address literal's
// brackets and IPv6: tag removed: "user@[IPv6:2001:db8::1]" returns
// "2001:db8::1".
func (ea EmailAddress) Hostname() Hostname {
	_, host := ea.split()
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
		if hasFoldPrefix(host, ipv6LiteralTag) {
			host = host[len(ipv6LiteralTag):]
		}
	}
	return Hostname(host)
}

// Equals reports whether ea and other name the same mailbox: the local parts
// match exactly and the hosts match ignoring ASCII case (RFC 5321 section
// 2.4).
func (ea EmailAddress) Equals(other EmailAddress) bool {
	lp1, host1 := ea.split()
	lp2, host2 := other.split()
	return lp1 == lp2 && strings.EqualFold(host1, host2)
}

// RequiresSMTPUTF8 reports whether ea holds non-ASCII text, so sending it
// needs the SMTPUTF8 extension (RFC 6531).
func (ea EmailAddress) RequiresSMTPUTF8() bool {
	return !isASCII(string(ea))
}

// Validate returns nil when ea is a valid email address, or the error
// ParseEmailAddress would return for it. It accepts a UTF-8 local part, since
// a value does not record how it was parsed, and rejects a Unicode host name,
// which ParseEmailAddress would have converted to A-labels.
func (ea EmailAddress) Validate() error {
	_, err := ParseEmailAddress(string(ea), &ParseEmailAddressArgs{AllowSMTPUTF8: true})
	return err
}

// ErrKV returns ea as error metadata under the key "email_address".
func (ea EmailAddress) ErrKV() ErrKV {
	return kv{k: "email_address", v: ea}
}

// ErrKV returns lp as error metadata under the key "email_local_part".
func (lp EmailLocalPart) ErrKV() ErrKV {
	return kv{k: "email_local_part", v: lp}
}

// split returns the text before and after the last "@".
func (ea EmailAddress) split() (localPart, host string) {
	at := strings.LastIndexByte(string(ea), '@')
	if at < 0 {
		return string(ea), ""
	}
	return string(ea)[:at], string(ea)[at+1:]
}

// MakeEmailAddress builds an EmailAddress from typed parts, failing when the
// result is not valid.
//
// It ships commented out: ParseEmailAddress already covers the need, and a
// Make* constructor is one more way to do the same thing. If you have a use
// that ParseEmailAddress serves poorly, open a pull request at
// https://github.com/mikeschinkel/go-dt that enables it and explains the use.
//
//func MakeEmailAddress(lp EmailLocalPart, h Hostname, args *ParseEmailAddressArgs) (ea EmailAddress, err error) {
//	host := string(h)
//	if _, ipErr := netip.ParseAddr(host); ipErr == nil {
//		host = "[" + host + "]"
//	}
//	return ParseEmailAddress(string(lp)+"@"+host, args)
//}

// emailLocalPartFault returns the rule s breaks as a local part, or nil.
func emailLocalPartFault(s string) (fault error) {
	switch {
	case s == "":
		fault = ErrEmpty
	case len(s) > maxEmailLocalPartLength:
		fault = ErrTooLong
	case !utf8.ValidString(s):
		fault = ErrInvalidUTF8
	case s[0] == '"':
		fault = quotedLocalPartFault(s)
	default:
		fault = dotAtomFault(s)
	}
	return fault
}

// dotAtomFault checks RFC 5322's dot-atom, with RFC 6532's UTF-8.
func dotAtomFault(s string) (fault error) {
	for _, atom := range strings.Split(s, ".") {
		if atom == "" {
			fault = ErrEmptyLabel
			goto end
		}
		for _, r := range atom {
			if !isAText(r) {
				fault = ErrInvalidCharacter
				goto end
			}
		}
	}
end:
	return fault
}

// quotedLocalPartFault checks RFC 5321's Quoted-string, with RFC 6532's
// UTF-8 in qtext.
func quotedLocalPartFault(s string) (fault error) {
	var escaped bool
	var inner string
	if len(s) < 2 || s[len(s)-1] != '"' {
		fault = ErrInvalidCharacter
		goto end
	}
	inner = s[1 : len(s)-1]
	for _, r := range inner {
		switch {
		case escaped:
			// quoted-pairSMTP: backslash then any printable ASCII or space
			if r < 32 || r > 126 {
				fault = ErrInvalidCharacter
				goto end
			}
			escaped = false
		case r == '\\':
			escaped = true
		case r == '"':
			fault = ErrInvalidCharacter
			goto end
		case r >= 32 && r <= 126, r >= 0x80:
			// qtextSMTP and UTF8-non-ascii
		default:
			fault = ErrControlCharacter
			goto end
		}
	}
	if escaped {
		fault = ErrInvalidCharacter
	}
end:
	return fault
}

// isAText reports whether r is RFC 5322 atext or RFC 6532 UTF8-non-ascii.
func isAText(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z', '0' <= r && r <= '9':
		return true
	case r >= 0x80:
		return true
	}
	return strings.ContainsRune("!#$%&'*+-/=?^_`{|}~", r)
}

// emailHostText returns the host part of an email address in its stored
// form: an IP bracketed as an RFC 5321 address literal, a host domain name as
// given or converted to A-labels.
func emailHostText(host string, toer ASCIIToer) (text string, err error) {
	var ip4 IPv4Address
	var ip6 IPv6Address
	var name, inner string

	switch {
	case host == "":
		err = NewErr(ErrInvalidHostname, ErrEmpty)
		goto end
	case strings.HasPrefix(host, "["):
		if !strings.HasSuffix(host, "]") {
			err = NewErr(ErrInvalidHostname, ErrInvalidCharacter, "input", host)
			goto end
		}
		inner = host[1 : len(host)-1]
		switch {
		case hasFoldPrefix(inner, ipv6LiteralTag):
			ip6, err = ParseIPv6Address(inner[len(ipv6LiteralTag):])
			text = "[" + ipv6LiteralTag + string(ip6) + "]"
		case strings.Contains(inner, ":"):
			err = NewErr(ErrInvalidHostname, ErrUnsupportedSyntax, ErrGeneralAddressLiteral, "type", "dt.EmailAddress", "input", host)
		default:
			ip4, err = ParseIPv4Address(inner)
			text = "[" + string(ip4) + "]"
		}
		goto end
	}
	if addr, ipErr := netip.ParseAddr(host); ipErr == nil {
		// A bare IP after the "@" is stored as an address literal
		if addr.Is4() {
			ip4, err = ParseIPv4Address(host)
			text = "[" + string(ip4) + "]"
			goto end
		}
		ip6, err = ParseIPv6Address(host)
		text = "[" + ipv6LiteralTag + string(ip6) + "]"
		goto end
	}
	name, err = hostToASCII(host, toer)
	if err != nil {
		err = NewErr(ErrInvalidHostname, "input", host, err)
		goto end
	}
	err = domainNameErr(name, hostDomainNameRules, ErrInvalidHostname, "dt.EmailAddress")
	text = name
end:
	if err != nil {
		text = ""
	}
	return text, err
}

// hasFoldPrefix reports whether s begins with prefix, ignoring ASCII case.
func hasFoldPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

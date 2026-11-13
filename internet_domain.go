package dt

import (
	"strings"
)

// InternetDomain is used for internet domains like google.com or
// www.example.com: a host domain name of two or more labels that is not under
// a special-use suffix, so it can name something on the public internet.
//
// ParseInternetDomain enforces that meaning; InternetDomain itself is a string
// type, so a conversion such as InternetDomain("x") is not checked. It
// converts to HostDomainName and Hostname with a plain type conversion.
//
// InternetDomain rejects everything HostDomainName rejects, and:
//   - A single-label name such as "localhost". SingleLabelHostname holds it.
//   - A name under a special-use suffix that RFCs reserve from the public
//     DNS: "test", "localhost" and "invalid" (RFC 6761), "local" (RFC 6762)
//     and "home.arpa" (RFC 8375). HostDomainName holds such names
//     ("repo.test", "mikes-mbp.local"). "example" is not on the list:
//     example.com resolves on the public internet.
//   - With a PublicSuffixClassifier, a name under a suffix the classifier
//     does not list ("foo.cromulent"), and a name that is itself a public
//     suffix ("co.uk", "github.io").
//   - A name under one of the caller's RejectedSuffixes.
//
// If you need a type that accepts what InternetDomain rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type InternetDomain string

// ParseInternetDomainArgs holds the options for ParseInternetDomain. A nil
// *ParseInternetDomainArgs, or a zero value, means syntax plus the built-in
// special-use list, with Unicode rejected.
type ParseInternetDomainArgs struct {
	// PublicSuffixClassifier, when set, rejects a name under an unlisted
	// suffix and a name that is itself a public suffix.
	PublicSuffixClassifier PublicSuffixClassifier

	// IDNAConverter, when set, uses IDNA (Internationalized Domain Names in
	// Applications; see IDNAConverter) to convert a Unicode name to A-labels
	// before the checks run; the parsed value holds the A-label form.
	IDNAConverter IDNAConverter

	// RejectedSuffixes rejects names at or under these suffixes, such as
	// "lan" or "corp", which no RFC reserves but a caller may not want.
	// Matching ignores ASCII case.
	RejectedSuffixes []HostDomainName

	// AcceptedSuffixes accepts names at or under these suffixes without
	// asking the PublicSuffixClassifier, for a suffix the classifier's list
	// does not know yet. It does not override the special-use list or
	// RejectedSuffixes. Matching ignores ASCII case.
	AcceptedSuffixes []HostDomainName
}

// specialUseSuffixes are the suffixes RFCs reserve from the public DNS.
var specialUseSuffixes = []HostDomainName{
	"test",      // RFC 6761 section 6.2
	"localhost", // RFC 6761 section 6.3
	"invalid",   // RFC 6761 section 6.4
	"local",     // RFC 6762
	"home.arpa", // RFC 8375
}

// ParseInternetDomain parses s as an internet domain. args may be nil; see
// ParseInternetDomainArgs. It fails with ErrInvalidInternetDomain, joined with
// ErrUnsupportedSyntax for a name some RFC allows but InternetDomain rejects;
// see InternetDomain.
func ParseInternetDomain(s string, args *ParseInternetDomainArgs) (d InternetDomain, err error) {
	var name, lower string
	var suffix HostDomainName
	var kind PublicSuffixKind

	if args == nil {
		args = &ParseInternetDomainArgs{}
	}
	name, err = hostToASCII(s, args.IDNAConverter)
	if err != nil {
		err = NewErr(ErrInvalidInternetDomain, "input", s, err)
		goto end
	}
	err = domainNameErr(name, hostDomainNameRules, ErrInvalidInternetDomain, "dt.InternetDomain")
	if err != nil {
		goto end
	}
	if !strings.Contains(name, ".") {
		err = NewErr(ErrInvalidInternetDomain, ErrSingleLabel, "input", s)
		goto end
	}
	lower = strings.ToLower(name)
	if hasAnySuffix(lower, specialUseSuffixes) {
		err = NewErr(ErrInvalidInternetDomain, ErrSpecialUseDomain, "input", s)
		goto end
	}
	if hasAnySuffix(lower, args.RejectedSuffixes) {
		err = NewErr(ErrInvalidInternetDomain, ErrRejectedSuffix, "input", s)
		goto end
	}
	if args.PublicSuffixClassifier == nil || hasAnySuffix(lower, args.AcceptedSuffixes) {
		goto done
	}
	suffix, kind = args.PublicSuffixClassifier.ClassifyPublicSuffix(HostDomainName(lower))
	switch {
	case kind == UnspecifiedPublicSuffixKind, kind == UnlistedPublicSuffix:
		err = NewErr(ErrInvalidInternetDomain, ErrUnlistedPublicSuffix, "input", s, "suffix", suffix)
		goto end
	case strings.EqualFold(string(suffix), lower):
		err = NewErr(ErrInvalidInternetDomain, ErrIsPublicSuffix, "input", s)
		goto end
	}
done:
	d = InternetDomain(name)
end:
	return d, err
}

// ErrKV returns d as error metadata under the key "internet_domain".
func (d InternetDomain) ErrKV() ErrKV {
	return kv{k: "internet_domain", v: d}
}

// hasAnySuffix reports whether lowerName equals, or is under, any of
// suffixes, ignoring ASCII case in the suffixes.
func hasAnySuffix(lowerName string, suffixes []HostDomainName) bool {
	for _, s := range suffixes {
		suffix := strings.ToLower(strings.Trim(string(s), "."))
		if suffix == "" {
			continue
		}
		if lowerName == suffix || strings.HasSuffix(lowerName, "."+suffix) {
			return true
		}
	}
	return false
}

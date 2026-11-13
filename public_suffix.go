package dt

// PublicSuffixKind says what kind of public suffix a PublicSuffixClassifier
// found under a domain. A public suffix is a domain under which anyone can
// register names, as the Public Suffix List (https://publicsuffix.org) records
// them: "com", "co.uk", "github.io".
//
// The constants have explicit values, not iota, so a stored value keeps its
// meaning if the list is reordered.
type PublicSuffixKind int

const (
	// UnspecifiedPublicSuffixKind is the zero value. ParseInternetDomain
	// treats it like UnlistedPublicSuffix.
	UnspecifiedPublicSuffixKind PublicSuffixKind = 0

	// UnlistedPublicSuffix means the domain's top-level label is on no list,
	// as in "foo.cromulent".
	UnlistedPublicSuffix PublicSuffixKind = 1

	// ICANNPublicSuffix means the suffix is in the list's ICANN section:
	// "com", "co.uk".
	ICANNPublicSuffix PublicSuffixKind = 2

	// PrivatePublicSuffix means the suffix is in the list's private section,
	// added by the organization that owns it: "github.io".
	PrivatePublicSuffix PublicSuffixKind = 3
)

// PublicSuffixClassifier finds the public suffix of a domain and says what
// kind it is. ParseInternetDomain uses one, when its args supply it, to reject
// a domain under an unlisted suffix and a domain that is itself a public
// suffix.
//
// go-dt ships no implementation, so its go.mod stays dependency-free. An
// implementation backed by golang.org/x/net/publicsuffix belongs in go-dt's
// dtx module, or in the caller's code.
//
// The interface pairs with ClassifyPublicSuffixFunc the way http.Handler pairs
// with http.HandlerFunc. It does not reuse net/http/cookiejar's
// PublicSuffixList: importing cookiejar pulls net/http and crypto/tls into
// every go-dt consumer, and its PublicSuffix(domain) string method cannot
// report that a suffix is unlisted. The method is named ClassifyPublicSuffix,
// not PublicSuffix, so it cannot collide with cookiejar's.
type PublicSuffixClassifier interface {
	// ClassifyPublicSuffix returns domain's public suffix and its kind. For
	// an unlisted suffix it returns the domain's top-level label and
	// UnlistedPublicSuffix. ParseInternetDomain passes a domain that has
	// passed HostDomainName's syntax check, in lower case and A-label form.
	ClassifyPublicSuffix(domain HostDomainName) (suffix HostDomainName, kind PublicSuffixKind)
}

// ClassifyPublicSuffixFunc adapts a function to PublicSuffixClassifier, as
// http.HandlerFunc adapts a function to http.Handler.
type ClassifyPublicSuffixFunc func(domain HostDomainName) (suffix HostDomainName, kind PublicSuffixKind)

// ClassifyPublicSuffix calls f(domain).
func (f ClassifyPublicSuffixFunc) ClassifyPublicSuffix(domain HostDomainName) (suffix HostDomainName, kind PublicSuffixKind) {
	return f(domain)
}

package dt

// IDNAConverter converts an internationalized domain name to the ASCII form
// that DNS, email and URLs carry.
//
// IDNA stands for "Internationalized Domain Names in Applications", the
// standard (RFC 5890 through RFC 5894) that lets a domain name contain
// non-ASCII letters, as in "bücher.example". DNS itself carries only ASCII, so
// IDNA encodes each non-ASCII label, which RFC 5890 calls a U-label, as an
// ASCII label beginning "xn--", called an A-label: "bücher.example" becomes
// "xn--bcher-kva.example". The encoding is Punycode (RFC 3492), and IDNA's
// tables also decide which Unicode characters a label may hold and how they
// are normalized first.
//
// ParseInternetDomain, ParseHostname, ParseEmailAddress and ParseMailbox take
// an IDNAConverter through their args; without one they reject a non-ASCII
// host name with ErrUnicodeHostname. go-dt ships no implementation, because
// IDNA's tables live in golang.org/x/net/idna and go-dt's go.mod stays
// dependency-free. An adapter belongs in go-dt's dtx module or in the
// caller's code.
type IDNAConverter interface {
	// ToASCIIHostDomainName returns name with every U-label converted to an
	// A-label, or an error when name is not a valid internationalized domain
	// name. The Parse function that called it still checks the result
	// against HostDomainName's syntax.
	ToASCIIHostDomainName(name string) (HostDomainName, error)
}

// ToASCIIHostDomainNameFunc adapts a function to IDNAConverter, as
// http.HandlerFunc adapts a function to http.Handler.
type ToASCIIHostDomainNameFunc func(name string) (HostDomainName, error)

// ToASCIIHostDomainName calls f(name).
func (f ToASCIIHostDomainNameFunc) ToASCIIHostDomainName(name string) (HostDomainName, error) {
	return f(name)
}

// hostToASCII converts a non-ASCII host name with conv, leaving ASCII input
// and a nil conv alone so the syntax check reports Unicode itself.
func hostToASCII(s string, conv IDNAConverter) (ascii string, err error) {
	var hdn HostDomainName
	ascii = s
	if conv == nil || isASCII(s) {
		goto end
	}
	hdn, err = conv.ToASCIIHostDomainName(s)
	ascii = string(hdn)
end:
	return ascii, err
}

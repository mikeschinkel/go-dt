package dt

import (
	"net/netip"
)

// Hostname holds the host part of a URL or email address without brackets:
// an IP address or a host domain name. It is RFC 3986's "host" rule with an
// IP literal's brackets removed, which is what url.URL.Hostname() returns, and
// what os.Hostname() returns for this machine.
//
// The name is Hostname rather than Host because it is more specific, and it
// matches url.URL.Hostname and os.Hostname. RFC 9499 notes that "host name"
// has no formal definition and is usually read to exclude IP addresses; the
// standard library includes them, as does RFC 1123's "identity of an Internet
// host", and so does Hostname.
//
// Its members, each of which converts to Hostname with a plain type
// conversion or NewHostname:
//
//	Hostname
//	├── IPAddress
//	│   ├── IPv4Address
//	│   └── IPv6Address
//	└── HostDomainName
//	    ├── SingleLabelHostname
//	    └── InternetDomain
//
// Hostname and DNSName are not alternatives: they hold different sets.
// ToDNSName and DNSName.ToHostname convert between them where they overlap.
//
// Hostname rejects everything IPAddress and HostDomainName reject, including
// a bracketed IP such as "[::1]" and an IPv6 zone. ParseHostname accepts a
// Unicode name only through an ASCIIToer, and stores the A-label form.
//
// If you need a type that accepts what Hostname rejects, define your own type,
// or open an issue at https://github.com/mikeschinkel/go-dt/issues to propose
// one.
type Hostname string

// ParseHostnameArgs holds the options for ParseHostname. A nil
// *ParseHostnameArgs, or a zero value, rejects Unicode host names.
type ParseHostnameArgs struct {
	// ASCIIToer, when set, converts a Unicode name to A-labels; the parsed
	// value holds the A-label form.
	ASCIIToer ASCIIToer
}

// ParseHostname parses s as an IP address or a host domain name. An IP
// address is stored as net/netip's canonical text. args may be nil; see
// ParseHostnameArgs. It fails with ErrInvalidHostname, joined with
// ErrUnsupportedSyntax for input some RFC allows but Hostname rejects; see
// Hostname.
func ParseHostname(s string, args *ParseHostnameArgs) (h Hostname, err error) {
	var name string
	var ip IPAddress

	if args == nil {
		args = &ParseHostnameArgs{}
	}
	if _, ipErr := netip.ParseAddr(s); ipErr == nil {
		ip, err = ParseIPAddress(s)
		if err != nil {
			err = NewErr(ErrInvalidHostname, "input", s, err)
			goto end
		}
		h = Hostname(ip)
		goto end
	}
	name, err = hostToASCII(s, args.ASCIIToer)
	if err != nil {
		err = NewErr(ErrInvalidHostname, "input", s, err)
		goto end
	}
	err = domainNameErr(name, hostDomainNameRules, ErrInvalidHostname, "dt.Hostname")
	if err != nil {
		goto end
	}
	h = Hostname(name)
end:
	return h, err
}

// NewHostname converts any member of the host family to a Hostname. It cannot
// fail. Its type parameter admits only the member types, so passing a plain
// string does not compile; use ParseHostname for a string.
func NewHostname[T IPAddress | IPv4Address | IPv6Address | HostDomainName | SingleLabelHostname | InternetDomain](h T) Hostname {
	return Hostname(h)
}

// IPAddress returns h as an IPAddress and true when it holds an IP address,
// or "" and false otherwise.
func (h Hostname) IPAddress() (ip IPAddress, ok bool) {
	ip, err := ParseIPAddress(string(h))
	if err != nil {
		goto end
	}
	ok = true
end:
	return ip, ok
}

// HostDomainName returns h as a HostDomainName and true when it holds a host
// domain name, or "" and false when it holds an IP address.
func (h Hostname) HostDomainName() (hdn HostDomainName, ok bool) {
	hdn, err := ParseHostDomainName(string(h))
	if err != nil {
		goto end
	}
	ok = true
end:
	return hdn, ok
}

// SingleLabelHostname returns h as a SingleLabelHostname and true when it
// holds a host domain name of one label, such as "localhost", or "" and false
// otherwise.
func (h Hostname) SingleLabelHostname() (slh SingleLabelHostname, ok bool) {
	slh, err := ParseSingleLabelHostname(string(h))
	if err != nil {
		goto end
	}
	ok = true
end:
	return slh, ok
}

// InternetDomain returns h as an InternetDomain and true when
// ParseInternetDomain with nil args accepts it, or "" and false otherwise. It
// checks syntax and the built-in special-use list only; call
// ParseInternetDomain with a PublicSuffixClassifier to check public suffixes.
func (h Hostname) InternetDomain() (d InternetDomain, ok bool) {
	d, err := ParseInternetDomain(string(h), nil)
	if err != nil {
		goto end
	}
	ok = true
end:
	return d, ok
}

// ToDNSName converts h to a DNSName. It fails with ErrInvalidDNSName and
// ErrIsIPAddress when h holds an IP address, which has no DNS name of its own.
func (h Hostname) ToDNSName() (dn DNSName, err error) {
	if _, ok := h.IPAddress(); ok {
		err = NewErr(ErrInvalidDNSName, ErrIsIPAddress, h.ErrKV())
		goto end
	}
	dn, err = ParseDNSName(string(h))
end:
	return dn, err
}

// ErrKV returns h as error metadata under the key "hostname".
func (h Hostname) ErrKV() ErrKV {
	return kv{k: "hostname", v: h}
}

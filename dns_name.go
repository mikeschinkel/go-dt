package dt

import (
	"strings"
)

// DNSName holds the name of a DNS record or zone, RFC 9499's "domain name",
// as x509.Certificate.DNSNames uses the term: "example.com",
// "_dmarc.example.com", "_sip._tcp.example.com.".
//
// DNSName adds two things to HostDomainName's syntax, both seen in practice:
// underscores in labels (RFC 8552's underscored attribute leaves), and a
// trailing root dot, which DNSName keeps, as zone files write it. The root
// name "." is accepted. Hostname and DNSName are not alternatives: they hold
// different sets, and ToHostname and Hostname.ToDNSName convert between them
// where they overlap.
//
// DNSName rejects:
//   - RFC 2181 section 11's arbitrary octets in a label. Writing them as text
//     needs RFC 1035 backslash escapes that no email address or URL accepts,
//     so a label may hold only ASCII letters, digits, hyphens and
//     underscores, and may not begin or end with a hyphen.
//   - A wildcard label such as "*.example.com".
//   - Unicode (RFC 5890 U-labels); write the name as A-labels.
//   - A final label made only of digits, an empty label, a label over 63
//     octets or a name over 253 octets, as HostDomainName does.
//
// If you need a type that accepts what DNSName rejects, define your own type,
// or open an issue at https://github.com/mikeschinkel/go-dt/issues to propose
// one.
type DNSName string

var dnsNameRules = domainNameRules{allowUnderscore: true, allowRootDot: true}

// ParseDNSName parses s as a DNS name. It fails with ErrInvalidDNSName,
// joined with ErrUnsupportedSyntax for a name some RFC allows but DNSName
// rejects; see DNSName.
func ParseDNSName(s string) (dn DNSName, err error) {
	err = domainNameErr(s, dnsNameRules, ErrInvalidDNSName, "dt.DNSName")
	if err != nil {
		goto end
	}
	dn = DNSName(s)
end:
	return dn, err
}

// ToHostname converts dn to a Hostname. It fails with ErrInvalidHostname when
// dn has no host name counterpart: a trailing root dot, or a label holding an
// underscore.
func (dn DNSName) ToHostname() (h Hostname, err error) {
	var hdn HostDomainName
	if strings.HasSuffix(string(dn), ".") {
		err = NewErr(ErrInvalidHostname, ErrRootDot, dn.ErrKV())
		goto end
	}
	hdn, err = ParseHostDomainName(string(dn))
	if err != nil {
		err = NewErr(ErrInvalidHostname, dn.ErrKV(), err)
		goto end
	}
	h = Hostname(hdn)
end:
	return h, err
}

// ErrKV returns dn as error metadata under the key "dns_name".
func (dn DNSName) ErrKV() ErrKV {
	return kv{k: "dns_name", v: dn}
}

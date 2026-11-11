package dt

import (
	"net/netip"
)

// IPAddress holds an IPv4 or IPv6 address as net/netip's canonical text.
//
// The name follows net/netip.Addr, which documents itself as "an IPv4 or IPv6
// address". IPAddress is a string type rather than a netip.Addr so that it
// converts to Hostname like every other member of the host family, and so
// that == compares two spellings of the same address as equal: the Parse
// functions store netip's canonical text, so "2001:DB8::0:1" is stored as
// "2001:db8::1" (RFC 5952).
//
// IPAddress rejects:
//   - An IPv6 zone index such as "fe80::1%en0" (RFC 4007, RFC 6874). A zone
//     names an interface on one machine, so an address carrying one does not
//     mean the same thing anywhere else. ScopedIPv6Address is declared for
//     that use.
//   - Brackets, as in "[::1]". Brackets belong to URL and email syntax, not
//     to the address.
//   - IPv4 with leading zeros, such as "192.168.01.1", which some parsers
//     read as octal.
//
// If you need a type that accepts what IPAddress rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type IPAddress string

// IPv4Address holds an IPv4 address in dotted-decimal form, as RFC 3986's
// IPv4address rule spells it: "192.0.2.1".
//
// IPv4Address rejects everything IPAddress rejects, and any IPv6 address,
// including an IPv4-mapped one such as "::ffff:192.0.2.1": that is IPv6 text,
// so it parses as IPv6Address.
//
// If you need a type that accepts what IPv4Address rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type IPv4Address string

// IPv6Address holds an IPv6 address (RFC 4291 "IPv6 address") in the
// canonical text form of RFC 5952, which net/netip produces: lower case, the
// longest run of zero groups compressed, and leading zeros dropped.
//
// IPv6Address accepts an IPv4-mapped address such as "::ffff:192.0.2.1". It
// rejects everything IPAddress rejects, and any IPv4 address.
//
// If you need a type that accepts what IPv6Address rejects, define your own
// type, or open an issue at https://github.com/mikeschinkel/go-dt/issues to
// propose one.
type IPv6Address string

// ParseIPAddress parses s as an IPv4 or IPv6 address and returns net/netip's
// canonical text for it. It fails with ErrInvalidIPAddress; see IPAddress for
// what it rejects.
func ParseIPAddress(s string) (ip IPAddress, err error) {
	var addr netip.Addr
	addr, err = parseIPAddr(s, ErrInvalidIPAddress, "dt.IPAddress")
	if err != nil {
		goto end
	}
	ip = IPAddress(addr.String())
end:
	return ip, err
}

// ParseIPv4Address parses s as an IPv4 address in dotted-decimal form. It
// fails with ErrInvalidIPv4Address; see IPv4Address for what it rejects.
func ParseIPv4Address(s string) (ip IPv4Address, err error) {
	var addr netip.Addr
	addr, err = parseIPAddr(s, ErrInvalidIPv4Address, "dt.IPv4Address")
	if err != nil {
		goto end
	}
	if !addr.Is4() {
		err = NewErr(ErrInvalidIPv4Address, ErrNotIPv4Address, "input", s)
		goto end
	}
	ip = IPv4Address(addr.String())
end:
	return ip, err
}

// ParseIPv6Address parses s as an IPv6 address and returns its RFC 5952
// canonical text: ParseIPv6Address("2001:DB8::0:1") returns "2001:db8::1".
// It fails with ErrInvalidIPv6Address; see IPv6Address for what it rejects.
func ParseIPv6Address(s string) (ip IPv6Address, err error) {
	var addr netip.Addr
	addr, err = parseIPAddr(s, ErrInvalidIPv6Address, "dt.IPv6Address")
	if err != nil {
		goto end
	}
	if !addr.Is6() {
		err = NewErr(ErrInvalidIPv6Address, ErrNotIPv6Address, "input", s)
		goto end
	}
	ip = IPv6Address(addr.String())
end:
	return ip, err
}

// NewIPAddress converts an IPv4Address or IPv6Address to an IPAddress. It
// cannot fail. Its type parameter admits only the member types, so passing a
// plain string does not compile; use ParseIPAddress for a string.
func NewIPAddress[T IPv4Address | IPv6Address](ip T) IPAddress {
	return IPAddress(ip)
}

// IPv4Address returns ip as an IPv4Address and true when it holds an IPv4
// address, or "" and false otherwise.
func (ip IPAddress) IPv4Address() (v4 IPv4Address, ok bool) {
	addr, err := netip.ParseAddr(string(ip))
	if err != nil || !addr.Is4() {
		goto end
	}
	v4, ok = IPv4Address(ip), true
end:
	return v4, ok
}

// IPv6Address returns ip as an IPv6Address and true when it holds an IPv6
// address, or "" and false otherwise.
func (ip IPAddress) IPv6Address() (v6 IPv6Address, ok bool) {
	addr, err := netip.ParseAddr(string(ip))
	if err != nil || !addr.Is6() || addr.Zone() != "" {
		goto end
	}
	v6, ok = IPv6Address(ip), true
end:
	return v6, ok
}

// ErrKV returns ip as error metadata under the key "ip_address".
func (ip IPAddress) ErrKV() ErrKV {
	return kv{k: "ip_address", v: ip}
}

// ErrKV returns ip as error metadata under the key "ipv4_address".
func (ip IPv4Address) ErrKV() ErrKV {
	return kv{k: "ipv4_address", v: ip}
}

// ErrKV returns ip as error metadata under the key "ipv6_address".
func (ip IPv6Address) ErrKV() ErrKV {
	return kv{k: "ipv6_address", v: ip}
}

// parseIPAddr parses s with net/netip and rejects a zone, wrapping any
// failure in sentinel. typeName names the type whose doc comment explains an
// unsupported rejection.
func parseIPAddr(s string, sentinel error, typeName string) (addr netip.Addr, err error) {
	addr, err = netip.ParseAddr(s)
	if err != nil {
		err = NewErr(sentinel, "input", s, err)
		goto end
	}
	if addr.Zone() != "" {
		err = NewErr(sentinel, ErrUnsupportedSyntax, ErrIPv6Zone, "type", typeName, "input", s)
		addr = netip.Addr{}
	}
end:
	return addr, err
}

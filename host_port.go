package dt

// The types in this file are declared but not yet implemented. Each Parse
// function returns ErrNotYetSupported. They are declared so they show up in
// godoc and editors, and each doc comment records the options weighed and the
// use cases that would decide them. If you have such a use case, open an issue
// at https://github.com/mikeschinkel/go-dt/issues describing it.

// ScopedIPv6Address will hold an IPv6 address with a zone index, RFC 4007's
// "zone index" written as RFC 6874's ZoneID: "fe80::1%en0". IPAddress and
// IPv6Address reject zones because a zone names an interface on one machine.
//
// Not yet supported. The open questions:
//   - Whether a scoped address belongs inside IPAddress and Hostname, so a
//     Hostname could hold "fe80::1%en0", or only inside HostPort, where a
//     dial target is machine-local anyway.
//   - Whether to store the zone as given or normalize it (an interface name
//     versus its index), which only the machine that owns it can do.
//
// A use case that dials or listens on a link-local address would decide them.
type ScopedIPv6Address string

// ParseScopedIPv6Address returns ErrNotYetSupported; see ScopedIPv6Address.
func ParseScopedIPv6Address(s string) (ip ScopedIPv6Address, err error) {
	err = NewErr(ErrInvalidScopedIPv6Address, ErrNotYetSupported, "input", s)
	return ip, err
}

// ErrKV returns ip as error metadata under the key "scoped_ipv6_address".
func (ip ScopedIPv6Address) ErrKV() ErrKV {
	return kv{k: "scoped_ipv6_address", v: ip}
}

// HostPort will hold a host and a port joined as net.SplitHostPort and
// net.JoinHostPort spell it, their "hostport": "example.com:443",
// "[2001:db8::1]:443", "example.com:https".
//
// Not yet supported. The open questions:
//   - Whether the port may be a ServiceName ("https"), as net.Dial accepts,
//     or only a PortNumber. HostPort declares both a generic Port method and
//     the typed PortNumber and ServiceName methods so either answer fits.
//   - How a HostPort holding a ScopedIPv6Address answers Hostname(), given
//     that Hostname rejects zones.
//   - Whether a missing port is an error or means a scheme's default.
//
// A use case that stores dial targets or listen addresses would decide them.
type HostPort string

// ParseHostPort returns ErrNotYetSupported; see HostPort.
func ParseHostPort(s string) (hp HostPort, err error) {
	err = NewErr(ErrInvalidHostPort, ErrNotYetSupported, "input", s)
	return hp, err
}

// Port will return hp's port as a PortNumber or a ServiceName, chosen by the
// type argument. It returns ErrNotYetSupported; see HostPort.
func (hp HostPort) Port[T PortNumber | ServiceName]() (port T, err error) {
	err = NewErr(ErrInvalidHostPort, ErrNotYetSupported, hp.ErrKV())
	return port, err
}

// PortNumber will return hp's port as a PortNumber. It returns
// ErrNotYetSupported; see HostPort.
func (hp HostPort) PortNumber() (PortNumber, error) {
	return hp.Port[PortNumber]()
}

// ServiceName will return hp's port as a ServiceName. It returns
// ErrNotYetSupported; see HostPort.
func (hp HostPort) ServiceName() (ServiceName, error) {
	return hp.Port[ServiceName]()
}

// ErrKV returns hp as error metadata under the key "host_port".
func (hp HostPort) ErrKV() ErrKV {
	return kv{k: "host_port", v: hp}
}

// PortNumber will hold a TCP or UDP port number, RFC 6335's "port number",
// 0 through 65535.
//
// Not yet supported. The open questions:
//   - Whether 0, which asks the system to choose a port when listening but is
//     reserved on the wire, is valid.
//   - Whether to distinguish RFC 6335's system, user and dynamic ranges.
//
// A use case that stores configured ports would decide them.
type PortNumber uint16

// ParsePortNumber returns ErrNotYetSupported; see PortNumber.
func ParsePortNumber(s string) (pn PortNumber, err error) {
	err = NewErr(ErrInvalidPortNumber, ErrNotYetSupported, "input", s)
	return pn, err
}

// ErrKV returns pn as error metadata under the key "port_number".
func (pn PortNumber) ErrKV() ErrKV {
	return kv{k: "port_number", v: pn}
}

// ServiceName will hold an IANA service name, RFC 6335's "service name":
// "https", "ssh", "xmpp-client".
//
// Not yet supported. The open questions:
//   - Whether to enforce RFC 6335 section 5.1's syntax only (1 to 15
//     characters, letters, digits and interior hyphens) or also require the
//     name be registered, which needs the IANA registry or /etc/services.
//   - Whether a ServiceName converts to a PortNumber, which needs the same
//     lookup.
//
// A use case that accepts ports by name would decide them.
type ServiceName string

// ParseServiceName returns ErrNotYetSupported; see ServiceName.
func ParseServiceName(s string) (sn ServiceName, err error) {
	err = NewErr(ErrInvalidServiceName, ErrNotYetSupported, "input", s)
	return sn, err
}

// ErrKV returns sn as error metadata under the key "service_name".
func (sn ServiceName) ErrKV() ErrKV {
	return kv{k: "service_name", v: sn}
}

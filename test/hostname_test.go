package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParseHostDomainName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		unsupported bool
		cause       error
	}{
		{name: "two labels", input: "example.com"},
		{name: "local suffix", input: "mikes-mbp.local"},
		{name: "test suffix", input: "repo.test"},
		{name: "single label", input: "localhost"},
		{name: "digit-first label (RFC 1123)", input: "3com.com"},
		{name: "all-digit inner label", input: "1.example.com"},
		{name: "A-label", input: "xn--bcher-kva.example"},
		{name: "case kept", input: "Example.COM"},
		{name: "63-octet label", input: strings.Repeat("a", 63) + ".com"},
		{name: "hyphen edges", input: "-bad-.com", wantErr: true, unsupported: true, cause: dt.ErrNonLDHLabel},
		{name: "underscore", input: "_dmarc.example.com", wantErr: true, unsupported: true, cause: dt.ErrNonLDHLabel},
		{name: "root dot", input: "example.com.", wantErr: true, unsupported: true, cause: dt.ErrRootDot},
		{name: "unicode", input: "bücher.example", wantErr: true, unsupported: true, cause: dt.ErrUnicodeHostname},
		{name: "numeric final label", input: "1.2.3.4", wantErr: true, cause: dt.ErrNumericFinalLabel},
		{name: "empty", input: "", wantErr: true, cause: dt.ErrEmpty},
		{name: "empty label", input: "a..com", wantErr: true, cause: dt.ErrEmptyLabel},
		{name: "leading dot", input: ".com", wantErr: true, cause: dt.ErrEmptyLabel},
		{name: "64-octet label", input: strings.Repeat("a", 64) + ".com", wantErr: true, cause: dt.ErrLabelTooLong},
		{name: "254-octet name", input: strings.Repeat("a.", 126) + "ab", wantErr: true, cause: dt.ErrTooLong},
		{name: "space", input: "a b.com", wantErr: true, unsupported: true, cause: dt.ErrNonLDHLabel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseHostDomainName(tt.input)
			want := dt.HostDomainName(tt.input)
			checkParse(t, got, err, want, tt.wantErr, dt.ErrInvalidHostDomainName, tt.unsupported)
			checkCause(t, err, tt.cause)
		})
	}
}

func TestParseSingleLabelHostname(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "localhost", input: "localhost"},
		{name: "buildbox", input: "buildbox"},
		{name: "two labels", input: "mikes-mbp.local", wantErr: true},
		{name: "numeric", input: "1234", wantErr: true},
		{name: "hyphen edge", input: "bad-", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseSingleLabelHostname(tt.input)
			if tt.wantErr {
				if !errors.Is(err, dt.ErrInvalidSingleLabelHostname) {
					t.Fatalf("got %q, %v; want ErrInvalidSingleLabelHostname", got, err)
				}
				return
			}
			if err != nil || got != dt.SingleLabelHostname(tt.input) {
				t.Fatalf("got %q, %v; want %q", got, err, tt.input)
			}
		})
	}
}

func TestParseHostname(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		args        *dt.ParseHostnameArgs
		want        dt.Hostname
		wantErr     bool
		unsupported bool
	}{
		{name: "IPv4", input: "192.0.2.1", want: "192.0.2.1"},
		{name: "IPv6 canonicalised", input: "2001:DB8::0:1", want: "2001:db8::1"},
		{name: "internet domain", input: "example.com", want: "example.com"},
		{name: "single label", input: "localhost", want: "localhost"},
		{name: "special-use", input: "repo.test", want: "repo.test"},
		{name: "bracketed IPv6 rejected", input: "[::1]", wantErr: true, unsupported: true},
		{name: "zone rejected", input: "fe80::1%en0", wantErr: true, unsupported: true},
		{name: "unicode without IDNAConverter", input: "bücher.example", wantErr: true, unsupported: true},
		{name: "unicode with IDNAConverter", input: "bücher.example", args: &dt.ParseHostnameArgs{IDNAConverter: fakeIDNA}, want: "xn--bcher-kva.example"},
		{name: "IDNAConverter error", input: "badé.example", args: &dt.ParseHostnameArgs{IDNAConverter: fakeIDNA}, wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseHostname(tt.input, tt.args)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidHostname, tt.unsupported)
		})
	}
}

func TestHostnameAccessors(t *testing.T) {
	tests := []struct {
		host     dt.Hostname
		ip       dt.IPAddress
		hdn      dt.HostDomainName
		single   dt.SingleLabelHostname
		internet dt.InternetDomain
	}{
		{host: "192.0.2.1", ip: "192.0.2.1"},
		{host: "2001:db8::1", ip: "2001:db8::1"},
		{host: "localhost", hdn: "localhost", single: "localhost"},
		{host: "example.com", hdn: "example.com", internet: "example.com"},
		{host: "mikes-mbp.local", hdn: "mikes-mbp.local"},
	}
	for _, tt := range tests {
		t.Run(string(tt.host), func(t *testing.T) {
			ip, ok := tt.host.IPAddress()
			if ip != tt.ip || ok != (tt.ip != "") {
				t.Errorf("IPAddress() = %q, %v", ip, ok)
			}
			hdn, ok := tt.host.HostDomainName()
			if hdn != tt.hdn || ok != (tt.hdn != "") {
				t.Errorf("HostDomainName() = %q, %v", hdn, ok)
			}
			single, ok := tt.host.SingleLabelHostname()
			if single != tt.single || ok != (tt.single != "") {
				t.Errorf("SingleLabelHostname() = %q, %v", single, ok)
			}
			internet, ok := tt.host.InternetDomain()
			if internet != tt.internet || ok != (tt.internet != "") {
				t.Errorf("InternetDomain() = %q, %v", internet, ok)
			}
		})
	}
}

func TestNewHostname(t *testing.T) {
	tests := []struct {
		name string
		got  dt.Hostname
		want dt.Hostname
	}{
		{name: "IPAddress", got: dt.NewHostname(dt.IPAddress("192.0.2.1")), want: "192.0.2.1"},
		{name: "IPv4Address", got: dt.NewHostname(dt.IPv4Address("192.0.2.1")), want: "192.0.2.1"},
		{name: "IPv6Address", got: dt.NewHostname(dt.IPv6Address("::1")), want: "::1"},
		{name: "HostDomainName", got: dt.NewHostname(dt.HostDomainName("a.local")), want: "a.local"},
		{name: "SingleLabelHostname", got: dt.NewHostname(dt.SingleLabelHostname("localhost")), want: "localhost"},
		{name: "InternetDomain", got: dt.NewHostname(dt.InternetDomain("example.com")), want: "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("NewHostname() = %q; want %q", tt.got, tt.want)
			}
		})
	}
}

func TestConvertUpFamily(t *testing.T) {
	// Converting up the family is a plain type conversion.
	d := dt.InternetDomain("example.com")
	h := dt.Hostname(dt.HostDomainName(d))
	if h != "example.com" {
		t.Errorf("conversion = %q", h)
	}
}

func TestHostnameToDNSName(t *testing.T) {
	dn, err := dt.Hostname("example.com").ToDNSName()
	if err != nil || dn != "example.com" {
		t.Errorf("ToDNSName() = %q, %v; want example.com", dn, err)
	}
	dn, err = dt.Hostname("192.0.2.1").ToDNSName()
	if !errors.Is(err, dt.ErrInvalidDNSName) || !errors.Is(err, dt.ErrIsIPAddress) || dn != "" {
		t.Errorf("ToDNSName() on IP = %q, %v; want ErrIsIPAddress", dn, err)
	}
}

func TestParseDNSName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		unsupported bool
	}{
		{name: "plain", input: "example.com"},
		{name: "underscored label (RFC 8552)", input: "_dmarc.example.com"},
		{name: "SRV name", input: "_sip._tcp.example.com"},
		{name: "root dot kept", input: "example.com."},
		{name: "root", input: "."},
		{name: "single label", input: "localhost"},
		{name: "wildcard", input: "*.example.com", wantErr: true, unsupported: true},
		{name: "hyphen edge", input: "-bad-.example.com", wantErr: true, unsupported: true},
		{name: "unicode", input: "bücher.example", wantErr: true, unsupported: true},
		{name: "double dot", input: "a..com", wantErr: true},
		{name: "two root dots", input: "example.com..", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseDNSName(tt.input)
			checkParse(t, got, err, dt.DNSName(tt.input), tt.wantErr, dt.ErrInvalidDNSName, tt.unsupported)
		})
	}
}

func TestDNSNameToHostname(t *testing.T) {
	h, err := dt.DNSName("example.com").ToHostname()
	if err != nil || h != "example.com" {
		t.Errorf("ToHostname() = %q, %v; want example.com", h, err)
	}
	for _, dn := range []dt.DNSName{"_dmarc.example.com", "example.com.", "."} {
		h, err = dn.ToHostname()
		if !errors.Is(err, dt.ErrInvalidHostname) || h != "" {
			t.Errorf("%q.ToHostname() = %q, %v; want ErrInvalidHostname", dn, h, err)
		}
	}
}

func TestHostErrKV(t *testing.T) {
	checkErrKV(t, dt.HostDomainName("a.b").ErrKV(), "host_domain_name")
	checkErrKV(t, dt.SingleLabelHostname("a").ErrKV(), "single_label_hostname")
	checkErrKV(t, dt.InternetDomain("a.com").ErrKV(), "internet_domain")
	checkErrKV(t, dt.Hostname("a").ErrKV(), "hostname")
	checkErrKV(t, dt.DNSName("a").ErrKV(), "dns_name")
}

// fakeIDNA stands in for golang.org/x/net/idna, which go-dt's test module
// does not require. It knows one U-label and fails on anything else
// non-ASCII.
var fakeIDNA = dt.ToASCIIHostDomainNameFunc(func(name string) (dt.HostDomainName, error) {
	out := strings.ReplaceAll(name, "bücher", "xn--bcher-kva")
	for _, r := range out {
		if r >= 0x80 {
			return "", errors.New("fake idna: unknown U-label")
		}
	}
	return dt.HostDomainName(out), nil
})

// checkCause asserts err matches cause, when cause is set.
func checkCause(t *testing.T, err, cause error) {
	t.Helper()
	if cause != nil && !errors.Is(err, cause) {
		t.Errorf("error %v does not match %v", err, cause)
	}
}

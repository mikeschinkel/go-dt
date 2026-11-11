package test

import (
	"errors"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParseIPAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        dt.IPAddress
		wantErr     bool
		unsupported bool
	}{
		{name: "IPv4", input: "192.0.2.1", want: "192.0.2.1"},
		{name: "IPv6 canonicalised", input: "2001:DB8::0:1", want: "2001:db8::1"},
		{name: "IPv4-mapped IPv6", input: "::ffff:192.0.2.1", want: "::ffff:192.0.2.1"},
		{name: "zone rejected", input: "fe80::1%en0", wantErr: true, unsupported: true},
		{name: "brackets rejected", input: "[::1]", wantErr: true},
		{name: "leading zero rejected", input: "192.168.01.1", wantErr: true},
		{name: "empty rejected", input: "", wantErr: true},
		{name: "name rejected", input: "example.com", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseIPAddress(tt.input)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidIPAddress, tt.unsupported)
		})
	}
}

func TestParseIPv4Address(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        dt.IPv4Address
		wantErr     bool
		unsupported bool
	}{
		{name: "IPv4", input: "192.0.2.1", want: "192.0.2.1"},
		{name: "IPv6 rejected", input: "2001:db8::1", wantErr: true},
		{name: "IPv4-mapped rejected", input: "::ffff:192.0.2.1", wantErr: true},
		{name: "zone rejected", input: "fe80::1%en0", wantErr: true, unsupported: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseIPv4Address(tt.input)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidIPv4Address, tt.unsupported)
		})
	}
}

func TestParseIPv6Address(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        dt.IPv6Address
		wantErr     bool
		unsupported bool
	}{
		{name: "canonical text", input: "2001:DB8::0:1", want: "2001:db8::1"},
		{name: "RFC 5952 compression", input: "2001:db8:0:0:1:0:0:1", want: "2001:db8::1:0:0:1"},
		{name: "IPv4-mapped accepted", input: "::ffff:192.0.2.1", want: "::ffff:192.0.2.1"},
		{name: "IPv4 rejected", input: "192.0.2.1", wantErr: true},
		{name: "zone rejected", input: "fe80::1%en0", wantErr: true, unsupported: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseIPv6Address(tt.input)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidIPv6Address, tt.unsupported)
		})
	}
}

func TestIPAddressAccessors(t *testing.T) {
	v4, ok := dt.IPAddress("192.0.2.1").IPv4Address()
	if !ok || v4 != "192.0.2.1" {
		t.Errorf("IPv4Address() = %q, %v; want 192.0.2.1, true", v4, ok)
	}
	if _, ok = dt.IPAddress("192.0.2.1").IPv6Address(); ok {
		t.Errorf("IPv6Address() on IPv4 = true; want false")
	}
	v6, ok := dt.IPAddress("2001:db8::1").IPv6Address()
	if !ok || v6 != "2001:db8::1" {
		t.Errorf("IPv6Address() = %q, %v; want 2001:db8::1, true", v6, ok)
	}
	if _, ok = dt.IPAddress("2001:db8::1").IPv4Address(); ok {
		t.Errorf("IPv4Address() on IPv6 = true; want false")
	}
	if _, ok = dt.IPAddress("::ffff:192.0.2.1").IPv4Address(); ok {
		t.Errorf("IPv4Address() on IPv4-mapped IPv6 = true; want false")
	}
}

func TestNewIPAddress(t *testing.T) {
	if got := dt.NewIPAddress(dt.IPv4Address("192.0.2.1")); got != "192.0.2.1" {
		t.Errorf("NewIPAddress(IPv4Address) = %q", got)
	}
	if got := dt.NewIPAddress(dt.IPv6Address("2001:db8::1")); got != "2001:db8::1" {
		t.Errorf("NewIPAddress(IPv6Address) = %q", got)
	}
}

func TestIPErrKV(t *testing.T) {
	checkErrKV(t, dt.IPAddress("192.0.2.1").ErrKV(), "ip_address")
	checkErrKV(t, dt.IPv4Address("192.0.2.1").ErrKV(), "ipv4_address")
	checkErrKV(t, dt.IPv6Address("::1").ErrKV(), "ipv6_address")
}

// checkParse asserts a Parse function's result: the value on success, or an
// error matching sentinel (and ErrUnsupportedSyntax when unsupported) on
// failure.
func checkParse[T comparable](t *testing.T, got T, err error, want T, wantErr bool, sentinel error, unsupported bool) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Fatalf("got %v, nil; want error", got)
		}
		if !errors.Is(err, sentinel) {
			t.Errorf("error %q does not match %q", err, sentinel)
		}
		if unsupported != errors.Is(err, dt.ErrUnsupportedSyntax) {
			t.Errorf("errors.Is(err, ErrUnsupportedSyntax) = %v; want %v (err: %v)",
				!unsupported, unsupported, err)
		}
		var zero T
		if got != zero {
			t.Errorf("got %v with error; want zero value", got)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

// checkErrKV asserts an ErrKV method's key.
func checkErrKV(t *testing.T, kv dt.ErrKV, key string) {
	t.Helper()
	if kv.Key() != key {
		t.Errorf("ErrKV key = %q; want %q", kv.Key(), key)
	}
}

package test

import (
	"errors"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestDeferredTypesNotYetSupported(t *testing.T) {
	hp := dt.HostPort("example.com:443")
	tests := []struct {
		name     string
		err      error
		sentinel error
	}{
		{name: "ParseScopedIPv6Address", err: second(dt.ParseScopedIPv6Address("fe80::1%en0")), sentinel: dt.ErrInvalidScopedIPv6Address},
		{name: "ParseHostPort", err: second(dt.ParseHostPort("example.com:443")), sentinel: dt.ErrInvalidHostPort},
		{name: "ParsePortNumber", err: second(dt.ParsePortNumber("443")), sentinel: dt.ErrInvalidPortNumber},
		{name: "ParseServiceName", err: second(dt.ParseServiceName("https")), sentinel: dt.ErrInvalidServiceName},
		{name: "HostPort.Port[PortNumber]", err: second(hp.Port[dt.PortNumber]()), sentinel: dt.ErrInvalidHostPort},
		{name: "HostPort.Port[ServiceName]", err: second(hp.Port[dt.ServiceName]()), sentinel: dt.ErrInvalidHostPort},
		{name: "HostPort.PortNumber", err: second(hp.PortNumber()), sentinel: dt.ErrInvalidHostPort},
		{name: "HostPort.ServiceName", err: second(hp.ServiceName()), sentinel: dt.ErrInvalidHostPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.err, dt.ErrNotYetSupported) || !errors.Is(tt.err, tt.sentinel) {
				t.Errorf("error = %v; want ErrNotYetSupported and %v", tt.err, tt.sentinel)
			}
		})
	}
}

func TestDeferredErrKV(t *testing.T) {
	checkErrKV(t, dt.ScopedIPv6Address("fe80::1%en0").ErrKV(), "scoped_ipv6_address")
	checkErrKV(t, dt.HostPort("a:1").ErrKV(), "host_port")
	checkErrKV(t, dt.PortNumber(443).ErrKV(), "port_number")
	checkErrKV(t, dt.ServiceName("https").ErrKV(), "service_name")
}

// second returns the error of a (value, error) pair.
func second[T any](_ T, err error) error {
	return err
}

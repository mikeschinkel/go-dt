package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParseEmailLocalPart(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "dot-atom", input: "john.doe"},
		{name: "atext specials", input: "a!#$%&'*+-/=?^_`{|}~z"},
		{name: "quoted with space", input: `"john doe"`},
		{name: "quoted with quoted-pair", input: `"a\"b"`},
		{name: "quoted with at", input: `"a@b"`},
		{name: "empty quoted", input: `""`},
		{name: "UTF-8 accepted", input: "jöhn"},
		{name: "64 octets", input: strings.Repeat("a", 64)},
		{name: "65 octets", input: strings.Repeat("a", 65), wantErr: true},
		{name: "empty", input: "", wantErr: true},
		{name: "double dot (obsolete)", input: "a..b", wantErr: true},
		{name: "leading dot", input: ".a", wantErr: true},
		{name: "trailing dot", input: "a.", wantErr: true},
		{name: "space", input: "john doe", wantErr: true},
		{name: "comment", input: "john(comment)", wantErr: true},
		{name: "unterminated quote", input: `"john`, wantErr: true},
		{name: "unescaped quote inside", input: `"a"b"`, wantErr: true},
		{name: "trailing backslash", input: `"a\"`, wantErr: true},
		{name: "control inside quotes", input: "\"a\tb\"", wantErr: true},
		{name: "invalid UTF-8", input: "a\xffb", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseEmailLocalPart(tt.input)
			checkParse(t, got, err, dt.EmailLocalPart(tt.input), tt.wantErr, dt.ErrInvalidEmailLocalPart, false)
		})
	}
}

func TestParseEmailAddress(t *testing.T) {
	utf8OK := &dt.ParseEmailAddressArgs{AllowSMTPUTF8: true}
	tests := []struct {
		name        string
		input       string
		args        *dt.ParseEmailAddressArgs
		want        dt.EmailAddress
		wantErr     bool
		unsupported bool
		cause       error
	}{
		{name: "plain", input: "john@example.com", want: "john@example.com"},
		{name: "quoted round-trips", input: `"john doe"@example.com`, want: `"john doe"@example.com`},
		{name: "quoted with at", input: `"a@b"@example.com`, want: `"a@b"@example.com`},
		{name: "case kept", input: "John@Example.COM", want: "John@Example.COM"},
		{name: "single-label host", input: "mike@localhost", want: "mike@localhost"},
		{name: "local host", input: "mike@mikes-mbp.local", want: "mike@mikes-mbp.local"},
		{name: "IPv4 literal", input: "user@[192.0.2.1]", want: "user@[192.0.2.1]"},
		{name: "IPv6 literal", input: "user@[IPv6:2001:db8::1]", want: "user@[IPv6:2001:db8::1]"},
		{name: "IPv6 literal canonicalised", input: "user@[ipv6:2001:DB8::0:1]", want: "user@[IPv6:2001:db8::1]"},
		{name: "bare IPv4 bracketed", input: "user@192.0.2.1", want: "user@[192.0.2.1]"},
		{name: "bare IPv6 bracketed", input: "user@2001:db8::1", want: "user@[IPv6:2001:db8::1]"},
		{name: "general address literal", input: "user@[x400:abc]", wantErr: true, unsupported: true, cause: dt.ErrGeneralAddressLiteral},
		{name: "IPv6 literal without tag", input: "user@[2001:db8::1]", wantErr: true, unsupported: true, cause: dt.ErrGeneralAddressLiteral},
		{name: "IPv4 in IPv6 tag", input: "user@[IPv6:192.0.2.1]", wantErr: true, cause: dt.ErrInvalidIPv6Address},
		{name: "unclosed literal", input: "user@[192.0.2.1", wantErr: true},
		{name: "zone", input: "user@fe80::1%en0", wantErr: true, unsupported: true, cause: dt.ErrIPv6Zone},
		{name: "bad label", input: "user@-bad-.com", wantErr: true, unsupported: true, cause: dt.ErrNonLDHLabel},
		{name: "no at", input: "user.example.com", wantErr: true, cause: dt.ErrMissingAtSign},
		{name: "empty host", input: "user@", wantErr: true},
		{name: "empty local part", input: "@example.com", wantErr: true, cause: dt.ErrInvalidEmailLocalPart},
		{name: "angle brackets", input: "<user@example.com>", wantErr: true},
		{name: "UTF-8 local part refused", input: "jöhn@example.com", wantErr: true, cause: dt.ErrSMTPUTF8NotAllowed},
		{name: "UTF-8 local part allowed", input: "jöhn@example.com", args: utf8OK, want: "jöhn@example.com"},
		{name: "unicode host without IDNAConverter", input: "john@bücher.example", args: utf8OK, wantErr: true, unsupported: true, cause: dt.ErrUnicodeHostname},
		{name: "unicode host with IDNAConverter", input: "john@bücher.example", args: &dt.ParseEmailAddressArgs{IDNAConverter: fakeIDNA}, want: "john@xn--bcher-kva.example"},
		{name: "IDNAConverter error", input: "john@badé.example", args: &dt.ParseEmailAddressArgs{IDNAConverter: fakeIDNA}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseEmailAddress(tt.input, tt.args)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidEmailAddress, tt.unsupported)
			checkCause(t, err, tt.cause)
		})
	}
}

func TestEmailAddressParts(t *testing.T) {
	tests := []struct {
		email     dt.EmailAddress
		localPart dt.EmailLocalPart
		hostname  dt.Hostname
	}{
		{email: "john@example.com", localPart: "john", hostname: "example.com"},
		{email: `"a@b"@example.com`, localPart: `"a@b"`, hostname: "example.com"},
		{email: "user@[192.0.2.1]", localPart: "user", hostname: "192.0.2.1"},
		{email: "user@[IPv6:2001:db8::1]", localPart: "user", hostname: "2001:db8::1"},
	}
	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			if got := tt.email.LocalPart(); got != tt.localPart {
				t.Errorf("LocalPart() = %q; want %q", got, tt.localPart)
			}
			if got := tt.email.Hostname(); got != tt.hostname {
				t.Errorf("Hostname() = %q; want %q", got, tt.hostname)
			}
		})
	}
}

func TestEmailAddressEquals(t *testing.T) {
	tests := []struct {
		a, b dt.EmailAddress
		want bool
	}{
		{a: "john@example.com", b: "john@example.com", want: true},
		{a: "john@example.com", b: "john@EXAMPLE.com", want: true},
		{a: "john@example.com", b: "John@example.com", want: false},
		{a: "john@example.com", b: "john@example.org", want: false},
		{a: "user@[IPv6:2001:db8::1]", b: "user@[ipv6:2001:db8::1]", want: true},
	}
	for _, tt := range tests {
		t.Run(string(tt.a)+"="+string(tt.b), func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.want {
				t.Errorf("Equals() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAddressRequiresSMTPUTF8(t *testing.T) {
	if dt.EmailAddress("john@example.com").RequiresSMTPUTF8() {
		t.Error("ASCII address requires SMTPUTF8")
	}
	if !dt.EmailAddress("jöhn@example.com").RequiresSMTPUTF8() {
		t.Error("UTF-8 address does not require SMTPUTF8")
	}
}

func TestEmailAddressValidate(t *testing.T) {
	for _, ea := range []dt.EmailAddress{"john@example.com", "jöhn@example.com", "user@[192.0.2.1]"} {
		if err := ea.Validate(); err != nil {
			t.Errorf("%q.Validate() = %v; want nil", ea, err)
		}
	}
	ea := dt.EmailAddress("john@-bad-.com")
	_, parseErr := dt.ParseEmailAddress(string(ea), nil)
	err := ea.Validate()
	if err == nil || err.Error() != parseErr.Error() {
		t.Errorf("Validate() = %v; want the Parse error %v", err, parseErr)
	}
	for _, sentinel := range []error{dt.ErrInvalidEmailAddress, dt.ErrUnsupportedSyntax, dt.ErrNonLDHLabel} {
		if !errors.Is(err, sentinel) {
			t.Errorf("Validate() error does not match %v", sentinel)
		}
	}
	if err = dt.EmailAddress("").Validate(); !errors.Is(err, dt.ErrInvalidEmailAddress) {
		t.Errorf("empty Validate() = %v; want ErrInvalidEmailAddress", err)
	}
}

func TestEmailErrKV(t *testing.T) {
	checkErrKV(t, dt.EmailLocalPart("a").ErrKV(), "email_local_part")
	checkErrKV(t, dt.EmailAddress("a@b").ErrKV(), "email_address")
	checkErrKV(t, dt.EmailDisplayName("A").ErrKV(), "email_display_name")
	checkErrKV(t, dt.Mailbox{Email: "a@b"}.ErrKV(), "mailbox")
}

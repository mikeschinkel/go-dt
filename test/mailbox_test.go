package test

import (
	"errors"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TestParseEmailDisplayName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		unsupported bool
	}{
		{name: "plain", input: "John Doe"},
		{name: "UTF-8", input: "Jöhn Döe"},
		{name: "quotes kept", input: `"Doe, John"`},
		{name: "less than", input: "John <Doe", wantErr: true, unsupported: true},
		{name: "greater than", input: "John> Doe", wantErr: true, unsupported: true},
		{name: "carriage return", input: "John\rDoe", wantErr: true},
		{name: "line feed", input: "John\nDoe", wantErr: true},
		{name: "tab", input: "John\tDoe", wantErr: true},
		{name: "invalid UTF-8", input: "John\xff", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseEmailDisplayName(tt.input)
			checkParse(t, got, err, dt.EmailDisplayName(tt.input), tt.wantErr, dt.ErrInvalidEmailDisplayName, tt.unsupported)
		})
	}
}

func TestParseMailbox(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		args    *dt.ParseMailboxArgs
		want    dt.Mailbox
		wantErr bool
		cause   error
	}{
		{name: "name and address", input: "John Doe <john@example.com>", want: dt.Mailbox{DisplayName: "John Doe", Email: "john@example.com"}},
		{name: "spaces trimmed", input: "  John Doe   <john@example.com> ", want: dt.Mailbox{DisplayName: "John Doe", Email: "john@example.com"}},
		{name: "angle address only", input: "<john@example.com>", want: dt.Mailbox{Email: "john@example.com"}},
		{name: "bare address", input: "john@example.com", want: dt.Mailbox{Email: "john@example.com"}},
		{name: "git default identity", input: "mike <mike@localhost>", want: dt.Mailbox{DisplayName: "mike", Email: "mike@localhost"}},
		{name: "bare IP bracketed", input: "A <a@192.0.2.1>", want: dt.Mailbox{DisplayName: "A", Email: "a@[192.0.2.1]"}},
		{name: "UTF-8 refused", input: "J <jöhn@example.com>", wantErr: true, cause: dt.ErrSMTPUTF8NotAllowed},
		{name: "UTF-8 allowed", input: "J <jöhn@example.com>", args: &dt.ParseMailboxArgs{AllowSMTPUTF8: true}, want: dt.Mailbox{DisplayName: "J", Email: "jöhn@example.com"}},
		{name: "unicode host with ASCIIToer", input: "J <j@bücher.example>", args: &dt.ParseMailboxArgs{ASCIIToer: fakeToASCII}, want: dt.Mailbox{DisplayName: "J", Email: "j@xn--bcher-kva.example"}},
		{name: "unicode host without ASCIIToer", input: "J <j@bücher.example>", wantErr: true, cause: dt.ErrUnicodeHostname},
		{name: "bad address", input: "John <john@-bad-.com>", wantErr: true, cause: dt.ErrInvalidEmailAddress},
		{name: "bad name", input: "Jo\nhn <john@example.com>", wantErr: true, cause: dt.ErrInvalidEmailDisplayName},
		{name: "unopened bracket", input: "john@example.com>", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseMailbox(tt.input, tt.args)
			if tt.wantErr {
				if !errors.Is(err, dt.ErrInvalidMailbox) {
					t.Fatalf("got %+v, %v; want ErrInvalidMailbox", got, err)
				}
				checkCause(t, err, tt.cause)
				if got != (dt.Mailbox{}) {
					t.Errorf("got %+v with error; want zero value", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v; want %+v", got, tt.want)
			}
		})
	}
}

func TestMailboxString(t *testing.T) {
	tests := []struct {
		mb   dt.Mailbox
		want string
	}{
		{mb: dt.Mailbox{DisplayName: "John Doe", Email: "john@example.com"}, want: "John Doe <john@example.com>"},
		{mb: dt.Mailbox{DisplayName: "Doe, John", Email: "john@example.com"}, want: "Doe, John <john@example.com>"},
		{mb: dt.Mailbox{Email: "john@example.com"}, want: "<john@example.com>"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mb.String(); got != tt.want {
				t.Errorf("String() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestMailboxValidate(t *testing.T) {
	valid := []dt.Mailbox{
		{DisplayName: "John Doe", Email: "john@example.com"},
		{Email: "john@example.com"},
		{DisplayName: "J", Email: "jöhn@example.com"},
	}
	for _, mb := range valid {
		if err := mb.Validate(); err != nil {
			t.Errorf("%+v.Validate() = %v; want nil", mb, err)
		}
	}
	tests := []struct {
		name   string
		mb     dt.Mailbox
		causes []error
	}{
		{name: "empty email", mb: dt.Mailbox{DisplayName: "John"}, causes: []error{dt.ErrInvalidEmailAddress}},
		{name: "bad email", mb: dt.Mailbox{Email: "john@-bad-.com"}, causes: []error{dt.ErrInvalidEmailAddress, dt.ErrUnsupportedSyntax}},
		{name: "bad name", mb: dt.Mailbox{DisplayName: "<x>", Email: "john@example.com"}, causes: []error{dt.ErrInvalidEmailDisplayName, dt.ErrAngleBracket}},
		{name: "both bad", mb: dt.Mailbox{DisplayName: "a\nb", Email: "x"}, causes: []error{dt.ErrInvalidEmailDisplayName, dt.ErrInvalidEmailAddress}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mb.Validate()
			if !errors.Is(err, dt.ErrInvalidMailbox) {
				t.Fatalf("Validate() = %v; want ErrInvalidMailbox", err)
			}
			for _, cause := range tt.causes {
				checkCause(t, err, cause)
			}
		})
	}
}

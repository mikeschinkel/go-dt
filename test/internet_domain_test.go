package test

import (
	"strings"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// fakePSL stands in for golang.org/x/net/publicsuffix, which go-dt's test
// module does not require. It knows a few suffixes and reports any other
// top-level label as unlisted.
var fakePSL = dt.ClassifyPublicSuffixFunc(func(domain dt.HostDomainName) (dt.HostDomainName, dt.PublicSuffixKind) {
	for _, s := range []struct {
		suffix dt.HostDomainName
		kind   dt.PublicSuffixKind
	}{
		{"github.io", dt.PrivatePublicSuffix},
		{"co.uk", dt.ICANNPublicSuffix},
		{"com", dt.ICANNPublicSuffix},
		{"uk", dt.ICANNPublicSuffix},
		{"io", dt.ICANNPublicSuffix},
	} {
		if domain == s.suffix || strings.HasSuffix(string(domain), "."+string(s.suffix)) {
			return s.suffix, s.kind
		}
	}
	return domain[strings.LastIndexByte(string(domain), '.')+1:], dt.UnlistedPublicSuffix
})

func TestParseInternetDomain(t *testing.T) {
	withPSL := &dt.ParseInternetDomainArgs{PublicSuffixClassifier: fakePSL}
	tests := []struct {
		name        string
		input       string
		args        *dt.ParseInternetDomainArgs
		want        dt.InternetDomain
		wantErr     bool
		unsupported bool
		cause       error
	}{
		{name: "two labels", input: "example.com", want: "example.com"},
		{name: "three labels", input: "www.example.com", want: "www.example.com"},
		{name: "example is not special-use", input: "example.example", want: "example.example"},
		{name: "zero args", input: "example.com", args: &dt.ParseInternetDomainArgs{}, want: "example.com"},
		{name: "single label", input: "localhost", wantErr: true, cause: dt.ErrSingleLabel},
		{name: "test", input: "repo.test", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "localhost", input: "a.localhost", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "invalid", input: "a.invalid", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "local", input: "mikes-mbp.local", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "local any case", input: "Mikes-MBP.LOCAL", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "home.arpa", input: "printer.home.arpa", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "home.arpa itself", input: "home.arpa", wantErr: true, cause: dt.ErrSpecialUseDomain},
		{name: "arpa is not home.arpa", input: "in-addr.arpa", want: "in-addr.arpa"},
		{name: "syntax", input: "-bad-.com", wantErr: true, unsupported: true, cause: dt.ErrNonLDHLabel},
		{name: "unicode without IDNAConverter", input: "bücher.example", wantErr: true, unsupported: true, cause: dt.ErrUnicodeHostname},
		{name: "unicode with IDNAConverter", input: "bücher.example", args: &dt.ParseInternetDomainArgs{IDNAConverter: fakeIDNA}, want: "xn--bcher-kva.example"},
		{name: "IDNAConverter error", input: "badé.example", args: &dt.ParseInternetDomainArgs{IDNAConverter: fakeIDNA}, wantErr: true},
		{name: "PSL listed", input: "example.com", args: withPSL, want: "example.com"},
		{name: "PSL under private suffix", input: "foo.github.io", args: withPSL, want: "foo.github.io"},
		{name: "PSL unlisted", input: "foo.cromulent", args: withPSL, wantErr: true, cause: dt.ErrUnlistedPublicSuffix},
		{name: "PSL is ICANN suffix", input: "co.uk", args: withPSL, wantErr: true, cause: dt.ErrIsPublicSuffix},
		{name: "PSL is private suffix", input: "github.io", args: withPSL, wantErr: true, cause: dt.ErrIsPublicSuffix},
		{name: "PSL is suffix any case", input: "GitHub.IO", args: withPSL, wantErr: true, cause: dt.ErrIsPublicSuffix},
		{
			name:  "unspecified kind is unlisted",
			input: "example.com",
			args: &dt.ParseInternetDomainArgs{PublicSuffixClassifier: dt.ClassifyPublicSuffixFunc(
				func(dt.HostDomainName) (dt.HostDomainName, dt.PublicSuffixKind) { return "com", dt.UnspecifiedPublicSuffixKind },
			)},
			wantErr: true,
			cause:   dt.ErrUnlistedPublicSuffix,
		},
		{
			name:  "accepted suffix skips classifier",
			input: "foo.cromulent",
			args:  &dt.ParseInternetDomainArgs{PublicSuffixClassifier: fakePSL, AcceptedSuffixes: []dt.HostDomainName{"cromulent"}},
			want:  "foo.cromulent",
		},
		{
			name:    "accepted suffix does not override special-use",
			input:   "repo.test",
			args:    &dt.ParseInternetDomainArgs{AcceptedSuffixes: []dt.HostDomainName{"test"}},
			wantErr: true,
			cause:   dt.ErrSpecialUseDomain,
		},
		{
			name:    "rejected suffix",
			input:   "nas.LAN",
			args:    &dt.ParseInternetDomainArgs{RejectedSuffixes: []dt.HostDomainName{"lan", "corp"}},
			wantErr: true,
			cause:   dt.ErrRejectedSuffix,
		},
		{
			name:    "rejected suffix beats accepted",
			input:   "nas.lan",
			args:    &dt.ParseInternetDomainArgs{RejectedSuffixes: []dt.HostDomainName{"lan"}, AcceptedSuffixes: []dt.HostDomainName{"lan"}},
			wantErr: true,
			cause:   dt.ErrRejectedSuffix,
		},
		{
			name:  "rejected suffix matches whole labels only",
			input: "example.plan",
			args:  &dt.ParseInternetDomainArgs{RejectedSuffixes: []dt.HostDomainName{"lan"}},
			want:  "example.plan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dt.ParseInternetDomain(tt.input, tt.args)
			checkParse(t, got, err, tt.want, tt.wantErr, dt.ErrInvalidInternetDomain, tt.unsupported)
			checkCause(t, err, tt.cause)
		})
	}
}

func TestPublicSuffixKindValues(t *testing.T) {
	// The values are stored, so they must not change.
	for kind, want := range map[dt.PublicSuffixKind]int{
		dt.UnspecifiedPublicSuffixKind: 0,
		dt.UnlistedPublicSuffix:        1,
		dt.ICANNPublicSuffix:           2,
		dt.PrivatePublicSuffix:         3,
	} {
		if int(kind) != want {
			t.Errorf("kind %d; want %d", kind, want)
		}
	}
}

func TestHookAdapters(t *testing.T) {
	var classifier dt.PublicSuffixClassifier = fakePSL
	if suffix, kind := classifier.ClassifyPublicSuffix("foo.github.io"); suffix != "github.io" || kind != dt.PrivatePublicSuffix {
		t.Errorf("ClassifyPublicSuffix() = %q, %v", suffix, kind)
	}
	var conv dt.IDNAConverter = fakeIDNA
	if got, err := conv.ToASCIIHostDomainName("bücher.example"); err != nil || got != "xn--bcher-kva.example" {
		t.Errorf("ToASCIIHostDomainName() = %q, %v", got, err)
	}
}

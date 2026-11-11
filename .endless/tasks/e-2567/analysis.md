Settled with Mike on 2026-11-10 and 2026-11-11. Each name below is the RFC's own term where one exists; the doc comment on each symbol repeats the reason so a reader who prefers another name can see why this one was chosen.

## The host family

    Hostname                 RFC 3986 "host" without brackets; url.URL.Hostname(), os.Hostname()
    ├── IPAddress            netip.Addr documents itself as "an IPv4 or IPv6 address"
    │   ├── IPv4Address      RFC 3986 IPv4address
    │   └── IPv6Address      RFC 4291 "IPv6 address"; canonical text per RFC 5952
    └── HostDomainName       RFC 1123 section 2.1 "host domain name"; RFC 1034 section 3.5 preferred name syntax as amended by RFC 1123
        ├── SingleLabelHostname   RFC 6762 "single-label host name": "localhost", "buildbox"
        └── InternetDomain        existing type: two or more labels, not under a special-use suffix

    DNSName                  RFC 9499 "domain name"; x509.Certificate.DNSNames; adds underscored labels (RFC 8552) and the root dot

- Hostname, not Host: more specific, and matches url.URL.Hostname and os.Hostname. RFC 9499 notes "host name" is not formally defined and usually excludes IPs; the standard library includes them, as does RFC 1123's "identity of an Internet host".
- SingleLabelHostname, not LocalHostname: macOS's local host name (Mikes-MBP.local) has two labels, and "local" suggests mDNS's .local. Mike dislikes adjective-adjective-noun but accepted it because RFC 6762 uses the phrase.
- HostDomainName holds names that are neither single-label nor Internet domains (mikes-mbp.local, repo.test), so the family needs no type per subset. Using it also avoids HostName beside Hostname.
- InternetDomain keeps its documented meaning. It rejects names under the special-use suffixes test, localhost, invalid (RFC 6761), local (RFC 6762) and home.arpa (RFC 8375). example is not on the list: example.com resolves on the public internet.
- DNSName exists for DNS record names (_dmarc.example.com, zone-file names with the root dot). It leaves out RFC 2181's arbitrary octets on purpose: writing them as text needs RFC 1035 backslash escapes that no email or URL accepts, and underscored labels are the looseness seen in practice. Hostname and DNSName are not alternatives; they hold different sets.
- Converting up the family (InternetDomain to HostDomainName to Hostname) is a plain type conversion. Converting sideways can fail, so it is a To method returning an error, following go-dt's ToTilde and x/net/idna's ToASCII: Hostname.ToDNSName fails for an IP; DNSName.ToHostname fails for an underscored label or root dot. Nothing in go-dt or the standard library uses Convert*.
- Type aliases were considered for context-specific names (EmailDomain). An email's domain part is a Hostname, which can be localhost or an IP, so it does not share InternetDomain's constraints; no alias is needed.
- IP types are go-dt string types holding netip's canonical text, not netip.Addr, so they convert to Hostname like every other member and == works across spellings. Zones are rejected: a zone names an interface on one machine.

## Dependency hooks

go-dt's README promises "dependency-free". Build tags do not keep that promise, because go mod tidy records requirements for every build tag. Hooks keep it, in the http.Handler and http.HandlerFunc shape (an interface plus a function type with the same method):

    type PublicSuffixKind int   // explicit values, no iota (Mike's standard: stored values must survive reordering)
    const (
        UnspecifiedPublicSuffixKind PublicSuffixKind = 0
        UnlistedPublicSuffix        PublicSuffixKind = 1   // "cromulent"
        ICANNPublicSuffix           PublicSuffixKind = 2   // "com", "co.uk"
        PrivatePublicSuffix         PublicSuffixKind = 3   // "github.io"
    )
    type PublicSuffixClassifier interface { ClassifyPublicSuffix(domain string) (suffix string, kind PublicSuffixKind) }
    type ClassifyPublicSuffixFunc func(domain string) (suffix string, kind PublicSuffixKind)

    type ASCIIToer interface { ToASCII(s string) (string, error) }   // idna.Lookup satisfies it
    type ToASCIIFunc func(s string) (string, error)                   // idna.ToASCII fits it

- Why not net/http/cookiejar.PublicSuffixList: importing cookiejar pulls net/http and crypto/tls into every go-dt consumer, and its PublicSuffix(domain) string signature cannot report that a suffix is unlisted, so foo.cromulent looks as valid as foo.com. x/net's publicsuffix.PublicSuffix returns an icann flag, but icann is false for both an unknown TLD and a real private suffix such as github.io; the rule that tells them apart is "icann false and no dot", and the PublicSuffixKind return keeps that rule in one adapter instead of in every implementation. The method is named ClassifyPublicSuffix, not PublicSuffix, so it cannot collide with cookiejar's method and go-dt can add a cookiejar-compatible one later.
- Nothing else in the standard library or golang.org/x offers a public-suffix interface.
- dtx, go-dt's nested module, holds the functions that wire golang.org/x/net in, so go-dt's root go.mod stays dependency-free.
- Without the hooks, InternetDomain checks syntax plus the built-in special-use list, and Unicode input is rejected. RejectedSuffixes lets a caller reject private suffixes no RFC reserves (lan, corp); AcceptedSuffixes lets a caller accept a suffix its classifier does not know yet.
- Unicode names travel as A-labels (RFC 5890, xn--...). Converting a U-label needs IDNA tables, hence the hook. Without it, an A-label is accepted as an LDH label but go-dt cannot tell a real A-label from a Fake A-label (RFC 5890).

## Email

    EmailLocalPart      RFC 5322 local-part: dot-atom or quoted-string, quotes kept; UTF-8 per RFC 6532
    EmailAddress        local-part "@" host, where host is a Hostname; RFC 5322 addr-spec; x509 EmailAddresses
    EmailDisplayName    RFC 5322 display-name, with the email adjective because LDAP, Windows and profile APIs use "display name" with other rules
    Mailbox             RFC 5322 mailbox: { DisplayName EmailDisplayName; Email EmailAddress }

- EmailAddress is modelled as EmailLocalPart @ Hostname, the way url.URL keeps "[::1]:80" in its text while Hostname() returns "::1". The text brackets an IP as RFC 5321 section 4.1.3 requires (user@[192.0.2.1], user@[IPv6:2001:db8::1]); Hostname() strips the brackets and the IPv6: tag. RFC 5321's general address literal ([tag:content]) is rejected; its IANA registry holds only IPv6.
- A bare IP after the @ is accepted and stored bracketed. Postfix rejects it by default but accepts it with resolve_numeric_domain, and Mike would rather not lock out such users. Storing the bracketed form keeps the stored text canonical.
- The text is otherwise kept as given. Equals compares with the host folded to lower case (RFC 5321: the domain is case-insensitive, the local part may not be), so == stays exact and Equals is the semantic comparison.
- UTF-8 local parts (RFC 6531 SMTPUTF8, RFC 6532) are allowed only when the caller opts in with AllowSMTPUTF8; without it they fail with a dedicated sentinel naming the reason. RequiresSMTPUTF8 reports the fact on any value, so a mail sender can refuse while git, which stores UTF-8 freely, need not care. Validate takes no arguments and accepts UTF-8, since a value does not record how it was parsed.
- Mailbox is RFC 5322's term, which matches Mike's choice. RFC 5321 uses "Mailbox" for the bare address; go-dt covers that meaning with EmailAddress, so there is no clash in practice. The field is DisplayName, the RFC term (mail.Address uses Name).
- Mailbox has no New* constructor: in Mike's standards a New* never returns an error, and exported fields mean no constructor can guarantee validity anyway (DisplayName("<x>") compiles). It is built by ParseMailbox or a struct literal, and Validate checks either. MakeMailbox and MakeEmailAddress, fallible constructors from typed parts, are deferred: their code ships commented out with a note inviting a pull request that justifies them, because the Parse functions already cover the need.
- String renders "DisplayName <address>" unquoted, as git and trailers such as Co-Authored-By want, and "<address>" when DisplayName is empty, because git's --author rejects a bare address.

## Errors and documentation

- Each type gets an ErrInvalidX sentinel. Input that some RFC allows but the type deliberately does not (for example a label such as -bad-, which DNS can store under RFC 2181 but SMTP forbids) fails with an unsupported sentinel instead, whose message points to the doc comment. Each type's doc comment lists what it rejects and why, and invites the reader to define their own type or open a GitHub issue proposing one.

## Deferred types, declared now

ScopedIPv6Address (RFC 4007 "zone index"; RFC 6874 ZoneID), HostPort (net.SplitHostPort's hostport), PortNumber and ServiceName (RFC 6335) have no use case yet. Mike wants them declared, not only described, so they are discoverable in godoc and editors: each Parse function returns ErrNotYetSupported, and each doc comment records the options weighed and the use cases that would decide them. HostPort carries both a generic Port[T PortNumber | ServiceName]() and the typed PortNumber() and ServiceName(); Go 1.27 adds generic methods (go.dev/blog/generic-methods). The open questions recorded there: whether a scoped address belongs inside IPAddress and Hostname or only inside HostPort, and how a HostPort holding a scoped host answers Hostname().

## Consumers

go-gitutil's commit ident (E-2548) and go-fsfix's RepoFixtureArgs Author and Committer (E-2568), whose defaults become Mailbox struct literals. Tag go-dt (with dtx) before go-gitutil and go-fsfix take the dependency. The separate task replacing iota in EntryStatus is meant to ship in the same go-dt release.

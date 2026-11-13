After go-dt root is tagged with E-2567's types, bump dtx/go.mod's go-dt require to that tag and add golang.org/x/net.

Add dtx.ClassifyPublicSuffix(domain dt.HostDomainName) (dt.HostDomainName, dt.PublicSuffixKind), wrapping publicsuffix.PublicSuffix: icann false and no dot means UnlistedPublicSuffix; icann false with a dot means PrivatePublicSuffix; icann true means ICANNPublicSuffix. A caller writes dt.ClassifyPublicSuffixFunc(dtx.ClassifyPublicSuffix).

Add dtx.ToASCIIHostDomainName(name string) (dt.HostDomainName, error), wrapping idna.Lookup.ToASCII, so a caller writes dt.ToASCIIHostDomainNameFunc(dtx.ToASCIIHostDomainName) for an args struct's IDNAConverter field. go-dt's hook is IDNAConverter.ToASCIIHostDomainName, not ToASCII, so idna.Lookup does not satisfy it directly.

Test in dtx/test: foo.cromulent is unlisted, foo.github.io private, example.co.uk ICANN; ToASCIIHostDomainName converts bücher.example to xn--bcher-kva.example. Document both in dtx's README, including what IDNA means.

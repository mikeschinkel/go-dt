# Approach

go-dt ships no adapter code. go-dt's root README gains a section, under InternetDomain, holding two short adapters for callers to copy into their own code:

- A PublicSuffixClassifier adapter over golang.org/x/net/publicsuffix.PublicSuffix: icann true means ICANNPublicSuffix; icann false and the suffix has a dot means PrivatePublicSuffix; icann false and no dot means UnlistedPublicSuffix. Wrapped as dt.ClassifyPublicSuffixFunc.
- An IDNAConverter adapter over golang.org/x/net/idna: idna.Lookup.ToASCII (RFC 5891 section 5), wrapped as dt.ToASCIIHostDomainNameFunc, with a one-line note that idna.Registration (RFC 5891 section 4) is the strict alternative.

The section says why go-dt ships neither: x/net/publicsuffix imports net/http/cookiejar, which pulls net/http and crypto/tls into every consumer, and x/net/idna depends on golang.org/x/text; "a little copying is better than a little dependency". It explains what IDNA means, as the IDNAConverter doc comment does.

A new example module, examples/hook_adapters (go.mod with replace ../.. to go-dt, plus golang.org/x/net), holds the same adapter code as the README and a main that parses foo.cromulent, foo.github.io, example.co.uk, co.uk and bücher.example through ParseInternetDomain with both hooks and prints each result. `just examples` and the Release workflow already build every examples/*/ module, so the copied code is compile-checked against go-dt and x/net. The README's two adapter blocks are identical to the example's adapter code.

No change to dtx, appinfo, dtglob, or any versioned go.mod; no dependency on a go-dt release.

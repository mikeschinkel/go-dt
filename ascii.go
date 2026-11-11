package dt

// ASCIIToer converts a Unicode domain name to its ASCII form, replacing each
// U-label with an A-label ("bücher.example" becomes "xn--bcher-kva.example"),
// as RFC 5890 and RFC 5891 define. ParseInternetDomain, ParseHostname,
// ParseEmailAddress and ParseMailbox accept one through their args; without
// one they reject Unicode host names.
//
// golang.org/x/net/idna's *Profile values, such as idna.Lookup, satisfy
// ASCIIToer. go-dt does not import x/net, so its go.mod stays dependency-free.
//
// The name follows Go's -er convention for single-method interfaces, applied
// to the method name ToASCII, which x/net/idna uses.
type ASCIIToer interface {
	// ToASCII returns s with every U-label converted to an A-label, or an
	// error when s is not a valid internationalized domain name.
	ToASCII(s string) (string, error)
}

// ToASCIIFunc adapts a function to ASCIIToer, as http.HandlerFunc adapts a
// function to http.Handler. golang.org/x/net/idna.ToASCII fits it.
type ToASCIIFunc func(s string) (string, error)

// ToASCII calls f(s).
func (f ToASCIIFunc) ToASCII(s string) (string, error) {
	return f(s)
}

// hostToASCII converts a non-ASCII host name with toer, leaving ASCII input
// and a nil toer alone so the syntax check reports Unicode itself.
func hostToASCII(s string, toer ASCIIToer) (ascii string, err error) {
	ascii = s
	if toer == nil || isASCII(s) {
		goto end
	}
	ascii, err = toer.ToASCII(s)
end:
	return ascii, err
}

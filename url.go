package dt

import (
	"net/http"
	"net/url"
)

// URL is a string that contains a syntactically valid Uniform Resource Locator
// A valid URL would be parsed without error by net/url.URL.Parse().
type URL string

// GET performs an HTTP get on the URL receiver given its
func (u URL) GET(c *http.Client) (resp *http.Response, err error) {
	return c.Get(string(u))
}

func (u URL) HTTPGet(c *http.Client) (resp *http.Response, err error) {
	return c.Get(string(u))
}

func (u URL) Parse() (*url.URL, error) {
	return url.Parse(string(u))
}

func ParseURL(s string) (u URL, err error) {
	// TODO Add some validation here
	u = URL(s)
	return u, err
}

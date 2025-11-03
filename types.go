package dt

type (
	// Filename with name and extension, if exists, but no path component
	Filename string

	// FileExt is a filename extension with leading period ('.')
	FileExt string

	// Identifier has a letter or underscore, then letters, digits, or underscores.
	Identifier string

	// URL is a string that contains a syntactically valid Uniform Resource Locator
	// A valid URL would be parsed without error by net/url.URL.Parse().
	URL string

	// Version is a string uses for a software version. It is mainly without
	// constraint as people have defined versions in many different ways over time.
	Version string

	// VolumeName returns the name of the mounted volume on Windows. It might be
	// "C:" or "\\server\share". On other platforms it is always an empty string.
	VolumeName string
)

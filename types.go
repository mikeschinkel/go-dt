package dt

type (
	// Filename with name and extension, if exists, but no path component
	Filename string

	// Filepath is an absolute or relativate filepath with filename including extension if applicable
	Filepath string

	// RelFilepath is an relativate filepath with filename including extension if applicable
	RelFilepath string

	// DirPath is an absolute or relativate directory path without or trailing slash
	DirPath string

	// PathSegments is one or more path segments for a filepath, dir path, or URL
	PathSegments string

	// Identifier has a letter or underscore, then letters, digits, or underscores.
	Identifier string

	// URL is a string that contains a syntactically valid Uniform Resource Locator
	// A valid URL would be parsed without error by net/url.URL.Parse().
	URL string

	// Version is a string uses for a software version. It is mainly without
	// constraint as people have defined versions in many different ways over time.
	Version string
)

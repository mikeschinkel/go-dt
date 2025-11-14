module github.com/mikeschinkel/go-dt/dtglob

go 1.25.3

replace github.com/mikeschinkel/go-dt => ../

replace github.com/mikeschinkel/go-dt/de => ../de

require (
	github.com/bmatcuk/doublestar/v4 v4.7.1
	github.com/mikeschinkel/go-dt v0.0.0
)

require github.com/mikeschinkel/go-dt/de v0.0.0-20251107040413-53a1559d69c5 // indirect

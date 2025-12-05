module github.com/mikeschinkel/go-dt/appinfo/test

go 1.25.3

replace github.com/mikeschinkel/go-dt/appinfo => ..

replace github.com/mikeschinkel/go-dt => ../..

require (
	github.com/mikeschinkel/go-dt v0.3.2
	github.com/mikeschinkel/go-dt/appinfo v0.2.1
)

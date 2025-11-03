module github.com/mikeschinkel/go-logutil

go 1.25.3

require (
	github.com/mikeschinkel/go-cliutil v0.0.0-20251027170801-82399064d27f
	github.com/mikeschinkel/go-dt/appinfo v0.0.0-00010101000000-000000000000
)

require (
	github.com/mikeschinkel/go-dt v0.0.0-20251027222746-b5ea4e0da9da // indirect
	github.com/mikeschinkel/go-dt/de v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	github.com/mikeschinkel/go-cliutil => ../go-cliutil
	github.com/mikeschinkel/go-dt => ../go-dt
	github.com/mikeschinkel/go-dt/appinfo => ../go-dt/appinfo
	github.com/mikeschinkel/go-dt/de => ../go-dt/de
)

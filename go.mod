module github.com/mikeschinkel/go-logutil

go 1.25.3

require (
	github.com/mikeschinkel/go-cliutil v0.0.0-20251103044509-f4857eae9f54
	github.com/mikeschinkel/go-dt/appinfo v0.0.0-20251103083857-4c80f1a95372
)

require (
	github.com/mikeschinkel/go-dt v0.0.0-20251103083857-4c80f1a95372 // indirect
	github.com/mikeschinkel/go-dt/de v0.0.0-20251103083857-4c80f1a95372 // indirect
)

replace (
	github.com/mikeschinkel/go-cliutil => ../go-cliutil
	github.com/mikeschinkel/go-dt => ../go-dt
	github.com/mikeschinkel/go-dt/appinfo => ../go-dt/appinfo
	github.com/mikeschinkel/go-dt/de => ../go-dt/de
)

package main

import (
	"scuffolding/pkg/shared/skelet"
	"scuffolding/pkg/web/di"
)

var (
	project = "scuffolding"
	name    = "web"
	version = "v0.0.0"
)

func main() {
	skelet.Daemon(
		project,
		name,
		version,
		new(di.Handler), // TODO(mpavlicek): - do it better
		di.DefaultConfig(),
		di.Providers()...,
	)
}

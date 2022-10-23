package main

import (
	"icbaat/pkg/shared/skelet"
	"icbaat/pkg/web/di"
)

var (
	project = "icbaat"
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

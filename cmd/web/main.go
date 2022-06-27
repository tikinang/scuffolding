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
		di.DefaultConfig(),
		&di.Handler{}, // FIXME: mpavlicek - do it better
		di.Providers()...,
	)
}

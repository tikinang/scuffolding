package main

import (
	"icbaat/pkg/api/di"
	"icbaat/pkg/shared/skelet"
)

var (
	project = "icbaat"
	name    = "api"
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

package main

import (
	"github.com/tikinang/tikigo/skelet"
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
		new(di.Handler),
		di.Providers()...,
	)
}

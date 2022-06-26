package main

import (
	"github.com/tikinang/tikigo/skelet"
	"icbaat/pkg/web/di"
)

var (
	Project = "icbaat"
	Name    = "web"
	Version = "v0.0.0"
)

func main() {
	skelet.Daemon(
		Project,
		Name,
		Version,
		di.DefaultConfig(),
		new(di.Handler),
		di.Providers()...,
	)
}

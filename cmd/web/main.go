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
	skelet.Daemon[di.Di](
		project,
		name,
		version,
		di.DefaultConfig(),
		di.Providers()...,
	)
}

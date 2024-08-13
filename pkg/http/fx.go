package http

import "go.uber.org/fx"

var HttpModule = fx.Module(
	"http",
	fx.Provide(NewHttpServer),
)

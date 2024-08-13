package main

import (
	"moov/config"
	gateway "moov/internal/api-gateway/src/infraestructure/controllers"
	"moov/pkg/http"
	"moov/pkg/rabittmq"
	"moov/pkg/validator"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		rabittmq.NewModuleProducer,
		http.HttpModule,
		validator.Module,
		gateway.RegisterControllers,
		fx.Invoke(http.RunHttpServer),
	).Run()
}

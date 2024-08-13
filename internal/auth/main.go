package main

import (
	presenters "moov/internal/auth/src/presenters/routing"

	"go.uber.org/fx"

	"moov/config"
	databaseMongodb "moov/pkg/database/mongodb"
	"moov/pkg/rabittmq"
)

func main() {
	fx.New(
		config.Module,
		databaseMongodb.Module,
		rabittmq.NewModuleConsumer("auth"),
		fx.Invoke(
			presenters.RoutingKey,
		),
	).Run()
}

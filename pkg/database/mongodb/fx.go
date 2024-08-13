package database

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongo_database",
	fx.Provide(NewDatabase),
)

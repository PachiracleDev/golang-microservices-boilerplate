package auth

import "go.uber.org/fx"

var AuthModule = fx.Module(
	"auth_controller",
	fx.Invoke(AuthController),
)

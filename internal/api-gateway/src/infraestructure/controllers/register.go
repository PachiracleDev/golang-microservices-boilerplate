package controllers

import (
	"go.uber.org/fx"

	"moov/internal/api-gateway/src/infraestructure/controllers/auth"
	upload "moov/internal/api-gateway/src/infraestructure/controllers/upload"
)

var RegisterControllers = fx.Options(
	// Todos los controladores de la API Gateway deben ser registrados aquí
	upload.UploadModule,
	auth.AuthModule,
)

package controllers

import (
	aws "moov/pkg/aws"

	"go.uber.org/fx"
)

var UploadModule = fx.Module(
	"upload_controller",
	fx.Provide(aws.NewSDKImplementation),
	fx.Invoke(UploadController),
)

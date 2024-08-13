package controllers

import (
	"github.com/gofiber/fiber/v2"

	//USECASE

	"moov/pkg"

	//UTILS
	"github.com/google/uuid"

	"moov/config"
	"moov/pkg/http"

	awsSdk "moov/pkg/aws"
)

func UploadController(http *http.HttpServer, conf *config.Config, awsSdk *awsSdk.AwsSdkImplementation) error {

	api := http.Group("/upload")

	api.Use(http.AuthMiddleware())

	api.Get("/get-authorization", func(c *fiber.Ctx) error {
		batchKey := uuid.New().String()
		credentialsS3, errS3 := awsSdk.GetS3Token(batchKey)

		if errS3 != nil {
			return c.JSON(pkg.ResponseError(errS3.Error()))
		}

		credentialsRekognition, err := awsSdk.GetRekognitionToken()

		if err != nil {
			return c.JSON(pkg.ResponseError(err.Error()))
		}

		return c.JSON(pkg.ResponseSuccess(
			map[string]interface{}{
				"S3Credentials":          credentialsS3,
				"RekognitionCredentials": credentialsRekognition,
				"BatchKey":               batchKey,
				"Bucket":                 conf.AWS.S3Bucket,
				"Region":                 conf.AWS.Region,
			}))
	})

	return nil
}

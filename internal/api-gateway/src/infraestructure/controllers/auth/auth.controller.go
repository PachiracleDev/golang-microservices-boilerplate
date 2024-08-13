package auth

import (
	"moov/internal/api-gateway/src/domain/constants"
	"moov/internal/api-gateway/src/infraestructure/dtos"
	"moov/pkg/http"
	"moov/pkg/rabittmq"
	"moov/pkg/validator"

	protobuffers "moov/protobuffers"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/proto"
)

func AuthController(http *http.HttpServer, validate *validator.XValidator, rbtmq *rabittmq.RabbitMQ) error {

	api := http.Group("/auth")

	api.Post("/signup", func(c *fiber.Ctx) error {

		signupBody := new(dtos.SignUpDto)
		if err := c.BodyParser(signupBody); err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: err.Error(),
			}
		}

		// Validation
		if errs := validate.Validate(signupBody); errs != nil {
			return errs
		}

		return c.JSON(map[string]string{
			"token": "",
		})
	})

	api.Post("/signin", func(c *fiber.Ctx) error {

		signInBody := new(dtos.SignInDto)
		if err := c.BodyParser(signInBody); err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: err.Error(),
			}
		}

		// Validation
		if errs := validate.Validate(signInBody); errs != nil {
			return errs
		}

		protoByte, _ := proto.Marshal(&protobuffers.SignInRequest{
			Email:    signInBody.Email,
			Password: signInBody.Password,
		})

		responseByte, _ := rbtmq.SendAndListen(rabittmq.SendEvent{
			Event:      protoByte,
			Service:    constants.Auth,
			RoutingKey: "signin",
		})

		var signInResponse protobuffers.SignInResponse
		err := proto.Unmarshal(responseByte, &signInResponse)

		if err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: err.Error(),
			}
		}

		return c.JSON(&signInResponse)
	})

	return nil
}

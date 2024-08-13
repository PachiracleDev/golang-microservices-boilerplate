package http

import (
	"fmt"

	"moov/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

type HttpServer struct {
	app       *fiber.App
	jwtSecret string
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHttpServer(conf *config.Config) *HttpServer {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})
	app.Use(cors.New())

	return &HttpServer{app: app, jwtSecret: conf.JWT.Secret}
}

func RunHttpServer(server *HttpServer, conf *config.Config) error {
	go server.app.Listen(fmt.Sprintf(":%d", conf.Server.Port))
	return nil
}

func (u *HttpServer) AuthMiddleware() func(*fiber.Ctx) error {
	// JWT Middleware

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(u.jwtSecret)},
	})
}

func (u *HttpServer) GenerateToken(claims jwt.MapClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (u *HttpServer) Group(path string) fiber.Router {
	return u.app.Group(path)
}

func (u *HttpServer) App() *fiber.App {
	return u.app
}

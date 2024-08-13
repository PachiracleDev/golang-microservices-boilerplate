package controllers

import "github.com/gofiber/fiber/v2"

func UserController(app *fiber.App) {
	api := app.Group("/users")

	api.Post("/create", func(c *fiber.Ctx) error {

		return c.SendString("Hello, World!")
	})

}

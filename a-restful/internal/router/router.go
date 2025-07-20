package router

import (
	"a-resetful/internal/robot"

	"github.com/gofiber/fiber/v2"
)

func Setup() *fiber.App {
	app := fiber.New()

	repo := robot.NewMemoryRepository()
	service := robot.NewService(repo)
	handler := robot.NewHandler(service)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/commands", handler.SubmitCommands)
	app.Get("/status/:taskID", handler.GetStatus)
	app.Delete("/commands/:taskID", handler.CancelTask)

	return app
}

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
		return c.SendString("Robot API is running")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/robots", handler.GetRobots)
	v1.Post("/commands", handler.SubmitCommands)
	v1.Get("/status/:taskID", handler.GetStatus)
	v1.Delete("/commands/:taskID", handler.CancelTask)

	return app
}

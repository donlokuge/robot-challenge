package router

import (
	_ "a-resetful/docs"
	"a-resetful/internal/robot"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(robot.ErrorResponse{Message: err.Error()})
		},
	})

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

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	return app
}

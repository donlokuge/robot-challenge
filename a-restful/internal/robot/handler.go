package robot

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetRobots(c *fiber.Ctx) error {
	robots := h.service.GetAllRobots()
	return c.JSON(robots)
}

func (h *Handler) SubmitCommands(c *fiber.Ctx) error {
	var req struct {
		Commands string `json:"commands"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	taskID, err := h.service.SubmitCommands(req.Commands)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"taskID": taskID,
	})
}

func (h *Handler) GetStatus(c *fiber.Ctx) error {
	taskID := c.Params("taskID")

	state, err := h.service.GetStatus(taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(state)
}

func (h *Handler) CancelTask(c *fiber.Ctx) error {
	taskID := c.Params("taskID")

	if err := h.service.CancelTask(taskID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Task cancelled",
	})
}

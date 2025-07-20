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

// @Summary      List all robots
// @Description  Returns a list of all registered robots and their current state
// @Tags         Robots
// @Produce      json
// @Success      200 {array} RobotState
// @Router       /api/v1/robots [get]
func (h *Handler) GetRobots(c *fiber.Ctx) error {
	robots := h.service.GetAllRobots()
	return c.JSON(robots)
}

// @Summary Submit a command sequence to the robot
// @Description Accepts a sequence of commands like "N E S W" and enqueues it for execution
// @Tags Commands
// @Accept json
// @Produce json
// @Param request body SubmitCommandsRequest true "Robot command input"
// @Success 202 {object} SubmitCommandsResponse
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 422 {object} ErrorResponse "Validation or robot boundary error"
// @Router /api/v1/commands [post]
func (h *Handler) SubmitCommands(c *fiber.Ctx) error {
	var req SubmitCommandsRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	taskID, err := h.service.SubmitCommands(req.Commands)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(SubmitCommandsResponse{
		TaskID: taskID,
	})
}

// @Summary      Get robot task status
// @Description  Returns the current robot position for the given task ID
// @Tags         Status
// @Produce      json
// @Param        taskID path string true "Task ID"
// @Success      200 {object} RobotState
// @Failure      404 {object} ErrorResponse
// @Router       /api/v1/status/{taskID} [get]
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	taskID := c.Params("taskID")

	state, err := h.service.GetStatus(taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(state)
}

// @Summary      Cancel robot task
// @Description  Cancels an active robot task by task ID
// @Tags         Commands
// @Produce      json
// @Param        taskID path string true "Task ID"
// @Success      200 {object} CancelTaskResponse
// @Failure      404 {object} ErrorResponse
// @Router       /api/v1/commands/{taskID} [delete]
func (h *Handler) CancelTask(c *fiber.Ctx) error {
	taskID := c.Params("taskID")

	if err := h.service.CancelTask(taskID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(CancelTaskResponse{
		Message: "Task cancelled",
	})
}

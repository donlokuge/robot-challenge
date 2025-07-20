package robot

import (
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	repo         Repository
	taskRegistry *TaskRegistry
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:         repo,
		taskRegistry: NewTaskRegistry(),
	}
}

func (s *Service) GetAllRobots() []Robot {
	return s.repo.GetRobots()
}

func (s *Service) SubmitCommands(commands string) (string, error) {
	if err := validateCommands(commands); err != nil {
		return "", err
	}

	robots := s.repo.GetRobots()
	if len(robots) == 0 {
		return "", errors.New("no robots available")
	}

	robot := robots[0]

	startState := robot.CurrentState()
	if !isValidPath(startState, commands) {
		return "", errors.New("commands would move robot out of bounds")
	}

	taskID, posCh, errCh := robot.EnqueueTask(commands)
	s.taskRegistry.Register(taskID, robot)

	go func() {
		for {
			select {
			case <-posCh:
			case err := <-errCh:
				if err != nil {
					s.taskRegistry.Remove(taskID)
				}
				return
			}
		}
	}()

	return taskID, nil
}

func (s *Service) GetStatus(taskID string) (RobotState, error) {
	robot, ok := s.taskRegistry.Get(taskID)
	if !ok {
		return RobotState{}, fmt.Errorf("task ID not found: %s", taskID)
	}
	return robot.CurrentState(), nil
}

func (s *Service) CancelTask(taskID string) error {
	robot, ok := s.taskRegistry.Get(taskID)
	if !ok {
		return fmt.Errorf("task ID not found: %s", taskID)
	}

	if err := robot.CancelTask(taskID); err != nil {
		return err
	}
	s.taskRegistry.Remove(taskID)
	return nil
}

func validateCommands(commands string) error {
	valid := map[string]bool{"N": true, "S": true, "E": true, "W": true}
	for _, cmd := range strings.Fields(commands) {
		if !valid[cmd] {
			return fmt.Errorf("invalid command: %s", cmd)
		}
	}
	return nil
}

func isValidPath(start RobotState, commands string) bool {
	x, y := int(start.X), int(start.Y)

	for _, cmd := range strings.Fields(commands) {
		switch cmd {
		case string(North):
			y++
		case string(South):
			y--
		case string(East):
			x++
		case string(West):
			x--
		}
		if x < 0 || x >= 10 || y < 0 || y >= 10 {
			return false
		}
	}
	return true
}

package robot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SubmitCommands(t *testing.T) {
	repo := &MemoryRepository{
		robots: []Robot{NewInMemoryRobot()},
	}
	service := NewService(repo)

	t.Run("should accept valid command sequence", func(t *testing.T) {
		taskID, err := service.SubmitCommands("N E S W")
		require.NoError(t, err)
		assert.NotEmpty(t, taskID)
	})

	t.Run("should reject invalid command", func(t *testing.T) {
		taskID, err := service.SubmitCommands("N X S")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid command")
		assert.Empty(t, taskID)
	})

	t.Run("should reject out-of-bounds path", func(t *testing.T) {
		// Move north 10 times → would exceed Y = 9
		taskID, err := service.SubmitCommands("N N N N N N N N N N N")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "out of bounds")
		assert.Empty(t, taskID)
	})
}

func TestService_GetStatus(t *testing.T) {
	repo := &MemoryRepository{
		robots: []Robot{NewInMemoryRobot()},
	}
	service := NewService(repo)

	taskID, err := service.SubmitCommands("N E S")
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	state, err := service.GetStatus(taskID)
	require.NoError(t, err)
	assert.IsType(t, RobotState{}, state)
}

func TestService_CancelTask(t *testing.T) {
	repo := &MemoryRepository{
		robots: []Robot{NewInMemoryRobot()},
	}
	service := NewService(repo)

	t.Run("should return error for unknown task", func(t *testing.T) {
		err := service.CancelTask("not-found-id")
		require.Error(t, err)
	})

	t.Run("should cancel a running task", func(t *testing.T) {
		taskID, _, _ := repo.robots[0].EnqueueTask("N N N N N N N N N")
		service.taskRegistry.Register(taskID, repo.robots[0])

		time.Sleep(300 * time.Millisecond)

		err := service.CancelTask(taskID)
		require.NoError(t, err)
	})
}

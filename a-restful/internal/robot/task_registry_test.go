package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskRegistry_RegisterGetRemove(t *testing.T) {
	registry := NewTaskRegistry()
	mockRobot := NewInMemoryRobot()
	taskID := "task-123"

	t.Run("should not find unregistered task", func(t *testing.T) {
		_, ok := registry.Get(taskID)
		assert.False(t, ok)
	})

	t.Run("should register and retrieve task", func(t *testing.T) {
		registry.Register(taskID, mockRobot)
		robot, ok := registry.Get(taskID)
		assert.True(t, ok)
		assert.Equal(t, mockRobot, robot)
	})

	t.Run("should remove task", func(t *testing.T) {
		registry.Remove(taskID)
		_, ok := registry.Get(taskID)
		assert.False(t, ok)
	})
}

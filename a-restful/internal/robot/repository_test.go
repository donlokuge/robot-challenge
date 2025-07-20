package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRepository_GetRobots(t *testing.T) {
	repo := NewMemoryRepository()

	robots := repo.GetRobots()

	assert.NotNil(t, robots, "expected non-nil robot slice")
	assert.Len(t, robots, 1, "expected exactly one robot in memory")
	assert.IsType(t, &InMemoryRobot{}, robots[0], "expected robot to be of type *InMemoryRobot")
}

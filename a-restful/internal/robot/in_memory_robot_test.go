package robot

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRobot_SquarePathReturnsToOrigin(t *testing.T) {
	assert := assert.New(t)
	robot := NewInMemoryRobot()
	_, posCh, errCh := robot.EnqueueTask("N E S W")

	for pos := range posCh {
		fmt.Printf("Robot moved to position: (%d, %d)\n", pos.X, pos.Y)
	}

	select {
	case err, ok := <-errCh:
		assert.False(ok, "errCh should be closed")
		assert.Nil(err, "expected no error when moving in a square")
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timeout waiting for robot to finish square move")
	}

	state := robot.CurrentState()
	assert.Equal(uint(0), state.X, "robot X should return to 0")
	assert.Equal(uint(0), state.Y, "robot Y should return to 0")
}

func TestInMemoryRobot_ShouldNotMoveOutsideWarehouse(t *testing.T) {
	tests := []struct {
		name     string
		setup    string
		move     string
		expected string
	}{
		{
			name:     "cannot move north past Y=9",
			setup:    "N N N N N N N N N", // Y = 9
			move:     "N",
			expected: "north",
		},
		{
			name:     "cannot move south past Y=0",
			setup:    "", // Y = 0
			move:     "S",
			expected: "south",
		},
		{
			name:     "cannot move east past X=9",
			setup:    "E E E E E E E E E", // X = 9
			move:     "E",
			expected: "east",
		},
		{
			name:     "cannot move west past X=0",
			setup:    "", // X = 0
			move:     "W",
			expected: "west",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			robot := NewInMemoryRobot()

			// Move to the boundary first
			if tt.setup != "" {
				_, posCh, setupErrCh := robot.EnqueueTask(tt.setup)
				assertNoAsyncError(t, setupErrCh, posCh)
			}

			// illegal move
			_, _, errCh := robot.EnqueueTask(tt.move)
			assertAsyncErrorContains(t, errCh, tt.expected)
		})
	}
}

func assertNoAsyncError(t *testing.T, errCh <-chan error, posCh <-chan RobotState) {
	done := make(chan struct{})

	go func() {
		for range posCh {
			// just consume to unblock sender
		}
		close(done)
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err, "unexpected error during setup")
	case <-time.After(3 * time.Second):
		require.FailNow(t, "timeout waiting for setup movement to complete")
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Log("timeout waiting for posCh to close")
	}
}

func assertAsyncErrorContains(t *testing.T, errCh <-chan error, expected string) {
	select {
	case err := <-errCh:
		require.Error(t, err)
		assert.Contains(t, err.Error(), expected)
	case <-time.After(time.Second):
		require.FailNow(t, "timeout waiting for error")
	}
}

func TestInMemoryRobot_CancelTask(t *testing.T) {
	t.Run("should cancel running task", func(t *testing.T) {
		robot := NewInMemoryRobot()

		taskID, posCh, errCh := robot.EnqueueTask("N N N N N N N N N") // long-running task

		// Let the robot move
		select {
		case <-posCh:
		case <-time.After(500 * time.Millisecond):
		}

		// Cancel the task mid-execution
		err := robot.CancelTask(taskID)
		require.NoError(t, err)

		// Expect error from errCh indicating cancellation
		select {
		case err := <-errCh:
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cancelled")
		case <-time.After(2 * time.Second):
			t.Fatal("expected cancellation error but timed out")
		}
	})

	t.Run("should return error for invalid taskID", func(t *testing.T) {
		robot := NewInMemoryRobot()

		err := robot.CancelTask("non-existent-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such task")
	})
}

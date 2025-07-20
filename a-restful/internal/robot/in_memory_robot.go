package robot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var gridSize = 10

type InMemoryRobot struct {
	mu          sync.Mutex
	state       RobotState
	taskID      string
	cancelFuncs map[string]context.CancelFunc
}

func NewInMemoryRobot() *InMemoryRobot {
	return &InMemoryRobot{
		state:       RobotState{X: 0, Y: 0},
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

func (r *InMemoryRobot) EnqueueTask(commands string) (string, chan RobotState, chan error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskID := uuid.NewString()
	r.taskID = taskID

	posCh := make(chan RobotState, 1)
	errCh := make(chan error, 1)

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFuncs[taskID] = cancel

	go func() {
		defer close(posCh)
		defer close(errCh)
		defer func() {
			r.mu.Lock()
			delete(r.cancelFuncs, taskID)
			r.mu.Unlock()
		}()

		for _, raw := range strings.Fields(commands) {
			select {
			case <-ctx.Done():
				errCh <- errors.New("task cancelled")
				return
			case <-time.After(200 * time.Millisecond): // simulate movement delay
			}

			dir := Direction(raw)

			r.mu.Lock()
			switch dir {
			case North:
				if r.state.Y < uint(gridSize-1) {
					r.state.Y++
				} else {
					errCh <- fmt.Errorf("cannot move north: out of bounds at (%d,%d)", r.state.X, r.state.Y)
					r.mu.Unlock()
					return
				}
			case South:
				if r.state.Y > 0 {
					r.state.Y--
				} else {
					errCh <- fmt.Errorf("cannot move south: out of bounds at (%d,%d)", r.state.X, r.state.Y)
					r.mu.Unlock()
					return
				}
			case East:
				if r.state.X < uint(gridSize-1) {
					r.state.X++
				} else {
					errCh <- fmt.Errorf("cannot move east: out of bounds at (%d,%d)", r.state.X, r.state.Y)
					r.mu.Unlock()
					return
				}
			case West:
				if r.state.X > 0 {
					r.state.X--
				} else {
					errCh <- fmt.Errorf("cannot move west: out of bounds at (%d,%d)", r.state.X, r.state.Y)
					r.mu.Unlock()
					return
				}
			default:
				errCh <- fmt.Errorf("invalid command: %s", dir)
				r.mu.Unlock()
				return
			}

			stateCopy := r.state
			r.mu.Unlock()

			posCh <- stateCopy
		}
	}()

	return taskID, posCh, errCh
}

func (r *InMemoryRobot) CancelTask(taskID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cancel, exists := r.cancelFuncs[taskID]
	if !exists {
		return fmt.Errorf("no such task: %s", taskID)
	}

	cancel()
	delete(r.cancelFuncs, taskID)
	return nil
}

func (r *InMemoryRobot) CurrentState() RobotState {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

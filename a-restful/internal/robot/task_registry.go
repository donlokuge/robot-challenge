package robot

import (
	"sync"
)

type TaskRegistry struct {
	mu    sync.RWMutex
	tasks map[string]Robot
}

func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		tasks: make(map[string]Robot),
	}
}

func (r *TaskRegistry) Register(taskID string, robot Robot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[taskID] = robot
}

func (r *TaskRegistry) Get(taskID string) (Robot, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	robot, ok := r.tasks[taskID]
	return robot, ok
}

func (r *TaskRegistry) Remove(taskID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tasks, taskID)
}

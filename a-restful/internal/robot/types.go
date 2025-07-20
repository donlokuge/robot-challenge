package robot

type Warehouse interface {
	Robots() []Robot
}

type Robot interface {
	EnqueueTask(commands string) (taskID string, position chan RobotState, err chan error)
	CancelTask(taskID string) error
	CurrentState() RobotState
}

type RobotState struct {
	X        uint `json:"x"`
	Y        uint `json:"y"`
	HasCrate bool `json:"hasCrate"`
}

type Direction string

const (
	North Direction = "N"
	South Direction = "S"
	East  Direction = "E"
	West  Direction = "W"
)

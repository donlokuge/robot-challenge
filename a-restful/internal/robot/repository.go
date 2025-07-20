package robot

type Repository interface {
	GetRobots() []Robot
}

type MemoryRepository struct {
	robots []Robot
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		robots: []Robot{NewInMemoryRobot()},
	}
}

func (r *MemoryRepository) GetRobots() []Robot {
	return r.robots
}

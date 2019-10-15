package scheduler

import "github.com/robfig/cron/v3"

type Scheduler struct {
	tasks []*task
}

type task struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		tasks: []*task{},
	}
}

func (s *Scheduler) NewTask() {
	newTask := &task{
		cron: cron.New(),
	}
	// s.tasks = append(s.tasks,
}

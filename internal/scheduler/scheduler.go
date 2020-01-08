package scheduler

import "github.com/robfig/cron/v3"

// Scheduler is a struct that defines an instance of a scheduler.
// It holds the tasks (`[]*tasks`).
type Scheduler struct {
	tasks []*task
	c     *cron.Cron
}

type task struct {
	cron *cron.Cron
}

// TaskFunc is the function type that should be used
// when adding tasks to the scheduler
type TaskFunc func()

// nilLogger is a logger that will be used with `cron`.
// We don't need extra logging so better explicitly not logging anything ¯\_(ツ)_/¯
type nilLogger struct{}

func (l *nilLogger) Info(msg string, keysAndValues ...interface{}) {
	return
}

func (l *nilLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	return
}

// New returns a new empty scheduler.
// To this scheduler one would need to add tasks using
// `NewTask` and then start it using `Start`
// 	s := scheduler.New()
// 	s.NewTask("* * * * *", func() { fmt.Println("Hello, world") })
// 	s.Start()
func New() *Scheduler {
	return &Scheduler{
		tasks: []*task{},
		c:     cron.New(cron.WithSeconds(), cron.WithLogger(&nilLogger{})),
	}
}

// Start starts the scheduler. This should be called after adding tasks
// using NewTask, as it's currently impossible to add tasks to an already started
// scheduler.
func (s *Scheduler) Start() error {
	s.c.Start()
	return nil
}

// NewTask adds a new scheduler task. It receives a cron formatted string
// and a function of the type `TaskFunc` that's defined in this package.
func (s *Scheduler) NewTask(cronfmt string, fn func()) error {
	_, err := s.c.AddFunc(cronfmt, fn)
	return err
}

func (s *Scheduler) Stop() {
	s.c.Stop()
}

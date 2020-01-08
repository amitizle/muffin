package scheduler

import (
	"testing"
)

func TestSchedulerRunningFuncs(t *testing.T) {
	s := New()
	i := 1
	c := make(chan bool)
	err := s.NewTask("* * * * * *", func() {
		i++
		c <- true
		c <- true
	})
	if err != nil {
		t.Fatalf("could not add new task to scheduler: %v", err)
	}
	s.Start()
	<-c
	if i != 2 {
		t.Fatalf("expected scheduler function to raise i by one")
	}
	s.Stop()
}

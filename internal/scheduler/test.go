package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	// Params
	schedulerRunTime := 1 * time.Second
	startTime := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	delay := time.Millisecond * 100

	// Create Job
	job := mockJob{}
	var js []Job
	js = append(js, &job)

	// Create Schedular
	s := CreateSchedular(startTime, delay, js)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(schedulerRunTime)
		s.End()
	}()
	s.Start()
	wg.Wait()

	// Check if job executed the expected number of times
	expected := int(schedulerRunTime) / int(delay)
	if job.Counter != expected {
		t.Errorf("Expected job to execute %d times, got %d", expected, job.Counter)
	}
}

type mockJob struct {
	Counter int
}

func (m *mockJob) Execute() {
	m.Counter++
}

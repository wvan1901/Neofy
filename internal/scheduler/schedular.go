package scheduler

import "time"

type Schedular struct {
	jobs      []Job
	startTime time.Time
	delay     time.Duration
	done      chan bool
}

type Job interface {
	Execute()
}

func CreateSchedular(startTime time.Time, delay time.Duration, jobs []Job) *Schedular {
	return &Schedular{
		jobs:      jobs,
		startTime: startTime,
		delay:     delay,
		done:      make(chan bool),
	}
}

// If you dont run as a go routine it will block your program
func (s *Schedular) Start() {
	if !s.startTime.IsZero() {
		diff := time.Until(s.startTime)
		if diff < 0 {
			total := diff - s.delay
			times := total / s.delay * -1

			s.startTime = s.startTime.Add(times * s.delay)
		}
	}

	// Delays goroutine until its ready for first start time
	<-time.After(time.Until(s.startTime))

	// Ticker for scheduled job
	ticker := time.NewTicker(s.delay)
	defer ticker.Stop()

	// Run the first inital run of jobs
	s.executeJobs()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.executeJobs()
		}
	}
}

func (s *Schedular) End() {
	s.done <- true
	close(s.done)
}

func (s *Schedular) executeJobs() {
	for _, job := range s.jobs {
		job.Execute()
	}
}

func (s *Schedular) AddJob(newJob Job) {
	s.jobs = append(s.jobs, newJob)
}

package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
)

var Sch *Scheduler

func InitScheduler() {
	var err error
	Sch, err = New()
	if err != nil {
		panic(err)
	}
}

type Job interface {
	Cron() string
	Interval() int
	ID() string
	JobType() JobType
	ScheduleType() ScheduleType
	Run()
}

type Scheduler struct {
	crons   *cron.Cron
	JobChan chan Job
}

func New() (*Scheduler, error) {
	sch := &Scheduler{
		crons:   cron.New(),
		JobChan: make(chan Job),
	}

	return sch, nil
}

func (s *Scheduler) Start() {
	s.crons.Start()

	for {
		select {
		case job := <-s.JobChan:
			switch job.ScheduleType() {
			case Cron:
				log.Printf("add cron job [%s]", job.ID())
				s.crons.AddJob(job.Cron(), job)
			}
		}
	}
}

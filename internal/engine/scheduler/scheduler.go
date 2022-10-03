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

	go Sch.Start()
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
	jobs    map[string]cron.EntryID
}

func New() (*Scheduler, error) {
	sch := &Scheduler{
		crons:   cron.New(),
		JobChan: make(chan Job),
		jobs:    make(map[string]cron.EntryID),
	}

	return sch, nil
}

func (s *Scheduler) Start() {
	s.crons.Start()
	log.Printf("start scheduler")

	for {
		select {
		case job := <-s.JobChan:
			log.Printf("rev job %s", job.ID())
			switch job.ScheduleType() {
			case Cron:
				log.Printf("add cron job [%s]", job.ID())
				if entryID, err := s.crons.AddJob(job.Cron(), job); err != nil {
					s.jobs[job.ID()] = entryID
				}
			}
		}
	}
}

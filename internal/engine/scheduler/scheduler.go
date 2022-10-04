package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
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

type CronJob interface {
	cron.Job

	GetID() string
	GetCron() string
}

type Scheduler struct {
	crons *cron.Cron

	sync.RWMutex
	jobs map[string]cron.EntryID

	BackupScheduleChan chan CronJob
	UnloadJobIDChan    chan string
}

func New() (*Scheduler, error) {
	sch := &Scheduler{
		crons: cron.New(),
		jobs:  make(map[string]cron.EntryID),

		BackupScheduleChan: make(chan CronJob),
		UnloadJobIDChan:    make(chan string),
	}

	return sch, nil
}

func (s *Scheduler) Start() {
	s.crons.Start()
	log.Printf("start scheduler")

	for {
		select {
		case j := <-s.BackupScheduleChan:
			log.Printf("rev backup schedule job, id: %s", j.GetID())

			entryID, err := s.crons.AddJob(j.GetCron(), j)
			if err != nil {
				log.Printf("add cron job err, msg: %v", err)
				continue
			}

			s.Lock()
			s.jobs[j.GetID()] = entryID
			s.Unlock()

		case id := <-s.UnloadJobIDChan:
			s.RLock()
			entryID := s.jobs[id]
			s.RUnlock()

			s.Lock()
			s.crons.Remove(entryID)
			delete(s.jobs, id)
			log.Printf("unload cron job, %s", id)
			s.Unlock()
		}
	}
}

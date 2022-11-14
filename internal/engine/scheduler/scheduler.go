package scheduler

import (
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/job"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/worker"
)

type Scheduler struct {
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Run() error {
	return s.run()
}

func (s *Scheduler) run() error {
	for {
		jobs, err := getSchedulingJobs()
		if err != nil {
			return err
		}

		for _, j := range jobs {
			switch j.JobType {
			case "backup":
				backupJob, err := job.LoadBackupJob(j.ID)
				if err != nil {
					return err
				}

				worker.Pool.Submit(backupJob)
			default:
				log.Printf("unknow job type %s", j.JobType)
				continue
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func getSchedulingJobs() (jobs []job.ScheduledJob, err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	js := []repository.ScheduledJob{}
	if result := tx.Where("status = ?", "queued").Find(&js); result.Error != nil {
		return nil, result.Error
	}

	for _, j := range js {
		item, err := job.LoadScheduledJobFromModel(tx, j)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, *item)
	}

	return jobs, nil
}

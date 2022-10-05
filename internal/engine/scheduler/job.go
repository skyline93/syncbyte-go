package scheduler

import (
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/backup"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/skyline93/syncbyte-go/internal/pkg/worker"
)

type ScheduledJob struct {
	ID             uint
	Type           JobType
	BackupPolicyID uint
}

func (j *ScheduledJob) Run() {
	switch j.Type {
	case Backup:
		backuper := backup.New(repository.Db)
		if _, _, err := backuper.StartBackup(j.BackupPolicyID); err != nil {
			log.Printf("start backup failed, msg: %v", err)
		}
	}
}

func DoScheduledJob() {
	go func() {
		for {
			doScheduledJob()
			time.Sleep(time.Duration(time.Second * 10))
		}
	}()
}

func doScheduledJob() {
	var jobs []repository.ScheduledJob

	if result := repository.Db.Where("status = ?", types.Queued).Order("id").Limit(500).Find(&jobs); result.Error != nil {
		log.Printf("error: %v", result.Error)
		return
	}

	for _, job := range jobs {
		j := &ScheduledJob{
			ID:             job.ID,
			Type:           JobType(job.JobType),
			BackupPolicyID: job.BackupPolicyID,
		}

		worker.Pool.Submit(j)
	}
}

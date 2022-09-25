package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/skyline93/syncbyte-go/internal/engine/backup"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
)

var Sch *Scheduler

func InitScheduler() {
	var err error
	Sch, err = New()
	if err != nil {
		panic(err)
	}
}

type BackupScheduleJob struct {
	BackupPolicyID uint
}

func NewBackupScheduleJob(policyID uint) *BackupScheduleJob {
	return &BackupScheduleJob{
		BackupPolicyID: policyID,
	}
}

func (j BackupScheduleJob) Run() {
	backuper := backup.New(repository.Db)
	jobID, setID, err := backuper.StartBackup(j.BackupPolicyID)
	if err != nil {
		log.Printf("start backup job error, policy id <%d>", j.BackupPolicyID)
		return
	}

	log.Printf("backup job <%d> is started, backup set is <%d>", jobID, setID)
}

type Scheduler struct {
	crons        *cron.Cron
	entityIDChan chan int
}

func New() (*Scheduler, error) {
	sch := &Scheduler{
		crons: cron.New(),
	}

	if err := sch.initCronScheduler(); err != nil {
		return nil, err
	}

	return sch, nil
}

func (s *Scheduler) initCronScheduler() (err error) {
	var schs []repository.BackupPolicy
	if result := repository.Db.Where("schedule_type = ?", types.Cron).Find(&schs); result.Error != nil {
		err = result.Error
		return
	}

	for _, sch := range schs {
		cjob := NewBackupScheduleJob(sch.ID)

		log.Printf("add cron job, <%d>(%s)", sch.ID, sch.Cron)
		s.crons.AddJob(sch.Cron, cjob)
	}

	return nil
}

func (s *Scheduler) Run() {
	s.crons.Start()

	for {
		select {
		case entityID := <-s.entityIDChan:
			// TODO
			log.Printf("update cron, entityID: %v", entityID)
		}
	}
}

func (s *Scheduler) Update() {
	// TODO
}

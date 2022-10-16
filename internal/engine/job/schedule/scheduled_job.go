package schedule

import (
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/policy"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/resource"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type JobType string

const (
	Backup JobType = "backup"
)

type JobStatus string

const (
	Queued    JobStatus = "queued"
	Running   JobStatus = "running"
	Completed JobStatus = "completed"
	Failed    JobStatus = "failed"
)

type ScheduledJob struct {
	ID        uint
	Type      JobType
	JobID     uint
	Status    JobStatus
	StartTime time.Time
	EndTime   time.Time
}

func (s *ScheduledJob) Start(db *gorm.DB) (err error) {
	item := &repository.ScheduledJob{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	startTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Running, "start_time": startTime}); result.Error != nil {
		return result.Error
	}

	s.Status = Running
	s.StartTime = startTime

	return nil
}

func (s *ScheduledJob) Complete(db *gorm.DB) (err error) {
	item := &repository.ScheduledJob{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	endTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Completed, "end_time": endTime}); result.Error != nil {
		return result.Error
	}

	s.Status = Running
	s.EndTime = endTime

	return nil
}

func (s *ScheduledJob) Fail(db *gorm.DB) (err error) {
	item := &repository.ScheduledJob{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	endTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Failed, "end_time": endTime}); result.Error != nil {
		return result.Error
	}

	s.Status = Running
	s.EndTime = endTime

	return nil
}

func (s *ScheduledJob) ScheduleBackupJob(plcOpts *policy.Policy, db *gorm.DB) (err error) {
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	bj, _, err := s.addBackupJob(plcOpts, tx)
	if err != nil {
		return err
	}

	_, err = s.add(bj.ID, Backup, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScheduledJob) add(jobID uint, jobType JobType, db *gorm.DB) (j *repository.ScheduledJob, err error) {
	j = &repository.ScheduledJob{
		JobID:   jobID,
		JobType: string(jobType),
		Status:  string(types.Queued),
	}

	if result := db.Create(j); result.Error != nil {
		return nil, result.Error
	}

	return j, nil
}

func (s *ScheduledJob) addBackupJob(plcOpts *policy.Policy, db *gorm.DB) (bj *repository.BackupJob, bs *repository.BackupSet, err error) {
	source, err := resource.GetSource(plcOpts.ResourceID, db)
	if err != nil {
		return nil, nil, err
	}

	allocater := resource.Allocater{}
	backend, err := allocater.AllocateBackend(types.BackendDataTypeMapping[source.DBType], db)
	if err != nil {
		return nil, nil, err
	}

	backupHost, err := allocater.AllocateBackupHost(db)
	if err != nil {
		return nil, nil, err
	}

	bj = &repository.BackupJob{
		StartTime:  time.Now(),
		Status:     types.Queued,
		ResourceID: plcOpts.ResourceID,
		BackendID:  backend.ID,
		PolicyID:   plcOpts.ID,
		HostID:     backupHost.ID,
	}
	if result := db.Create(bj); result.Error != nil {
		return nil, nil, result.Error
	}

	bs = &repository.BackupSet{
		// DataSetName: datasetName,	// TODO 实际备份片生成后才会生成
		IsCompress:  plcOpts.IsCompress,
		BackupJobID: bj.ID,
		BackupTime:  bj.StartTime,
		ResourceID:  plcOpts.ResourceID,
		BackendID:   backend.ID, // XXX 仅在备份任务上
		Retention:   plcOpts.Retention,
	}
	if result := db.Create(bs); result.Error != nil {
		return nil, nil, result.Error
	}

	return bj, bs, nil
}

func Get(jobID uint, db *gorm.DB) (j *ScheduledJob, err error) {
	item := &repository.ScheduledJob{}
	if result := db.Where("id = ?", jobID).First(item); result.Error != nil {
		return nil, result.Error
	}

	return &ScheduledJob{
		ID:        item.ID,
		Type:      JobType(item.JobType),
		JobID:     item.JobID,
		Status:    JobStatus(item.Status),
		StartTime: item.StartTime,
		EndTime:   item.EndTime,
	}, nil
}

func GetSchedulingJobs() (js []ScheduledJob, err error) {
	items := []repository.ScheduledJob{}

	if result := repository.Db.Where("status = ?", types.Queued).Order("id").Find(&items); result.Error != nil {
		return nil, result.Error
	}

	for _, i := range items {
		j := ScheduledJob{
			ID:     i.ID,
			Type:   JobType(i.JobType),
			JobID:  i.JobID,
			Status: JobStatus(i.Status),
		}

		js = append(js, j)
	}

	return
}

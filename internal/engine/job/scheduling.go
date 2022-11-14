package job

import (
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/entity"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/gorm"
)

type ScheduledJob struct {
	ID              uint
	Status          string
	StartTime       time.Time
	EndTime         time.Time
	JobType         string
	ResourceType    string
	StorageUnitType string
	BackupSetID     uint
	Host            *entity.Host
	BackupPolicy    *entity.BackupPolicy
	StorageUnit     *entity.StorageUnit
}

func LoadScheduledJobFromModel(db *gorm.DB, j repository.ScheduledJob) (*ScheduledJob, error) {
	h, err := entity.GetHost(j.HostID, db)
	if err != nil {
		return nil, err
	}

	p, err := entity.GetBackupPolicy(j.HostID, db)
	if err != nil {
		return nil, err
	}

	s, err := entity.GetStorageUnit(j.HostID, db)
	if err != nil {
		return nil, err
	}

	return &ScheduledJob{
		ID:              j.ID,
		Status:          j.Status,
		JobType:         j.JobType,
		ResourceType:    j.ResourceType,
		StorageUnitType: j.StorageUnitType,
		BackupSetID:     j.BackupSetID,
		Host:            h,
		BackupPolicy:    p,
		StorageUnit:     s,
	}, nil
}

func (j *ScheduledJob) Start() error {
	if result := repository.Db.Model(&repository.ScheduledJob{}).Where("id = ?", j.ID).Updates(map[string]interface{}{
		"status":     "running",
		"start_time": time.Now(),
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (j *ScheduledJob) Fail() error {
	if result := repository.Db.Model(&repository.ScheduledJob{}).Where("id = ?", j.ID).Updates(map[string]interface{}{
		"status":   "failed",
		"end_time": time.Now(),
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (j *ScheduledJob) Success() error {
	if result := repository.Db.Model(&repository.ScheduledJob{}).Where("id = ?", j.ID).Updates(map[string]interface{}{
		"status":   "successed",
		"end_time": time.Now(),
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

func ScheduleBackupJob(policyID uint, db *gorm.DB) error {
	policy, err := entity.GetBackupPolicy(policyID, db)
	if err != nil {
		return err
	}

	host, err := entity.AllocateHost(db, "backup")
	if err != nil {
		return err
	}

	stu, err := entity.AllocateStu(db, entity.Local)
	if err != nil {
		return err
	}

	sj := repository.ScheduledJob{
		JobType:         "backup",
		HostID:          host.ID,
		BackupPolicyID:  policy.ID,
		StorageUnitID:   stu.ID,
		Status:          "queued",
		ResourceType:    policy.Resource.ResourceType,
		StorageUnitType: stu.StuType,
	}

	if result := db.Create(&sj); result.Error != nil {
		return result.Error
	}

	return nil
}

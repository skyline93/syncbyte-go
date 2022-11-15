package job

import (
	"encoding/json"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/client"
	"github.com/skyline93/syncbyte-go/internal/engine/entity"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
)

type JobResult struct {
	Status        string
	BackupSetSize uint64
}

type BackupJob struct {
	ScheduledJob
	client client.Client
}

func (j *BackupJob) createBackupSet() error {
	bs := repository.BackupSet{
		IsCompress: j.BackupPolicy.IsCompress,
		IsValid:    false,
		BackupTime: time.Now(),
		Retention:  j.BackupPolicy.Retention,
		ResourceID: j.BackupPolicy.Resource.ID,
	}

	if result := repository.Db.Create(&bs); result.Error != nil {
		return result.Error
	}

	if result := repository.Db.Model(&repository.ScheduledJob{}).Where("id = ?", j.ID).Updates(map[string]interface{}{
		"backup_set_id": bs.ID,
	}); result.Error != nil {
		return result.Error
	}

	j.BackupSetID = bs.ID
	return nil
}

func (j *BackupJob) UpdateBackupSetSize(size uint64) error {
	if result := repository.Db.Model(&repository.BackupSet{}).Where("id = ?", j.BackupSetID).Updates(map[string]interface{}{
		"size": size,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (j *BackupJob) Run() error {
	if err := j.Start(); err != nil {
		return err
	}

	if err := j.createBackupSet(); err != nil {
		return err
	}

	jobID, err := j.client.StartBackup(
		j.BackupPolicy.IsCompress,
		j.ResourceType,
		j.StorageUnitType,
		j.BackupPolicy.Resource.Args,
		j.StorageUnit.Args,
	)
	if err != nil {
		return err
	}

	res, err := j.client.WaitJobComplete(jobID)
	if err != nil {
		return err
	}

	jobResult := &JobResult{}
	if err = json.Unmarshal(res, jobResult); err != nil {
		return err
	}

	if err = j.UpdateBackupSetSize(jobResult.BackupSetSize); err != nil {
		return err
	}

	if err := j.Success(); err != nil {
		return err
	}
	return nil
}

func LoadBackupJob(id uint) (j *BackupJob, err error) {
	db := repository.Db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		}
		db.Commit()
	}()

	sj := repository.ScheduledJob{}

	if result := db.Where("id = ?", id).First(&sj); result.Error != nil {
		err = result.Error
		return nil, err
	}

	host, err := entity.GetHost(sj.HostID, db)
	if err != nil {
		return nil, err
	}

	pl, err := entity.GetBackupPolicy(sj.BackupPolicyID, db)
	if err != nil {
		return nil, err
	}

	stu, err := entity.GetStorageUnit(sj.StorageUnitID, db)
	if err != nil {
		return nil, err
	}

	j = &BackupJob{
		ScheduledJob: ScheduledJob{
			ID:              sj.ID,
			Status:          sj.Status,
			Host:            host,
			BackupPolicy:    pl,
			StorageUnit:     stu,
			ResourceType:    pl.Resource.ResourceType,
			StorageUnitType: stu.StuType,
		},
		client: client.NewGRPClient("127.0.0.1:50051"),
	}

	return j, nil
}

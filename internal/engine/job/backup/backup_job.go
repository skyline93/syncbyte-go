package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/agent"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/resource"
	"gorm.io/gorm"
)

type JobStatus string

const (
	Queued    JobStatus = "queued"
	Running   JobStatus = "running"
	Successed JobStatus = "successed"
	Failed    JobStatus = "failed"
)

type BackupJob struct {
	ID         uint
	Status     JobStatus
	StartTime  time.Time
	EndTime    time.Time
	PolicyID   uint
	HostID     uint
	BackendID  uint
	ResourceID uint
}

func (b *BackupJob) Start(db *gorm.DB) (err error) {
	item := &repository.BackupJob{}
	if result := db.Where("id = ?", b.ID).First(item); result.Error != nil {
		return result.Error
	}

	startTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Running, "start_time": startTime}); result.Error != nil {
		return result.Error
	}

	b.Status = Running
	b.StartTime = startTime

	return nil
}

func (b *BackupJob) Success(db *gorm.DB) (err error) {
	item := &repository.BackupJob{}
	if result := db.Where("id = ?", b.ID).First(item); result.Error != nil {
		return result.Error
	}

	endTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Successed, "end_time": endTime}); result.Error != nil {
		return result.Error
	}

	b.Status = Successed
	b.EndTime = endTime

	return nil
}

func (b *BackupJob) Fail(db *gorm.DB) (err error) {
	item := &repository.BackupJob{}
	if result := db.Where("id = ?", b.ID).First(item); result.Error != nil {
		return result.Error
	}

	endTime := time.Now()

	if result := db.Model(item).Updates(map[string]interface{}{"status": Successed, "end_time": endTime}); result.Error != nil {
		return result.Error
	}

	b.Status = Failed
	b.EndTime = endTime

	return nil
}

func (b *BackupJob) Run(ctx context.Context, db *gorm.DB) (err error) {
	policy, err := GetPolicy(b.PolicyID, db)
	if err != nil {
		return err
	}

	source, err := resource.GetSource(policy.ResourceID, db)
	if err != nil {
		return err
	}

	backend, err := resource.GetBackend(b.BackendID, db)
	if err != nil {
		return err
	}

	backupHost, err := resource.GetBackupHost(b.HostID, db)
	if err != nil {
		return err
	}

	agent, err := agent.New(fmt.Sprintf("%s:50051", backupHost.IP))
	if err != nil {
		return err
	}
	defer agent.Close()

	sourceOpts := options.SourceOption{
		Name:     source.Name,
		Server:   source.Server,
		User:     source.User,
		Password: source.Password,
		DbName:   source.DBName,
		Version:  source.Version,
		DbType:   source.DBType,
		Port:     source.Port,
	}

	backendOpts := options.BackendOption{
		EndPoint:  backend.EndPoint,
		AccessKey: backend.AccessKey,
		SecretKey: backend.SecretKey,
		Bucket:    backend.Bucket,
	}

	rep, err := agent.StartBackup(policy.IsCompress, sourceOpts, backendOpts)
	if err != nil {
		return err
	}
}

func Get(jobID uint, db *gorm.DB) (j *BackupJob, err error) {
	item := &repository.BackupJob{}
	if result := db.Where("id = ?", jobID).First(item); result.Error != nil {
		return nil, result.Error
	}

	return &BackupJob{
		ID:        item.ID,
		Status:    JobStatus(item.Status),
		StartTime: item.StartTime,
		EndTime:   item.EndTime,
	}, nil
}

package restore

import (
	"fmt"
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/agent"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
)

type Restorer struct {
	db *repository.Repository
}

func New(db *repository.Repository) *Restorer {
	return &Restorer{
		db: db,
	}
}

func (r *Restorer) createRestoreJob(backupSetID uint, opts *options.DestOptions) (rjID, destResourceID uint, err error) {
	tx := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	rj := &repository.RestoreJob{
		StartTime:   time.Now(),
		Status:      types.Running,
		BackupSetID: backupSetID,
	}
	if result := tx.Create(rj); result.Error != nil {
		err = result.Error
		return
	}

	destResource := &repository.RestoreDBResource{
		Name:         opts.Name,
		DBType:       opts.DBType,
		Version:      opts.Version,
		Server:       opts.Server,
		Port:         opts.Port,
		User:         opts.User,
		Password:     opts.Password,
		DBName:       opts.DBName,
		RestoreJobID: rj.ID,
		RestoreTime:  rj.StartTime,
	}
	if result := tx.Create(destResource); result.Error != nil {
		err = result.Error
		return
	}

	return rj.ID, destResource.ID, nil
}

func (r *Restorer) failRestoreJob(rjID uint) (err error) {
	if result := r.db.Model(&repository.RestoreJob{}).Where("id = ?", rjID).Updates(map[string]interface{}{
		"status":   types.Failed,
		"end_time": time.Now(),
	}); result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *Restorer) successRestoreJob(rjID, destResourceID uint) (err error) {
	tx := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if result := tx.Model(&repository.RestoreJob{}).Where("id = ?", rjID).Updates(map[string]interface{}{
		"status":   types.Successed,
		"end_time": time.Now(),
	}); result.Error != nil {
		err = result.Error
		return
	}

	if result := tx.Model(&repository.RestoreDBResource{}).Where("id = ?", destResourceID).Updates(map[string]interface{}{
		"is_valid": true,
	}); result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *Restorer) StartRestore(backupSetID, agentID uint, destOpts *options.DestOptions) (rjID, destResourceID uint, err error) {
	backupSet := &repository.BackupSet{}
	if result := r.db.Where("id = ?", backupSetID).First(backupSet); result.Error != nil {
		err = result.Error
		return
	}

	ag := repository.Agent{}
	if result := r.db.Where("id = ?", agentID).First(&ag); result.Error != nil {
		err = result.Error
		return
	}

	s3 := &repository.S3Backend{}
	if result := r.db.Where("id = ?", backupSet.BackendID).First(s3); result.Error != nil {
		err = result.Error
		return
	}

	backendOpts := options.BackendOption{
		EndPoint:  s3.EndPoint,
		AccessKey: s3.AccessKey,
		SecretKey: s3.SecretKey,
		Bucket:    s3.Bucket,
	}

	rjID, destResourceID, err = r.createRestoreJob(backupSetID, destOpts)
	if err != nil {
		return
	}

	agentAddr := fmt.Sprintf("%s:%d", ag.IP, ag.Port)
	go r.startRestore(rjID, destResourceID, agentAddr, backupSet.DataSetName, backupSet.IsCompress, destOpts, &backendOpts)

	return
}

func (r *Restorer) startRestore(jobID, destResourceID uint, agentAddr string, datasetName string, isUnCompress bool, destOpts *options.DestOptions, bakOpts *options.BackendOption) (err error) {
	agent, err := agent.New(agentAddr)
	if err != nil {
		return err
	}
	defer agent.Close()

	rep, err := agent.StartRestore(datasetName, isUnCompress, destOpts, bakOpts)
	if err != nil {
		return err
	}

	log.Printf("start restore job to agent, jobID: %v", rep)

	if err != nil {
		log.Printf("start restore failed to agent, error: %v", err)
		return r.failRestoreJob(jobID)
	}

	for {
		time.Sleep(time.Second * 5)
		result, err := agent.GetJobStatus(rep.Jobid)
		if err != nil {
			log.Printf("get job status error, %v", err)
			return err
		}

		log.Printf("client job status: %v", result.Status)
		switch result.Status {
		case string(types.Successed):
			r.successRestoreJob(jobID, destResourceID)
		case string(types.Failed):
			r.failRestoreJob(jobID)
		case string(types.Running):
			continue
		default:
			continue
		}

		break
	}

	log.Printf("restore completed, %s", rep.Jobid)
	return nil
}

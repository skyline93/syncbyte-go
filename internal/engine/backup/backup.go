package backup

import (
	"fmt"
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/agent"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
)

func genDataSetName(policyID uint, sourceName string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%d-%s", sourceName, policyID, t.Format("20060102150405"))
}

type Backuper struct {
	db *repository.Repository
}

func New(db *repository.Repository) *Backuper {
	return &Backuper{db: db}
}

func (b *Backuper) createBackupJob(policyID, resourceID, backendID uint, retention int, datasetName string, isCompress bool) (backupJobID, backupSetID uint, err error) {
	tx := b.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	bj := repository.BackupJob{
		StartTime:  time.Now(),
		Status:     types.Running,
		ResourceID: resourceID,
		BackendID:  backendID,
		PolicyID:   policyID,
	}
	if result := tx.Create(&bj); result.Error != nil {
		err = result.Error
		return
	}

	bs := repository.BackupSet{
		DataSetName: datasetName,
		IsCompress:  isCompress,
		BackupJobID: bj.ID,
		BackupTime:  bj.StartTime,
		ResourceID:  resourceID,
		BackendID:   backendID,
		Retention:   retention,
	}
	if result := tx.Create(&bs); result.Error != nil {
		err = result.Error
		return
	}

	return bj.ID, bs.ID, nil
}

func (b *Backuper) failBackupJob(bjID uint) (err error) {
	if result := b.db.Model(&repository.BackupJob{}).Where("id = ?", bjID).Updates(map[string]interface{}{
		"status":   types.Failed,
		"end_time": time.Now(),
	}); result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (b *Backuper) successBackupJob(bjID, bsID uint, size int64) (err error) {
	tx := b.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if result := tx.Model(&repository.BackupJob{}).Where("id = ?", bjID).Updates(map[string]interface{}{
		"status":   types.Successed,
		"end_time": time.Now(),
	}); result.Error != nil {
		err = result.Error
		return
	}

	if result := tx.Model(&repository.BackupSet{}).Where("id = ?", bsID).Updates(map[string]interface{}{
		"size":     size,
		"is_valid": true,
	}); result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (b *Backuper) StartBackup(policyID uint) (jobID, setID uint, err error) {
	pl := repository.BackupPolicy{}
	if result := b.db.Where("id = ?", policyID).First(&pl); result.Error != nil {
		err = result.Error
		return
	}

	r := repository.DBResource{}
	if result := b.db.First(&r, pl.ResourceID); result.Error != nil {
		err = result.Error
		return
	}

	s3 := repository.S3Backend{}
	if result := b.db.Where("data_type = ?", types.BackendDataTypeMapping[r.DBType]).First(&s3); result.Error != nil {
		err = result.Error
		return
	}

	ag := repository.Agent{}
	if result := b.db.Where("id = ?", pl.AgentID).First(&ag); result.Error != nil {
		err = result.Error
		return
	}

	datasetName := genDataSetName(pl.ID, r.Name)

	if jobID, setID, err = b.createBackupJob(pl.ID, pl.ResourceID, s3.ID, pl.Retention, datasetName, pl.IsCompress); err != nil {
		log.Printf("create backup job failed, error: %v", err)
		return
	}

	sourceOpts := options.SourceOption{
		Name:     r.Name,
		Server:   r.Server,
		User:     r.User,
		Password: r.Password,
		DbName:   r.DBName,
		Version:  r.Version,
		DbType:   r.DBType,
		Port:     r.Port,
	}

	backendOpts := options.BackendOption{
		EndPoint:  s3.EndPoint,
		AccessKey: s3.AccessKey,
		SecretKey: s3.SecretKey,
		Bucket:    s3.Bucket,
	}

	log.Println("start backup to backgroupd")

	agentAddr := fmt.Sprintf("%s:%d", ag.IP, ag.Port)
	b.startBackup(jobID, setID, agentAddr, datasetName, pl.IsCompress, &sourceOpts, &backendOpts)

	return
}

func (b *Backuper) startBackup(jobID, setID uint, agentAddr string, datasetName string, isCompress bool, srcOpts *options.SourceOption, bakOpts *options.BackendOption) error {
	agent, err := agent.New(agentAddr)
	if err != nil {
		return err
	}
	defer agent.Close()

	rep, err := agent.StartBackup(datasetName, isCompress, srcOpts, bakOpts)
	if err != nil {
		return err
	}

	log.Printf("start backup job to agent, jobID: %v", rep)

	if err != nil {
		log.Printf("start backup failed to agent, error: %v", err)
		return b.failBackupJob(jobID)
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
			b.successBackupJob(jobID, setID, 0)
		case string(types.Failed):
			b.failBackupJob(jobID)
		case string(types.Running):
			continue
		default:
			continue
		}

		break
	}

	log.Printf("backup completed, %s", rep.Jobid)
	return nil
}

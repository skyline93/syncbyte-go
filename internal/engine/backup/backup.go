package backup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func genDataSetName(sourceName string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", sourceName, t.Format("20060102150405"))
}

type Backuper struct {
	db *repository.Repository
}

func New(db *repository.Repository) *Backuper {
	return &Backuper{db: db}
}

func (b *Backuper) createBackupJob(resourceID, backendID uint, datasetName string, isCompress bool) (backupJobID, backupSetID uint, err error) {
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

func (b *Backuper) StartBackup(resourceID uint, isCompress bool) (jobID, setID uint, err error) {
	r := repository.DBResource{}
	if result := b.db.First(&r, resourceID); result.Error != nil {
		err = result.Error
		return
	}

	s3 := repository.S3Backend{}
	if result := b.db.Where("data_type = ?", types.BackendDataTypeMapping[r.DBType]).First(&s3); result.Error != nil {
		err = result.Error
		return
	}

	datasetName := genDataSetName(r.Name)

	if jobID, setID, err = b.createBackupJob(resourceID, s3.ID, datasetName, isCompress); err != nil {
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
	go b.startBackup(jobID, setID, datasetName, isCompress, &sourceOpts, &backendOpts)

	return
}

func (b *Backuper) startBackup(jobID, setID uint, datasetName string, isCompress bool, srcOpts *options.SourceOption, bakOpts *options.BackendOption) error {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("connect grpc server failed, error: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60*60*2)
	defer cancel()

	rep, err := client.StartBackup(ctx, &pb.BackupRequest{
		Datasetname: datasetName,
		Iscompress:  isCompress,
		SourceOpts: &pb.SourceOptions{
			Name:     srcOpts.Name,
			Server:   srcOpts.Server,
			User:     srcOpts.User,
			Password: srcOpts.Password,
			Dbname:   srcOpts.DbName,
			Version:  srcOpts.Version,
			Dbtype:   string(srcOpts.DbType),
			Port:     int32(srcOpts.Port),
		},
		BackendOpts: &pb.BackendOptions{
			Endpoint:  bakOpts.EndPoint,
			Accesskey: bakOpts.AccessKey,
			Secretkey: bakOpts.SecretKey,
			Bucket:    bakOpts.Bucket,
		},
	})

	log.Printf("start backup job to agent, jobID: %v", rep)

	if err != nil {
		log.Printf("start backup failed to agent, error: %v", err)
		return b.failBackupJob(jobID)
	}

	for {
		time.Sleep(time.Second * 5)
		rep, err := client.GetJobStatus(ctx, &pb.GetJobRequest{Jobid: rep.Jobid})
		if err != nil {
			log.Printf("get job status error, %v", err)
			return err
		}

		log.Printf("client job status: %v", rep.Status)
		switch rep.Status {
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

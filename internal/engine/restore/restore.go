package restore

import (
	"context"
	"log"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (r *Restorer) StartRestore(backupSetID uint, destOpts *options.DestOptions) (rjID, destResourceID uint, err error) {
	backupSet := &repository.BackupSet{}
	if result := r.db.Where("id = ?", backupSetID).First(backupSet); result.Error != nil {
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

	go r.startRestore(rjID, destResourceID, backupSet.DataSetName, backupSet.IsCompress, destOpts, &backendOpts)

	return
}

func (r *Restorer) startRestore(jobID, destResourceID uint, datasetName string, isUnCompress bool, destOpts *options.DestOptions, bakOpts *options.BackendOption) (err error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("connect grpc server failed, error: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60*60*2)
	defer cancel()

	rep, err := client.StartRestore(ctx, &pb.RestoreRequest{
		Datasetname:  datasetName,
		Isuncompress: isUnCompress,
		DestOpts: &pb.DestOptions{
			Name:     destOpts.Name,
			Server:   destOpts.Server,
			User:     destOpts.User,
			Password: destOpts.Password,
			Dbname:   destOpts.DBName,
			Version:  destOpts.Version,
			Dbtype:   string(destOpts.DBType),
			Port:     int32(destOpts.Port),
		},
		BackendOpts: &pb.BackendOptions{
			Endpoint:  bakOpts.EndPoint,
			Accesskey: bakOpts.AccessKey,
			Secretkey: bakOpts.SecretKey,
			Bucket:    bakOpts.Bucket,
		},
	})

	log.Printf("start restore job to agent, jobID: %v", rep)

	if err != nil {
		log.Printf("start restore failed to agent, error: %v", err)
		return r.failRestoreJob(jobID)
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

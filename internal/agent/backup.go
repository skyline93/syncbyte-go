package agent

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/skyline93/syncbyte-go/internal/agent/backend"
	"github.com/skyline93/syncbyte-go/internal/agent/source"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
)

func cleanFile(abs string) error {
	_, err := os.Stat(abs)
	if err != nil {
		log.Printf("dest file %s is not exist", abs)
		return err
	}

	log.Printf("remove dest file %s", abs)
	if err := os.RemoveAll(abs); err != nil {
		log.Printf("remove dest file err: %v", err)
		return err
	}

	return nil
}

type BackupJob struct {
	JobID       string
	Source      source.Source
	Backend     backend.Backend
	MountPoint  string
	DataSetName string
	IsCompress  bool
}

func NewBackupJob(source source.Source, backend backend.Backend, mountPoint string, datasetName string, isCompress bool) (*BackupJob, error) {
	jobID, ok := Jobs.PutWithUuid(JobInfo{Status: types.Running})
	if !ok {
		return nil, errors.New("gen job id error")
	}

	return &BackupJob{
		JobID:       jobID,
		Source:      source,
		Backend:     backend,
		MountPoint:  mountPoint,
		DataSetName: datasetName,
		IsCompress:  isCompress,
	}, nil
}

func (b *BackupJob) destFile() string {
	return filepath.Join(b.destDir(), b.DataSetName)
}

func (b *BackupJob) destDir() string {
	dir := filepath.Join(b.MountPoint, b.JobID)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(dir, 0766)
	}

	return dir
}

func (b *BackupJob) Run() (err error) {
	var size int64
	dumpFile := b.destFile()

	defer func() {
		v := Jobs.Get(b.JobID)
		jobInfo := v.(JobInfo)

		if err != nil {
			log.Printf("run backup job failed, error: %v", err)
			jobInfo.Status = types.Failed
			Jobs.Put(b.JobID, jobInfo, 60)
			return
		}

		jobInfo.Size = size
		jobInfo.Status = types.Successed
		Jobs.Put(b.JobID, jobInfo, 60)
	}()

	if err = b.Source.Dump(dumpFile); err != nil {
		log.Printf("dump file failed, error: %v", err)
		return err
	}

	if b.IsCompress {
		dumpFile, err = utils.Compress(dumpFile, true)
		if err != nil {
			log.Printf("compress dumpFile failed, error: %v", err)
			return err
		}
	}
	defer cleanFile(b.destDir())

	log.Printf("put file %s", dumpFile)
	size, err = b.Backend.Put(dumpFile)
	if err != nil {
		log.Printf("put file failed, error: %v", err)
		return err
	}

	return nil
}

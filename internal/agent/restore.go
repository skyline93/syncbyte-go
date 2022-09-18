package agent

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/skyline93/syncbyte-go/internal/agent/backend"
	"github.com/skyline93/syncbyte-go/internal/agent/dest"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
)

type RestoreJob struct {
	JobID        string
	Dest         dest.Dest
	Backend      backend.Backend
	MountPoint   string
	DataSetName  string
	IsUnCompress bool
}

func NewRestoreJob(dest dest.Dest, backend backend.Backend, mountPoint string, datasetName string, isUnCompress bool) (*RestoreJob, error) {
	jobID, ok := Jobs.PutWithUuid(JobInfo{Status: types.Running})
	if !ok {
		return nil, errors.New("gen job id error")
	}

	return &RestoreJob{
		JobID:        jobID,
		Dest:         dest,
		Backend:      backend,
		MountPoint:   mountPoint,
		DataSetName:  datasetName,
		IsUnCompress: isUnCompress,
	}, nil
}

func (r *RestoreJob) destFile() string {
	return filepath.Join(r.MountPoint, r.DataSetName)
}

func (r *RestoreJob) Run() (err error) {
	destFile := r.destFile()

	defer func() {
		v := Jobs.Get(r.JobID)
		jobInfo := v.(JobInfo)

		if err != nil {
			log.Printf("run restore job failed, error: %v", err)
			jobInfo.Status = types.Failed
			Jobs.Put(r.JobID, jobInfo, 60)
			return
		}

		jobInfo.Status = types.Successed
		Jobs.Put(r.JobID, jobInfo, 60)
	}()

	if r.IsUnCompress {
		compressFile := fmt.Sprintf("%s.gz", destFile)

		log.Printf("get file %s", compressFile)
		if err = r.Backend.Get(compressFile); err != nil {
			return err
		}

		if destFile, err = utils.UnCompress(compressFile, true); err != nil {
			return err
		}
	} else {
		if err = r.Backend.Get(destFile); err != nil {
			return err
		}
	}

	defer func() {
		var err error
		_, err = os.Stat(destFile)
		if err == nil {
			if err = os.Remove(destFile); err != nil {
				log.Printf("remove dest file err: %v", err)
				return
			}
			log.Printf("remove dest file %s", destFile)
		}
	}()

	if err = r.Dest.Build(destFile); err != nil {
		return err
	}

	return nil
}

package api

import (
	"github.com/skyline93/syncbyte-go/internal/engine/entity"
	"github.com/skyline93/syncbyte-go/internal/engine/job"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
)

func CreateBackupPolicy(
	resourceName string, retention int, isCompress bool, resourceType entity.ResourceType, scheduleType string, resourceOpts interface{}, scheduleOpts ...interface{},
) (err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = entity.CreateBackupPolicy(
		tx, resourceName, resourceType, resourceOpts, retention, isCompress, scheduleType, scheduleOpts...,
	)

	return nil
}

func AddHost(hostname, ip, hostType string) (err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = entity.AddHost(tx, hostname, ip, hostType)

	return
}

func AddStorageUnit(name string, stuType entity.StuType, options interface{}) (err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = entity.AddStorageUnit(tx, name, stuType, options)

	return
}

func StartBackupJob(policyID uint) (err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = job.ScheduleBackupJob(policyID, tx)

	return
}

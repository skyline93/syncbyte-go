package entity

import (
	"encoding/json"
	"errors"

	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ResourceType string

const (
	LocalFS  ResourceType = "localfs"
	Database ResourceType = "database"
)

type LocalFSOptions struct {
	MountPoint string
}

type PolicyStatus string

const (
	Implemented PolicyStatus = "implemented"
)

type DatabaseOptions struct {
	DbType   string
	Version  string
	Server   string
	Port     int
	User     string
	Password string
	DbName   string
}

type Resource struct {
	ID           uint
	Name         string
	ResourceType string
	Args         []byte
}

type BackupPolicy struct {
	ID         uint
	Retention  int
	IsCompress bool
	Status     string
	Resource   *Resource
}

type ScheduleOptions struct {
	Cron     string
	Interval int
}

func addResource(db *gorm.DB, name string, resType ResourceType, options interface{}) (uint, error) {
	args, err := json.Marshal(options)
	if err != nil {
		return 0, err
	}

	res := &repository.Resource{
		Name:         name,
		ResourceType: string(resType),
		Args:         datatypes.JSON(args),
	}

	if result := db.Create(res); result.Error != nil {
		return 0, result.Error
	}

	return res.ID, nil
}

func addBackupSchedule(db *gorm.DB, scheduleType string, policyID uint, args ...interface{}) error {
	sch := &repository.BackupSchedule{
		ScheduleType: scheduleType,
		IsActive:     true,
		PolicyID:     policyID,
	}

	switch scheduleType {
	case "cron":
		cron, ok := args[0].(string)
		if !ok {
			return errors.New("cron params error")
		}

		sch.Cron = cron
	case "interval":
		interval, ok := args[1].(int)
		if !ok {
			return errors.New("interval params error")
		}

		sch.Interval = interval
	default:
		return errors.New("schedule args error")
	}

	if result := db.Create(sch); result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateBackupPolicy(
	db *gorm.DB, resourceName string, resourceType ResourceType, resourceOpts interface{},
	retention int, isCompress bool, scheduleType string, scheduleArgs ...interface{},
) error {
	resourceID, err := addResource(db, resourceName, resourceType, resourceOpts)
	if err != nil {
		return err
	}

	pl := &repository.BackupPolicy{
		Retention:  retention,
		IsCompress: isCompress,
		Status:     string(Implemented),
		ResourceID: resourceID,
	}

	if result := db.Create(pl); result.Error != nil {
		return result.Error
	}

	if err := addBackupSchedule(db, scheduleType, pl.ID, scheduleArgs...); err != nil {
		return err
	}

	return nil
}

func GetBackupPolicy(id uint, db *gorm.DB) (policy *BackupPolicy, err error) {
	pl := repository.BackupPolicy{}

	if result := db.Where("id = ?", id).First(&pl); result.Error != nil {
		err = result.Error
		return nil, err
	}

	res := repository.Resource{}

	if result := db.Where("id = ?", pl.ResourceID).First(&res); result.Error != nil {
		err = result.Error
		return nil, err
	}

	policy = &BackupPolicy{
		ID:         pl.ID,
		Retention:  pl.Retention,
		IsCompress: pl.IsCompress,
		Status:     pl.Status,
		Resource: &Resource{
			ID:           res.ID,
			Name:         res.Name,
			ResourceType: res.ResourceType,
			Args:         res.Args,
		},
	}

	return policy, nil
}

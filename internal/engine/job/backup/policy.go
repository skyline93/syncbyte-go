package backup

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/gorm"
)

type Policy struct {
	ID         uint
	ResourceID uint
	Retention  int
	IsCompress bool
	Status     string
}

func (p *Policy) Add(resourceID uint, retention int, isCompress bool, schType string, cron string, interval int, db *gorm.DB) (err error) {
	item := repository.BackupPolicy{
		ResourceID:   resourceID,
		Retention:    retention,
		IsCompress:   isCompress,
		ScheduleType: schType,
		Cron:         cron,
		Frequency:    interval,
	}

	if result := db.Create(&item); result.Error != nil {
		return result.Error
	}

	return nil
}

func GetPolicy(policyID uint, db *gorm.DB) (policy *Policy, err error) {
	item := &repository.BackupPolicy{}

	if result := db.Where("id = ?", policyID).First(item); result.Error != nil {
		return nil, result.Error
	}

	return &Policy{
		ID:         item.ID,
		ResourceID: item.ResourceID,
		Retention:  item.Retention,
		IsCompress: item.IsCompress,
		Status:     item.Status,
	}, nil
}

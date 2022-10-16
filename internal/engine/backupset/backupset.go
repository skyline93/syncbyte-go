package backupset

import (
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/gorm"
)

type BackupSet struct {
	ID          uint
	Name        string
	Size        int
	BackupTime  time.Time
	Retention   int
	IsCompress  bool
	IsValid     bool
	ResourceID  uint
	BackupJobID uint
	BackendID   uint
}

func (s *BackupSet) SetBackupTime(backupTime time.Time, db *gorm.DB) (err error) {
	item := &repository.BackupSet{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	if result := db.Model(item).Updates(map[string]interface{}{"backup_time": backupTime}); result.Error != nil {
		return result.Error
	}

	s.BackupTime = backupTime

	return nil
}

func (s *BackupSet) SetSize(size int, db *gorm.DB) (err error) {
	item := &repository.BackupSet{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	if result := db.Model(item).Updates(map[string]interface{}{"size": size}); result.Error != nil {
		return result.Error
	}

	s.Size = size

	return nil
}

func (s *BackupSet) SetValidity(valid bool, db *gorm.DB) (err error) {
	item := &repository.BackupSet{}
	if result := db.Where("id = ?", s.ID).First(item); result.Error != nil {
		return result.Error
	}

	if result := db.Model(item).Updates(map[string]interface{}{"is_valid": valid}); result.Error != nil {
		return result.Error
	}

	s.IsValid = valid

	return nil
}

package repository

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model
	Name         string
	ResourceType string
	Args         datatypes.JSON `gorm:"type:jsonb" json:"args"`
}

type BackupPolicy struct {
	gorm.Model
	Retention  int
	IsCompress bool
	Status     string

	ResourceID uint
}

type BackupSchedule struct {
	gorm.Model
	ScheduleType string
	Cron         string
	Interval     int
	IsActive     bool

	PolicyID uint
}

type StorageUnit struct {
	gorm.Model
	Name    string
	StuType string

	Args datatypes.JSON `gorm:"type:jsonb" json:"args"`
}

type BackupSet struct {
	gorm.Model
	IsCompress bool
	IsValid    bool
	BackupTime time.Time
	Retention  int
	Expiration time.Time
	Size       uint64

	ResourceID uint
}

type Host struct {
	gorm.Model
	Ip       string
	HostName string
	HostType string
}

type ScheduledJob struct {
	gorm.Model
	StartTime       time.Time
	EndTime         time.Time
	Status          string
	JobType         string
	ResourceType    string
	StorageUnitType string
	Args            datatypes.JSON `gorm:"type:jsonb" json:"args"`

	HostID         uint
	BackupSetID    uint
	BackupPolicyID uint
	StorageUnitID  uint
}

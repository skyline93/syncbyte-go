package repository

import (
	"time"

	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type DBResource struct {
	gorm.Model
	Name     string
	DBType   types.DBType
	Version  string
	Server   string
	Port     int
	User     string
	Password string
	DBName   string
	Args     string
}

type S3Backend struct {
	gorm.Model
	EndPoint    string
	AccessKey   string
	SecretKey   string
	Bucket      string
	StorageType string
	DataType    types.BackendDataType
}

type RestoreJob struct {
	gorm.Model
	StartTime   time.Time
	EndTime     time.Time
	Status      types.JobStatus
	BackupSetID uint
}

type RestoreDBResource struct {
	gorm.Model
	Name         string
	DBType       types.DBType
	Version      string
	Server       string
	Port         int
	User         string
	Password     string
	DBName       string
	Args         string
	RestoreJobID uint
	IsValid      bool `gorm:"default:false"`
	RestoreTime  time.Time
}

type Agent struct {
	gorm.Model
	IP       string
	Port     int
	HostName string
	HostType string
}

// ==========================================================
type Resource struct {
	gorm.Model
	Name         string
	ResourceType string
}

type Database struct {
	gorm.Model
	DBType   string
	Version  string
	Server   string
	Port     int
	User     string
	Password string
	DbName   string

	ResourceID uint `gorm:"unique;not null"`
	Resource   Resource
}

type BackupPolicy struct {
	gorm.Model
	Retention  int
	IsCompress bool `gorm:"default:false"`
	Status     string

	ResourceID uint `gorm:"unique;not null"`
	Resource   Resource

	// to be delete
	AgentID      uint
	ScheduleType string
	Cron         string
	Frequency    int
	StartTime    types.LocalTime
	EndTime      types.LocalTime
	IsActive     bool `gorm:"default:true"`
}

type BackupSchedule struct {
	gorm.Model
	ScheduleType string
	Cron         string
	Frequency    int
	StartTime    types.LocalTime
	EndTime      types.LocalTime
	IsActive     bool `gorm:"default:true"`

	BackupPolicyID uint `gorm:"unique;not null"`
	BackupPolicy   BackupPolicy
}

type DataStorage struct {
	gorm.Model

	StorageType string
}

type S3 struct {
	gorm.Model
	EndPoint  string
	AccessKey string
	SecretKey string
	Bucket    string

	DataStorageID uint `gorm:"unique;not null"`
	DataStorage   DataStorage
}

type Local struct {
	gorm.Model
	MountPoint string

	DataStorageID uint `gorm:"unique;not null"`
	DataStorage   DataStorage
}

type BackupJob struct {
	gorm.Model
	StartTime time.Time
	EndTime   time.Time
	Status    types.JobStatus
	Args      string

	BackupPolicyID uint `gorm:"unique;not null"`
	BackupPolicy   BackupPolicy
	BackupSetID    uint `gorm:"unique;not null"`
	BackupSet      BackupSet

	// to be delete
	ResourceID uint
	PolicyID   uint
	BackendID  uint
}

type BackupSet struct {
	gorm.Model

	IsCompress bool
	IsValid    bool `gorm:"default:false"`
	Size       int
	BackupTime time.Time
	ResourceID uint
	Retention  int

	// to be delete
	BackupJobID uint
	BackendID   uint
	DataSetName string
}

type Host struct {
	gorm.Model
	IP       string
	HostName string
	HostType string
}

type ScheduledJob struct {
	gorm.Model

	JobType    string
	Status     string
	JobID      string
	ResourceID uint

	// to be delete
	BackupPolicyID uint
}

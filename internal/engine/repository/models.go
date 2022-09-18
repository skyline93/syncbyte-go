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

type BackupJob struct {
	gorm.Model
	StartTime  time.Time
	EndTime    time.Time
	Status     types.JobStatus
	ResourceID uint
	BackendID  uint
}

type BackupSet struct {
	gorm.Model
	DataSetName string
	IsCompress  bool
	IsValid     bool `gorm:"default:false"`
	Size        int
	BackupJobID uint
	BackupTime  time.Time
	ResourceID  uint
	BackendID   uint
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

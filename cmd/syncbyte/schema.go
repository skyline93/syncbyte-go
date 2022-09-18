package main

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
)

type S3Backend struct {
	ID          uint                  `json:"id"`
	EndPoint    string                `json:"endpoint"`
	AccessKey   string                `json:"accessKey"`
	SecretKey   string                `json:"secretKey"`
	Bucket      string                `json:"bucket"`
	StorageType string                `json:"storageType"`
	DataType    types.BackendDataType `json:"dataType"`
}

type Source struct {
	ID       uint         `json:"id"`
	Name     string       `json:"name"`
	Server   string       `json:"server"`
	Port     int          `json:"port"`
	User     string       `json:"user"`
	Password string       `json:"password"`
	DbName   string       `json:"dbname"`
	Extend   string       `json:"extend"`
	Version  string       `json:"version"`
	DbType   types.DBType `json:"type"`
}

type BackupJob struct {
	ID         uint            `json:"id"`
	StartTime  time.Time       `json:"start_time"`
	EndTime    time.Time       `json:"end_time"`
	Status     types.JobStatus `json:"status"`
	ResourceID uint            `json:"resource_id"`
	BackendID  uint            `json:"backend_id"`
}

type BackupSet struct {
	ID          uint      `json:"id"`
	DataSetName string    `json:"dataset_name"`
	IsCompress  bool      `json:"is_compress"`
	IsValid     bool      `json:"is_valid"`
	Size        int       `json:"size"`
	BackupJobID uint      `json:"backup_job_id"`
	BackupTime  time.Time `json:"backup_time"`
	ResourceID  uint      `json:"resource_id"`
	BackendID   uint      `json:"backend_id"`
}

type RestoreJob struct {
	ID          uint            `json:"id"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
	Status      types.JobStatus `json:"status"`
	BackupSetID uint            `json:"backup_set_id"`
}

type RestoreDBResource struct {
	ID           uint         `json:"id"`
	Name         string       `json:"name"`
	DBType       types.DBType `json:"db_type"`
	Version      string       `json:"version"`
	Server       string       `json:"server"`
	Port         int          `json:"port"`
	User         string       `json:"user"`
	Password     string       `json:"password"`
	DBName       string       `json:"dbname"`
	Args         string       `json:"args"`
	RestoreJobID uint         `json:"restore_job_id"`
	IsValid      bool         `json:"is_valid"`
	RestoreTime  time.Time    `json:"restore_time"`
}

type SystemResource struct {
	Kind types.SystemReSourceType `yaml:"kind"`
	Spec []map[string]interface{} `yaml:"spec"`
}

func NewSystemResource(kind types.SystemReSourceType) *SystemResource {
	return &SystemResource{
		Kind: kind,
		Spec: []map[string]interface{}{},
	}
}

func (s *SystemResource) ToString(in interface{}) (string, error) {
	item, err := utils.EncodeToMap(in)
	if err != nil {
		return "", err
	}
	s.Spec = append(s.Spec, item)

	result, err := utils.EncodeToBytes(s)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}

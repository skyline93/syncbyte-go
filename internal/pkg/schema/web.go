package schema

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
)

const (
	SUCCESS = iota
	ERROR
)

type ResponseBody struct {
	Code int         `json:"code"`
	Err  string      `json:"error"`
	Data interface{} `json:"data"`
}

func (r *ResponseBody) IsOk() bool {
	return r.Code == 0
}

func (r *ResponseBody) Error() error {
	return errors.New(r.Err)
}

func Response(ctx *gin.Context, data interface{}, err error) {
	switch {
	case err == nil:
		result := &ResponseBody{Err: "", Data: data, Code: SUCCESS}
		ctx.JSON(200, result)
	default:
		result := &ResponseBody{Err: err.Error(), Data: nil, Code: ERROR}
		ctx.JSON(400, result)
	}
}

type AddS3BackendRequest struct {
	EndPoint    string                `json:"endpoint"`
	AccessKey   string                `json:"accessKey"`
	SecretKey   string                `json:"secretKey"`
	Bucket      string                `json:"bucket"`
	StorageType string                `json:"storageType"`
	DataType    types.BackendDataType `json:"dataType"`
}

type AddS3BackendResponse struct {
	ID uint `json:"id"`
}

type BackupPolicyItem struct {
	Retention    int                      `json:"retention"`
	ScheduleType types.BackupScheduleType `json:"scheduleType"`
	Cron         string                   `json:"cron"`
	Frequency    int                      `json:"frequency"`
	StartTime    types.LocalTime          `json:"startTime"`
	EndTime      types.LocalTime          `json:"endTime"`
	IsCompress   bool                     `json:"isCompress"`
}

type AddSourceRequest struct {
	Name         string           `json:"name"`
	Server       string           `json:"server"`
	Port         int              `json:"port"`
	User         string           `json:"user"`
	Password     string           `json:"password"`
	DbName       string           `json:"dbname"`
	Extend       string           `json:"extend"`
	Version      string           `json:"version"`
	DbType       types.DBType     `json:"dbType"`
	BackupPolicy BackupPolicyItem `json:"backupPolicy"`
}

type AddSourceResponse struct {
	ID uint `json:"id"`
}

type StartBackupRequest struct {
	BackupPolicyID uint `json:"backupPolicyID"`
}

type StartBackupResponse struct {
	BackupJobID uint `json:"backupJobID" mapstructure:"backupJobID"`
	BackupSetID uint `json:"backupSetID" mapstructure:"backupSetID"`
}

type StartRestoreRequest struct {
	BackupSetID uint         `json:"backup_set_id"`
	Name        string       `json:"name"`
	DBType      types.DBType `json:"db_type"`
	Version     string       `json:"version"`
	Server      string       `json:"server"`
	Port        int          `json:"port"`
	User        string       `json:"user"`
	Password    string       `json:"password"`
	DBName      string       `json:"dbname"`
}

type StartRestoreResponse struct {
	RestoreJobID      uint `json:"restore_job_id" mapstructure:"restore_job_id"`
	RestoreResourceID uint `json:"restore_resource_id" mapstructure:"restore_resource_id"`
}

type AddAgentRequest struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	HostName string `json:"hostName"`
	HostType string `json:"hostType"`
}

type AddAgentResponse struct {
	AgentID uint `json:"agentID"`
}

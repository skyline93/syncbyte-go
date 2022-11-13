package webserver

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/entity"
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

type ResourceOptions struct {
	DbType   string `json:"db_type"`
	Version  string `json:"version"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`

	MountPoint string `json:"mountpoint"`
}

type ScheduleOptions struct {
	Cron     string `json:"cron"`
	Interval int    `json:"interval"`
}

type CreateBackupPolicyRequest struct {
	ResourceName    string              `json:"resource_name"`
	ResourceType    entity.ResourceType `json:"resource_type"`
	ResourceOptions ResourceOptions     `json:"resource_options"`
	Retention       int                 `json:"retention"`
	IsCompress      bool                `json:"is_compress"`
	ScheduleType    string              `json:"schedule_type"`
	ScheduleOptions ScheduleOptions     `json:"schedule_options"`
}

type AddHostRequest struct {
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
	HostType string `json:"host_type"`
}

type StuOptions struct {
	EndPoint  string `json:"endpoint"`
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	Bucket    string `json:"bucket"`

	MountPoint string `json:"mountpoint"`
}

type AddStorageUnitRequest struct {
	Name    string     `json:"name"`
	StuType string     `json:"stu_type"`
	Options StuOptions `json:"options"`
}

type StartBackupJobRequest struct {
	PolicyID uint `json:"policy_id"`
}

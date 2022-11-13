package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/api"
	"github.com/skyline93/syncbyte-go/internal/engine/entity"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateBackupPolicy(c *gin.Context) {
	req := CreateBackupPolicyRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, nil, err)
		return
	}

	if err := api.CreateBackupPolicy(
		req.ResourceName,
		req.Retention,
		req.IsCompress,
		req.ResourceType,
		req.ScheduleType,
		req.ResourceOptions,
		[]interface{}{req.ScheduleOptions.Cron, req.ScheduleOptions.Interval}...,
	); err != nil {
		Response(c, nil, err)
		return
	}

	Response(c, nil, nil)
}

func (h *Handler) AddHost(c *gin.Context) {
	req := AddHostRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, nil, err)
		return
	}

	if err := api.AddHost(req.HostName, req.IP, req.HostType); err != nil {
		Response(c, nil, err)
		return
	}

	Response(c, nil, nil)
}

func (h *Handler) AddStorageUnit(c *gin.Context) {
	req := AddStorageUnitRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, nil, err)
		return
	}

	if err := api.AddStorageUnit(req.Name, entity.StuType(req.StuType), req.Options); err != nil {
		Response(c, nil, err)
		return
	}

	Response(c, nil, nil)
}

func (h *Handler) StartBackupJob(c *gin.Context) {
	req := StartBackupJobRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, nil, err)
		return
	}

	if err := api.StartBackupJob(req.PolicyID); err != nil {
		Response(c, nil, err)
		return
	}

	Response(c, nil, nil)
}

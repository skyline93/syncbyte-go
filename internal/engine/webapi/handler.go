package webapi

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/backup"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/restore"
	"github.com/skyline93/syncbyte-go/internal/pkg/schema"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
	"gorm.io/gorm"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) StartBackup(c *gin.Context) {
	req := schema.StartBackupRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	backuper := backup.New(repository.Db)
	jobID, setID, err := backuper.StartBackup(req.BackupPolicyID)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	schema.Response(c, &schema.StartBackupResponse{BackupJobID: jobID, BackupSetID: setID}, nil)
}

func (h *Handler) StartRestore(c *gin.Context) {
	req := schema.StartRestoreRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	log.Printf("request body: %v", req)

	restorer := restore.New(repository.Db)
	rjID, resID, err := restorer.StartRestore(
		req.BackupSetID,
		req.AgentID,
		&options.DestOptions{
			Name:     req.Name,
			Server:   req.Server,
			User:     req.User,
			Password: req.Password,
			DBName:   req.DBName,
			Version:  req.Version,
			DBType:   req.DBType,
			Port:     req.Port,
		},
	)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	schema.Response(c, &schema.StartRestoreResponse{RestoreJobID: rjID, RestoreResourceID: resID}, nil)
}

func (h *Handler) AddS3Backend(c *gin.Context) {
	req := schema.AddS3BackendRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	s3 := repository.S3Backend{
		EndPoint:    req.EndPoint,
		AccessKey:   req.AccessKey,
		SecretKey:   req.SecretKey,
		Bucket:      req.Bucket,
		StorageType: req.StorageType,
		DataType:    req.DataType,
	}

	result := repository.Db.Create(&s3)
	if result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, s3.ID, nil)
}

func (h *Handler) AddAgent(c *gin.Context) {
	req := schema.AddAgentRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	agent := repository.Agent{
		IP:       req.IP,
		Port:     req.Port,
		HostName: req.HostName,
		HostType: req.HostType,
	}

	result := repository.Db.Create(&agent)
	if result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, agent.ID, nil)
}

func (h *Handler) AddDBResource(c *gin.Context) {
	var err error
	req := schema.AddSourceRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	resourceID, err := addBackupResource(tx, &req)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	_, err = backup.CreatePolicy(
		tx, resourceID, req.BackupPolicy.AgentID, req.BackupPolicy.Retention, req.BackupPolicy.IsCompress,
		req.BackupPolicy.ScheduleType, req.BackupPolicy.Cron, req.BackupPolicy.Frequency,
	)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	schema.Response(c, resourceID, nil)
}

func addBackupResource(db *gorm.DB, req *schema.AddSourceRequest) (resourceID uint, err error) {
	dbResource := repository.DBResource{
		Name:     req.Name,
		DBType:   req.DbType,
		Version:  req.Version,
		Server:   req.Server,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		DBName:   req.DbName,
		Args:     req.Extend,
	}

	if result := db.Create(&dbResource); result.Error != nil {
		return 0, result.Error
	}

	return dbResource.ID, nil
}

func (h *Handler) ListS3Backends(c *gin.Context) {
	var backends []repository.S3Backend
	if result := repository.Db.Find(&backends); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, backends, nil)
}

func (h *Handler) ListResources(c *gin.Context) {
	var resources []repository.DBResource

	limit, err := utils.StrToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	if result := repository.Db.Order("id desc").Limit(limit).Find(&resources); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, resources, nil)
}

func (h *Handler) ListBackupJobs(c *gin.Context) {
	var backupJobs []repository.BackupJob

	limit, err := utils.StrToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	if result := repository.Db.Order("start_time desc").Limit(limit).Find(&backupJobs); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, backupJobs, nil)
}

func (h *Handler) ListBackupSets(c *gin.Context) {
	var backupSets []repository.BackupSet

	limit, err := utils.StrToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	if result := repository.Db.Order("backup_time desc").Limit(limit).Find(&backupSets); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, backupSets, nil)
}

func (h *Handler) ListRestoreJobs(c *gin.Context) {
	var restoreJobs []repository.RestoreJob

	limit, err := utils.StrToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	if result := repository.Db.Order("start_time desc").Limit(limit).Find(&restoreJobs); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, restoreJobs, nil)
}

func (h *Handler) ListRestoreResources(c *gin.Context) {
	var restoreResources []repository.RestoreDBResource

	limit, err := utils.StrToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	if result := repository.Db.Order("restore_time desc").Limit(limit).Find(&restoreResources); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, restoreResources, nil)
}

func (h *Handler) ListAgents(c *gin.Context) {
	var agents []repository.Agent
	if result := repository.Db.Find(&agents); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, agents, nil)
}

func (h *Handler) EnableBackupScheduler(c *gin.Context) {
	var err error
	req := []uint{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	pls, err := backup.GetPolicies(tx, req)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	for _, pl := range pls {
		log.Printf("activate backup scheduler %d", pl.Scheduler.PolicyID)
		if err := pl.Scheduler.Activate(tx); err != nil {
			schema.Response(c, nil, err)
			return
		}
	}

	schema.Response(c, nil, nil)
}

func (h *Handler) DisableBackupScheduler(c *gin.Context) {
	var err error
	req := []uint{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	pls, err := backup.GetPolicies(tx, req)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	for _, pl := range pls {
		log.Printf("inactivate backup scheduler %d", pl.Scheduler.PolicyID)
		if err := pl.Scheduler.Inactivate(tx); err != nil {
			schema.Response(c, nil, err)
			return
		}
	}

	schema.Response(c, nil, nil)
}

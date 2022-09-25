package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/backup"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/restore"
	"github.com/skyline93/syncbyte-go/internal/pkg/schema"
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

	restorer := restore.New(repository.Db)
	rjID, resID, err := restorer.StartRestore(
		req.BackupSetID,
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

func (h *Handler) AddDBResource(c *gin.Context) {
	req := schema.AddSourceRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		schema.Response(c, nil, err)
		return
	}

	resourceID, err := addBackupResource(&req)
	if err != nil {
		schema.Response(c, nil, err)
		return
	}

	schema.Response(c, resourceID, nil)
}

func addBackupResource(req *schema.AddSourceRequest) (resourceID uint, err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

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

	if result := repository.Db.Create(&dbResource); result.Error != nil {
		return 0, result.Error
	}

	bp := req.BackupPolicy
	backupPolicy := repository.BackupPolicy{
		ResourceID:   dbResource.ID,
		Retention:    bp.Retention,
		ScheduleType: bp.ScheduleType,
		Cron:         bp.Cron,
		Frequency:    bp.Frequency,
		StartTime:    bp.StartTime,
		EndTime:      bp.EndTime,
		IsCompress:   bp.IsCompress,
	}

	if result := repository.Db.Create(&backupPolicy); result.Error != nil {
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
	if result := repository.Db.Find(&resources); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, resources, nil)
}

func (h *Handler) ListBackupJobs(c *gin.Context) {
	var backupJobs []repository.BackupJob
	if result := repository.Db.Find(&backupJobs); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, backupJobs, nil)
}

func (h *Handler) ListBackupSets(c *gin.Context) {
	var backupSets []repository.BackupSet
	if result := repository.Db.Find(&backupSets); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, backupSets, nil)
}

func (h *Handler) ListRestoreJobs(c *gin.Context) {
	var restoreJobs []repository.RestoreJob
	if result := repository.Db.Find(&restoreJobs); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, restoreJobs, nil)
}

func (h *Handler) ListRestoreResources(c *gin.Context) {
	var restoreResources []repository.RestoreDBResource
	if result := repository.Db.Find(&restoreResources); result.Error != nil {
		schema.Response(c, nil, result.Error)
		return
	}

	schema.Response(c, restoreResources, nil)
}

package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
)

type Server struct {
	*gin.Engine
}

func New() *Server {
	engine := gin.New()

	return &Server{
		Engine: engine,
	}
}

func (s *Server) Run() error {
	s.initRouter()

	return s.Engine.Run(options.Opts.ListenAddr)
}

func (s *Server) initRouter() {
	h := NewHandler()

	v1 := s.Engine.Group("/api/v1")

	backupGv1 := v1.Group("/backup")
	backupGv1.POST("", h.StartBackup)
	backupGv1.GET("/job", h.ListBackupJobs)
	backupGv1.GET("/set", h.ListBackupSets)

	restoreGv1 := v1.Group("/restore")
	restoreGv1.POST("", h.StartRestore)
	restoreGv1.GET("/job", h.ListRestoreJobs)
	restoreGv1.GET("/resource", h.ListRestoreResources)

	backendGv1 := v1.Group("/backend")
	backendGv1.POST("", h.AddS3Backend)
	backendGv1.GET("", h.ListS3Backends)

	sourceGv1 := v1.Group("/source")
	sourceGv1.POST("", h.AddDBResource)
	sourceGv1.GET("", h.ListResources)
}

package webapi

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
)

type Server struct {
	*gin.Engine
}

func New() *Server {
	engine := gin.New()
	gin.SetMode(gin.DebugMode)

	return &Server{
		Engine: engine,
	}
}

func (s *Server) Run() error {
	if err := s.initLogger(); err != nil {
		return err
	}

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

	agentGv1 := v1.Group("/agent")
	agentGv1.POST("", h.AddAgent)
	agentGv1.GET("", h.ListAgents)

	policyGv1 := v1.Group("/policy")
	policyGv1.POST("/enable", h.EnableBackupScheduler)
	policyGv1.POST("/disable", h.DisableBackupScheduler)

	manageGv1 := v1.Group("/manage")
	manageGv1.POST("/pool", h.SetPoolSize)
	manageGv1.GET("/pool", h.ListPoolWorker)
}

func (s *Server) initLogger() error {

	logFile, err := os.OpenFile(path.Join(options.Opts.LogPath, "engine.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	f, err := os.OpenFile(path.Join(options.Opts.LogPath, "access.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	s.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	return nil
}

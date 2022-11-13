package webserver

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

	backupPolicyGv1 := v1.Group("/backup/policy")
	backupPolicyGv1.POST("", h.CreateBackupPolicy)

	backupJobGv1 := v1.Group("/backup/job")
	backupJobGv1.POST("", h.StartBackupJob)

	hostGv1 := v1.Group("/host")
	hostGv1.POST("", h.AddHost)

	stuGv1 := v1.Group("/storageunit")
	stuGv1.POST("", h.AddStorageUnit)
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

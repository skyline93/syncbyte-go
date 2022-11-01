package agent

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/skyline93/syncbyte-go/internal/agent/backend"
	"github.com/skyline93/syncbyte-go/internal/agent/dest"
	"github.com/skyline93/syncbyte-go/internal/agent/source"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"github.com/skyline93/syncbyte-go/pkg/cache"
	"google.golang.org/grpc"
)

var Cache = *cache.New(4096, time.Second*60*60*6, time.Second*60)

type RPCServer struct {
	pb.UnimplementedAgentServer
}

func New() *RPCServer {
	return &RPCServer{}
}

func (s *RPCServer) StartBackup(ctx context.Context, in *pb.BackupRequest) (*pb.BackupResponse, error) {
	var (
		src source.Source
		bak backend.Backend
	)

	srcOpts := in.SourceOpts
	bakOpts := in.BackendOpts

	if in.SourceOpts.Dbtype == string(types.PostgreSQL) {
		src = source.NewPostgreSQL(&source.Options{
			Name:     srcOpts.Name,
			Server:   srcOpts.Server,
			User:     srcOpts.User,
			Password: srcOpts.Password,
			DBName:   srcOpts.Dbname,
			Version:  srcOpts.Version,
			DBType:   types.DBType(srcOpts.Dbtype),
			Port:     int(srcOpts.Port),
		})

		bak = backend.NewS3(&backend.Options{
			EndPoint:  bakOpts.Endpoint,
			AccessKey: bakOpts.Accesskey,
			SecretKey: bakOpts.Secretkey,
			Bucket:    bakOpts.Bucket,
		})
	}

	job, err := NewBackupJob(src, bak, Opts.Backup.MountPoint, in.Datasetname, in.Iscompress)
	if err != nil {
		return nil, err
	}

	go job.Run()

	return &pb.BackupResponse{Jobid: job.JobID}, nil
}

func (s *RPCServer) StartRestore(ctx context.Context, in *pb.RestoreRequest) (*pb.RestoreResponse, error) {
	var (
		d   dest.Dest
		bak backend.Backend
	)

	destOpts := in.DestOpts
	bakOpts := in.BackendOpts

	if in.DestOpts.Dbtype == string(types.PostgreSQL) {
		d = dest.NewPostgreSQL(&dest.Options{
			Name:     destOpts.Name,
			Server:   destOpts.Server,
			User:     destOpts.User,
			Password: destOpts.Password,
			DBName:   destOpts.Dbname,
			Version:  destOpts.Version,
			DBType:   types.DBType(destOpts.Dbtype),
			Port:     int(destOpts.Port),
		})

		bak = backend.NewS3(&backend.Options{
			EndPoint:  bakOpts.Endpoint,
			AccessKey: bakOpts.Accesskey,
			SecretKey: bakOpts.Secretkey,
			Bucket:    bakOpts.Bucket,
		})
	}

	job, err := NewRestoreJob(d, bak, Opts.Restore.MountPoint, in.Datasetname, in.Isuncompress)
	if err != nil {
		return nil, err
	}

	go job.Run()

	return &pb.RestoreResponse{Jobid: job.JobID}, nil
}

func (s *RPCServer) GetJobStatus(ctx context.Context, in *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	log.Printf("get job status, job_id: %v", in.Jobid)
	v := Cache.Get(in.Jobid)
	jobInfo := v.(JobInfo)

	return &pb.GetJobResponse{Status: string(jobInfo.Status)}, nil
}

func (s *RPCServer) Run() error {
	logFile, err := os.OpenFile(path.Join(Opts.LogPath, "agent.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	lis, err := net.Listen("tcp", Opts.ListenAddr)
	if err != nil {
		return errors.Wrapf(err, "listen failed, addr: %s", Opts.ListenAddr)
	}

	r := grpc.NewServer()
	pb.RegisterAgentServer(r, s)

	log.Printf("server listening at %v", lis.Addr())
	if err := r.Serve(lis); err != nil {
		return errors.Wrap(err, "serve error")
	}

	return nil
}

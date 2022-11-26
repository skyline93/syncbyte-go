package agent

import (
	"context"
	"encoding/json"

	"errors"

	"log"
	"net"

	"github.com/skyline93/syncbyte-go/internal/agent/backup"
	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"google.golang.org/grpc"
)

type RPCServer struct {
	pb.UnimplementedSyncbyteServiceServer
}

func New() *RPCServer {
	return &RPCServer{}
}

func (s *RPCServer) Call(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	req := &pb.StartBackupRequest{}

	if err := json.Unmarshal(in.Options, req); err != nil {
		return nil, err
	}

	log.Printf("action: %s, options: %v", in.Action, req)

	switch in.Action {
	case "backup":
		backuper, err := backup.NewBackuper(req.IsCompress, req.ResourceType, req.ResourceOpts, req.StuType, req.StuOpts)
		if err != nil {
			return nil, err
		}

		if err := backuper.Run(); err != nil {
			return nil, err
		}

		res, err := json.Marshal(pb.StartBackupResponse{JobID: 1})
		if err != nil {
			return nil, err
		}

		log.Printf("response backup result: %v", string(res))
		return &pb.Response{Result: res}, nil
	case "jobResult":
		result := struct {
			Status        string
			BackupSetSize uint64
		}{
			Status:        "",
			BackupSetSize: 1024,
		}

		v, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(pb.GetJobResultResponse{Status: "successed", Result: v})
		if err != nil {
			return nil, err
		}

		log.Printf("response jobResult result: %v", string(res))
		return &pb.Response{Result: res}, nil
	default:
		return nil, errors.New("unknow action")
	}
}

func (s *RPCServer) Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	r := grpc.NewServer()
	pb.RegisterSyncbyteServiceServer(r, s)

	log.Printf("server listening at %v", lis.Addr())
	if err := r.Serve(lis); err != nil {
		return err
	}

	return nil
}

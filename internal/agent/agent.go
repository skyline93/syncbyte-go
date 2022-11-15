package agent

import (
	"context"
	"encoding/json"

	"errors"

	"log"
	"net"

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

	switch in.Action {
	case "backup":
		res, err := json.Marshal(pb.StartBackupResponse{JobID: 1})
		if err != nil {
			return nil, err
		}

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

package agent

import (
	"context"
	"net"

	"github.com/skyline93/syncbyte-go/pkg/logging"
	"google.golang.org/grpc"

	pb "github.com/skyline93/syncbyte-go/internal/proto"
)

var logger = logging.GetSugaredLogger("backup")

type syncbyteServer struct {
	pb.UnimplementedSyncbyteServer
}

func (s *syncbyteServer) Backup(req *pb.BackupRequest, stream pb.Syncbyte_BackupServer) error {
	ctx := context.TODO()
	nas := NewNasVolume(req.DataStoreParams.MountPoint)
	mgmt := NewBackupManager(nas, ctx)

	fiChan := make(chan FileInfo)

	go mgmt.Backup(req.BackupParams.SourcePath, fiChan)

	for fi := range fiChan {
		var infos []*pb.PartInfo

		for _, i := range fi.PartInfos {
			info := &pb.PartInfo{
				Index: int32(i.Index),
				MD5:   i.MD5,
				Size:  i.Size,
			}

			infos = append(infos, info)
		}

		if err := stream.Send(&pb.BackupResponse{
			Name:      fi.Name,
			Path:      fi.Path,
			Size:      fi.Size,
			MD5:       fi.MD5,
			PartInfos: infos,
		}); err != nil {
			return err
		}
	}

	return nil
}

func newServer() *syncbyteServer {
	return &syncbyteServer{}
}

func RunServer(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Errorf("run server failed, err: %v", err)
		return err
	}

	logger.Infof("listen at %s", address)

	grpcServer := grpc.NewServer()
	pb.RegisterSyncbyteServer(grpcServer, newServer())

	return grpcServer.Serve(lis)
}

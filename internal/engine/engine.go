package engine

import (
	"context"
	"encoding/json"
	"io"
	"time"

	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"github.com/skyline93/syncbyte-go/pkg/logging"
	"github.com/skyline93/syncbyte-go/pkg/mongodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = logging.GetSugaredLogger("engine")

const uri = "mongodb://mongoadmin:123456@10.168.1.202:27017/?maxPoolSize=20&w=majority"

type PartInfo struct {
	Index int    `json:"index"`
	MD5   string `json:"md5"`
	Size  int64  `json:"size"`
}

type FileInfo struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Size      int64       `json:"size"`
	MD5       string      `json:"md5"`
	PartInfos []*PartInfo `json:"part_info"`
}

func (fi *FileInfo) String() string {
	v, _ := json.Marshal(fi)
	return string(v)
}

func backup(client pb.SyncbyteClient, fiChan chan FileInfo, sourcePath, mountPoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := pb.BackupRequest{
		BackupParams:    &pb.BackupParams{SourcePath: sourcePath},
		DataStoreParams: &pb.DataStoreParams{MountPoint: mountPoint},
	}

	stream, err := client.Backup(ctx, &req)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		var partInfos []*PartInfo
		for _, i := range resp.PartInfos {
			partInfo := PartInfo{
				Index: int(i.Index),
				MD5:   i.MD5,
				Size:  i.Size,
			}

			partInfos = append(partInfos, &partInfo)
		}

		fi := FileInfo{
			Name:      resp.Name,
			Path:      resp.Path,
			Size:      resp.Size,
			MD5:       resp.MD5,
			PartInfos: partInfos,
		}

		fiChan <- fi
	}

	close(fiChan)
	return nil
}

func Backup(address string, sourcePath, mountPoint string) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	mongoClient, err := mongodb.NewClient(uri)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Close()

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewSyncbyteClient(conn)

	fiChan := make(chan FileInfo)
	go backup(client, fiChan, sourcePath, mountPoint)

	col := mongoClient.GetCollection("backup01")

	for fi := range fiChan {
		if _, err := col.InsertOne(context.TODO(), fi); err != nil {
			logger.Errorf("insert error, err: %v", err)
			continue
		}

		logger.Debugf("fi: %s", fi.String())
	}

	return nil
}

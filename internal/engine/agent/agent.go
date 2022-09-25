package agent

import (
	"context"
	"time"

	"github.com/skyline93/syncbyte-go/internal/engine/options"
	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Agent struct {
	client  pb.AgentClient
	conn    *grpc.ClientConn
	context context.Context
	cancel  context.CancelFunc
}

func New(endpoint string) (*Agent, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60*60*2)

	return &Agent{
		client:  client,
		conn:    conn,
		context: ctx,
		cancel:  cancel,
	}, nil
}

func (a *Agent) Close() {
	a.cancel()
	a.conn.Close()
}

func (a *Agent) StartBackup(datasetName string, isCompress bool, srcOpts *options.SourceOption, bakOpts *options.BackendOption) (*pb.BackupResponse, error) {
	return a.client.StartBackup(a.context, &pb.BackupRequest{
		Datasetname: datasetName,
		Iscompress:  isCompress,
		SourceOpts: &pb.SourceOptions{
			Name:     srcOpts.Name,
			Server:   srcOpts.Server,
			User:     srcOpts.User,
			Password: srcOpts.Password,
			Dbname:   srcOpts.DbName,
			Version:  srcOpts.Version,
			Dbtype:   string(srcOpts.DbType),
			Port:     int32(srcOpts.Port),
		},
		BackendOpts: &pb.BackendOptions{
			Endpoint:  bakOpts.EndPoint,
			Accesskey: bakOpts.AccessKey,
			Secretkey: bakOpts.SecretKey,
			Bucket:    bakOpts.Bucket,
		},
	})
}

func (a *Agent) GetJobStatus(jobID string) (*pb.GetJobResponse, error) {
	return a.client.GetJobStatus(a.context, &pb.GetJobRequest{Jobid: jobID})
}

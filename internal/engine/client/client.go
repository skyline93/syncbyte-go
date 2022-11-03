package client

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	pb "github.com/skyline93/syncbyte-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	StartBackup(IsCompress bool, resourceType, stuType string, resourceOpts, stuOpts []byte) (jobID int, err error)
	WaitJobComplete(jobID int) (result []byte, err error)
}

type GRPClient struct {
	endpoint string
}

func NewGRPClient(endpoint string) *GRPClient {
	return &GRPClient{
		endpoint: endpoint,
	}
}

func (c *GRPClient) StartBackup(IsCompress bool, resourceType, stuType string, resourceOpts, stuOpts []byte) (jobID int, err error) {
	conn, err := grpc.Dial(c.endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewSyncbyteServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60*60*2)
	defer cancel()

	reqParams := pb.StartBackupRequest{
		IsCompress:   IsCompress,
		ResourceType: resourceType,
		ResourceOpts: resourceOpts,
		StuType:      stuType,
		StuOpts:      stuOpts,
	}

	p, err := json.Marshal(reqParams)
	if err != nil {
		return 0, err
	}

	resp, err := client.Call(ctx, &pb.Request{Action: "backup", Options: p})
	if err != nil {
		return 0, err
	}

	respParams := &pb.StartBackupResponse{}
	if err = json.Unmarshal(resp.Result, respParams); err != nil {
		return 0, err
	}

	return respParams.JobID, nil
}

func (c *GRPClient) WaitJobComplete(jobID int) (result []byte, err error) {
	conn, err := grpc.Dial(c.endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewSyncbyteServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60*60*2)
	defer cancel()

	for {
		status, result, err := c.getJobResult(ctx, client, jobID)
		if err != nil {
			return nil, err
		}

		if status == "successed" {
			return result, nil
		}

		if status == "failed" {
			return nil, errors.New("job failed")
		}

		time.Sleep(time.Second * 5)
	}
}

func (c *GRPClient) getJobResult(ctx context.Context, client pb.SyncbyteServiceClient, jobID int) (status string, result []byte, err error) {
	reqParams := pb.GetJobStatusRequest{
		JobID: jobID,
	}

	p, err := json.Marshal(reqParams)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Call(ctx, &pb.Request{Action: "jobResult", Options: p})
	if err != nil {
		return "", nil, err
	}

	respParams := &pb.GetJobResultResponse{}
	if err = json.Unmarshal(resp.Result, respParams); err != nil {
		return "", nil, err
	}

	return respParams.Status, respParams.Result, nil
}

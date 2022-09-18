package backend

import (
	"context"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Options struct {
	EndPoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type S3 struct {
	Opts *Options
}

func NewS3(opts *Options) *S3 {
	return &S3{Opts: opts}
}

func (s *S3) Put(destFile string) (size int64, err error) {
	objName := filepath.Base(destFile)

	client, err := minio.New(
		s.Opts.EndPoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(s.Opts.AccessKey, s.Opts.SecretKey, ""),
			Secure: s.Opts.UseSSL,
		},
	)
	if err != nil {
		return 0, err
	}

	objInfo, err := client.FPutObject(context.Background(), s.Opts.Bucket, objName, destFile, minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}

	return objInfo.Size, nil
}

func (s *S3) Get(destFile string) (err error) {
	objName := filepath.Base(destFile)

	client, err := minio.New(s.Opts.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.Opts.AccessKey, s.Opts.SecretKey, ""),
		Secure: s.Opts.UseSSL,
	})
	if err != nil {
		return err
	}

	err = client.FGetObject(context.Background(), s.Opts.Bucket, objName, destFile, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

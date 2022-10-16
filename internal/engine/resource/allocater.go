package resource

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type Allocater struct {
}

func (a *Allocater) AllocateBackupHost(db *gorm.DB) (host *BackupHost, err error) {
	item := &repository.Host{}
	if result := db.Where("host_type = ?", Backup).First(item); result.Error != nil {
		return nil, result.Error
	}

	return &BackupHost{
		IP:       item.IP,
		HostName: item.HostName,
		HostType: HostType(item.HostType),
	}, nil
}

func (a *Allocater) AllocateBackend(dataType types.BackendDataType, db *gorm.DB) (backend *Backend, err error) {
	item := &repository.Backend{}
	if result := db.Where("data_type = ?", dataType).First(item); result.Error != nil {
		return nil, result.Error
	}

	return &Backend{
		EndPoint:    item.EndPoint,
		AccessKey:   item.AccessKey,
		SecretKey:   item.SecretKey,
		Bucket:      item.Bucket,
		StorageType: item.StorageType,
		DataType:    item.DataType,
	}, nil
}

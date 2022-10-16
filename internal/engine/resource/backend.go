package resource

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type Backend struct {
	ID          uint
	EndPoint    string
	AccessKey   string
	SecretKey   string
	Bucket      string
	StorageType string
	DataType    types.BackendDataType
}

func GetBackend(backendID uint, db *gorm.DB) (backend *Backend, err error) {
	item := repository.Backend{}
	if result := db.First(&item, backendID); result.Error != nil {
		return nil, result.Error
	}

	return &Backend{
		ID:          item.ID,
		EndPoint:    item.EndPoint,
		AccessKey:   item.AccessKey,
		SecretKey:   item.SecretKey,
		Bucket:      item.Bucket,
		StorageType: item.StorageType,
		DataType:    item.DataType,
	}, nil
}

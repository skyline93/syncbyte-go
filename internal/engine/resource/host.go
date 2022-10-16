package resource

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/gorm"
)

type HostType string

const (
	Backup  HostType = "backup"
	Restore HostType = "restore"
)

type BackupHost struct {
	ID       uint
	IP       string
	HostName string
	HostType HostType
}

func GetBackupHost(hostID uint, db *gorm.DB) (backend *BackupHost, err error) {
	item := repository.Host{}
	if result := db.First(&item, hostID); result.Error != nil {
		return nil, result.Error
	}

	return &BackupHost{
		ID:       item.ID,
		IP:       item.IP,
		HostName: item.HostName,
		HostType: HostType(item.HostType),
	}, nil
}

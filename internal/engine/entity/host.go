package entity

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/gorm"
)

type Host struct {
	ID       uint
	Ip       string
	HostName string
	HostType string
}

func AddHost(db *gorm.DB, hostname, ip, hostType string) (err error) {
	host := &repository.Host{
		Ip:       ip,
		HostName: hostname,
		HostType: hostType,
	}

	if result := db.Create(host); result.Error != nil {
		return result.Error
	}

	return nil
}

func AllocateHost(db *gorm.DB, hostType string) (host *Host, err error) {
	h := repository.Host{}
	if result := db.Where("host_type = ?", hostType).First(&h); result.Error != nil {
		return nil, result.Error
	}

	return &Host{
		ID:       h.ID,
		Ip:       h.Ip,
		HostName: h.HostName,
		HostType: h.HostType,
	}, nil
}

func GetHost(id uint, db *gorm.DB) (host *Host, err error) {
	h := repository.Host{}
	if result := db.Where("id = ?", id).First(&h); result.Error != nil {
		err = result.Error
		return nil, err
	}

	return &Host{
		ID:       h.ID,
		Ip:       h.Ip,
		HostName: h.HostName,
		HostType: h.HostType,
	}, nil
}

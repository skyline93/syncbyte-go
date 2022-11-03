package entity

import (
	"encoding/json"

	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StuType string

const (
	Local StuType = "local"
	S3    StuType = "s3"
)

type StorageUnit struct {
	ID      uint
	Name    string
	StuType string
	Args    []byte
}

type LocalOptions struct {
	MountPoint string
}

type S3Options struct {
	EndPoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

func AddStorageUnit(db *gorm.DB, name string, stuType StuType, options interface{}) (err error) {
	args, err := json.Marshal(options)
	if err != nil {
		return err
	}

	stu := &repository.StorageUnit{
		Name:    name,
		StuType: string(stuType),
		Args:    datatypes.JSON(args),
	}

	if result := db.Create(stu); result.Error != nil {
		return result.Error
	}

	return nil
}

func AllocateStu(db *gorm.DB, stuType StuType) (stu *StorageUnit, err error) {
	h := repository.StorageUnit{}
	if result := db.Where("stu_type = ?", stuType).First(&h); result.Error != nil {
		return nil, result.Error
	}

	return &StorageUnit{
		ID:      h.ID,
		Name:    h.Name,
		StuType: h.StuType,
		Args:    h.Args,
	}, nil
}

func GetStorageUnit(id uint, db *gorm.DB) (stu *StorageUnit, err error) {
	s := repository.StorageUnit{}
	if result := db.Where("id = ?", id).First(&s); result.Error != nil {
		err = result.Error
		return nil, err
	}

	return &StorageUnit{
		ID:      s.ID,
		Name:    s.Name,
		StuType: s.StuType,
		Args:    s.Args,
	}, nil
}

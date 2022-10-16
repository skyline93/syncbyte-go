package resource

import (
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type Source struct {
	ID       uint
	Name     string
	DBType   types.DBType
	Version  string
	Server   string
	Port     int
	User     string
	Password string
	DBName   string
	Args     string
}

func GetSource(sourceID uint, db *gorm.DB) (source *Source, err error) {
	item := repository.Source{}
	if result := db.First(&item, sourceID); result.Error != nil {
		return nil, result.Error
	}

	return &Source{
		ID:       item.ID,
		Name:     item.Name,
		DBType:   item.DBType,
		Version:  item.Version,
		Server:   item.Server,
		Port:     item.Port,
		User:     item.User,
		Password: item.Password,
		DBName:   item.DBName,
		Args:     item.Args,
	}, nil
}

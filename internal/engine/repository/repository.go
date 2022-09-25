package repository

import (
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"gorm.io/gorm"
)

var Db *Repository

func InitDB() {
	var err error

	Db, err = New((*options.DatabaseOptions)(options.Opts.Database))
	if err != nil {
		panic(err)
	}

	if err = Db.AutoMigrate(
		&DBResource{},
		&S3Backend{},
		&BackupJob{},
		&BackupSet{},
		&RestoreJob{},
		&RestoreDBResource{},
		&BackupPolicy{},
	); err != nil {
		panic(err)
	}
}

type Repository struct {
	*gorm.DB
}

func New(opts *options.DatabaseOptions) (*Repository, error) {
	dia := options.Dialector(opts)

	db, err := gorm.Open(dia, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{DB: db}, nil
}

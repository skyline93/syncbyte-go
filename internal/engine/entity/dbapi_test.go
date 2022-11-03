package entity_test

import (
	"fmt"
	"testing"

	"github.com/skyline93/syncbyte-go/internal/engine/entity"
	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
)

func addPolicy() (err error) {
	tx := repository.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	options := entity.DatabaseOptions{
		DbType:   "postgresql",
		Version:  "14.5",
		Server:   "192.168.1.131",
		Port:     5432,
		User:     "syncbyte",
		Password: "lyp82nLF!?",
		DbName:   "syncbytego",
	}

	err = entity.CreateBackupPolicy(
		tx, "core-cms", entity.Database, options, 7, true, "cron", "* * * * *",
	)

	return nil
}

func TestAddResource(t *testing.T) {
	options.InitConfig()
	repository.InitDB()

	if err := addPolicy(); err != nil {
		panic(err)
	}

	pl, err := entity.GetBackupPolicy(1, repository.Db.DB)
	if err != nil {
		panic(err)
	}

	fmt.Println(pl)
}

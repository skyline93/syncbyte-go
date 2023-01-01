package agent

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skyline93/syncbyte-go/file"
	"github.com/skyline93/syncbyte-go/pkg/logging"
)

var logger = logging.GetSugaredLogger("backup")

func scanDir(root string, fChan chan string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fChan <- path
		return nil
	})

	close(fChan)
	return err
}

type BackupManager struct {
	dataStore file.DataStore
}

func NewBackupManager(store file.DataStore) *BackupManager {
	return &BackupManager{dataStore: store}
}

func (b *BackupManager) Backup(dir string) error {
	fChan := make(chan string)

	logger.Infof("scan dir in %s", dir)
	go scanDir(dir, fChan)

	for f := range fChan {
		var fi file.FileInfo

		if err := b.dataStore.UploadFile(f, &fi); err != nil {
			log.Printf("upload file failed, err: %v", err)
			continue
		}

		logger.Debugf("fileInfo: %s", fi.String())
	}

	return nil
}

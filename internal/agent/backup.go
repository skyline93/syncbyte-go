package agent

import (
	"context"
	"os"
	"path/filepath"
)

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
	dataStore DataStore
}

func NewBackupManager(store DataStore, ctx context.Context) *BackupManager {
	return &BackupManager{dataStore: store}
}

func (b *BackupManager) Backup(dir string, fiChan chan FileInfo) error {
	fChan := make(chan string)

	logger.Infof("scan dir in %s", dir)
	go scanDir(dir, fChan)

	for f := range fChan {
		var fi FileInfo

		if err := b.dataStore.UploadFile(f, &fi); err != nil {
			logger.Errorf("upload file failed, err: %v", err)
			continue
		}

		fiChan <- fi
	}

	close(fiChan)

	return nil
}

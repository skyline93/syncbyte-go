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
	wp        *WorkerPool
}

func NewBackupManager(store DataStore, ctx context.Context) *BackupManager {
	return &BackupManager{
		dataStore: store,
		wp:        NewWorkerPool(ctx, 10),
	}
}

func (b *BackupManager) Backup(dir string) error {
	fChan := make(chan string)

	logger.Infof("scan dir in %s", dir)
	go scanDir(dir, fChan)

	for f := range fChan {
		var fi FileInfo

		if err := b.dataStore.UploadFile(f, &fi); err != nil {
			logger.Errorf("upload file failed, err: %v", err)
			continue
		}

		logger.Debugf("fileInfo: %s", fi.String())
	}

	return nil
}

func (b *BackupManager) Stop() {
	b.wp.cancel()
}

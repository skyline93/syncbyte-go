package agent

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skyline93/syncbyte-go/file"
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
	dataStor file.DataStor
}

func NewBackupManager(stor file.DataStor) *BackupManager {
	return &BackupManager{dataStor: stor}
}

func (b *BackupManager) Backup(dir string) error {
	fChan := make(chan string)
	go scanDir(dir, fChan)

	for f := range fChan {
		if err := b.dataStor.UploadFile(f); err != nil {
			log.Printf("upload file failed, err: %v", err)
			continue
		}
	}

	return nil
}

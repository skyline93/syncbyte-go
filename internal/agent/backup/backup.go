package backup

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/skyline93/syncbyte-go/internal/agent/source"
	"github.com/skyline93/syncbyte-go/internal/agent/storageunit"
)

type Backuper struct {
	IsCompress  bool
	Source      source.Source
	StorageUnit storageunit.StorageUnit
}

type FSSourceOptions struct {
	MountPoint string `json:"mountpoint"`
}

type LocalDestOptions struct {
	Path string `json:"mountpoint"`
}

func NewBackuper(isCompress bool, resourceType string, resourceOpts []byte, stuType string, stuOpts []byte) (*Backuper, error) {
	var src source.Source
	var stu storageunit.StorageUnit

	// TODO other type
	if resourceType == "database" {
		fsOpts := FSSourceOptions{}
		if err := json.Unmarshal(resourceOpts, &fsOpts); err != nil {
			return nil, err
		}

		src = &source.FS{Path: fsOpts.MountPoint}
	}

	if stuType == "local" {
		localOpts := LocalDestOptions{}
		if err := json.Unmarshal(stuOpts, &localOpts); err != nil {
			return nil, err
		}

		stu = &storageunit.Local{Path: localOpts.Path}
	}

	return &Backuper{
		IsCompress:  isCompress,
		Source:      src,
		StorageUnit: stu,
	}, nil
}

func (b *Backuper) genBackupTimePoint() time.Time {
	return time.Now()
}

func (b *Backuper) Run() error {
	bt := b.genBackupTimePoint()

	if err := b.Source.BeforeBackup(); err != nil {
		return err
	}

	if err := b.transport(bt); err != nil {
		return err
	}

	return nil
}

func (b *Backuper) transport(backupTime time.Time) error {
	sourcePath := b.Source.GetSourcePath()

	fChan := make(chan string)

	go ScanDir(sourcePath, fChan)

	func() {
		for f := range fChan {
			fp, err := os.Open(f)
			if err != nil {
				log.Printf("open file [%s] failed", f)
				continue
			}

			log.Printf("save file [%s]", f)
			if err := b.StorageUnit.Save(path.Base(f), fp); err != nil {
				log.Printf("stu save file [%s] failed, err: %v", f, err)
				continue
			}
		}
	}()

	return nil
}

func ScanDir(root string, fChan chan string) error {
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

package backup

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
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
	path := b.Source.GetSourcePath()

	fChan := make(chan fs.FileInfo)

	go func() {
		for fi := range fChan {
			log.Printf("open file [%s]", fi.Name())
			fp, err := os.Open(fi.Name())
			if err != nil {
				log.Printf("open file [%s] failed", fi.Name())
				continue
			}

			log.Printf("save file [%s]", fi.Name())
			if err := b.StorageUnit.Save(fi.Name(), fp); err != nil {
				log.Printf("stu save file [%s] failed", fi.Name())
				continue
			}
		}
	}()

	b.ScanDir(path, fChan)

	close(fChan)
	return nil
}

func (b *Backuper) ScanDir(root string, fChan chan fs.FileInfo) error {
	if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return b.ScanDir(path, fChan)
		}

		log.Printf("scaned file [%v]", info)
		fChan <- info
		return nil
	}); err != nil {
		return err
	}

	return nil
}

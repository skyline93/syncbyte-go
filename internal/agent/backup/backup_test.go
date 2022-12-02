package backup_test

import (
	"testing"

	"github.com/skyline93/syncbyte-go/internal/agent/backup"
)

func TestBackupTransport(t *testing.T) {
	var (
		isCompress   bool   = true
		resourceType string = "fs"
		resourceOpts []byte = []byte(`{"mountpoint": "/home/greene/workspace/syncbyte-go/output"}`)
		stuType      string = "local"
		stuOpts      []byte = []byte(`{"mountpoint": "/home/greene/workspace/syncbyte-go/testgo"}`)
	)

	backuper, err := backup.NewBackuper(isCompress, resourceType, resourceOpts, stuType, stuOpts)
	if err != nil {
		t.Fatal("new backuper failed")
	}

	if err := backuper.Run(); err != nil {
		t.Fatal("bakcuper run failed")
	}
}

package main

import (
	"github.com/skyline93/syncbyte-go/agent"
	"github.com/skyline93/syncbyte-go/file"
	// "github.com/skyline93/syncbyte-go/internal/agent"
)

func main() {
	dir := "/home/greene/workspace/syncbyte-go/pkg/cache"
	// filename := "/home/greene/Python-3.11.0.tar.xz"
	// filename := "/home/greene/workspace/syncbyte-go/go.mod"
	mountPoint := "/home/greene/workspace/syncbyte-go/datastore"

	nas := file.NewNasVolume(mountPoint)
	// err := nas.UploadFile(filename)
	// if err != nil {
	// 	panic(err)
	// }

	mgmt := agent.NewBackupManager(nas)
	err := mgmt.Backup(dir)
	if err != nil {
		panic(err)
	}
}

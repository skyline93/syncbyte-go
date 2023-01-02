package main

import (
	"context"

	"github.com/skyline93/syncbyte-go/internal/agent"
	"github.com/skyline93/syncbyte-go/pkg/logging"
)

func main() {
	logger := logging.GetSugaredLogger("backup")

	dir := "/etc"

	mountPoint := "/home/greene/workspace/syncbyte-go/datastore"

	nas := agent.NewNasVolume(mountPoint)

	ctx := context.TODO()
	mgmt := agent.NewBackupManager(nas, ctx)

	logger.Infof("start backup..")
	err := mgmt.Backup(dir)
	if err != nil {
		panic(err)
	}
}

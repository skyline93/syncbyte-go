package main

import "github.com/skyline93/syncbyte-go/internal/agent/backuper"

func main() {
	dir := "/Users/greene/Workspace/syncbyte-go/internal"
	dest := "/Users/greene/Workspace/syncbyte-go/datastore"

	dataStore := backuper.NewLocalFileSystem(dest)
	b := backuper.NewBackuper(true, dataStore)

	if err := b.Backup(dir); err != nil {
		panic(err)
	}
}

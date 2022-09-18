package backend

type LocalDisk struct {
	MountPoint string
}

func NewLocalDisk(mountPoint string) *LocalDisk {
	return &LocalDisk{MountPoint: mountPoint}
}

func (l *LocalDisk) Save(destFile string) (size int64, err error) {
	// TODO
	return
}

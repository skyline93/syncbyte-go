package source

type Source interface {
	BeforeBackup() (err error)
	GetSourcePath() (path string)
}

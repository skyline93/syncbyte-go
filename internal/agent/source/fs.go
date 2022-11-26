package source

type FS struct {
	Path string
}

func (f *FS) GetSourcePath() (path string) {
	return f.Path
}

func (f *FS) BeforeBackup() (err error) {

	return nil
}

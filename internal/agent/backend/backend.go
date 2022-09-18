package backend

type Backend interface {
	Put(destFile string) (size int64, err error)
	Get(destFile string) (err error)
}

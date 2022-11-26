package storageunit

import "io"

type StorageUnit interface {
	Save(fileName string, read io.Reader) error
}

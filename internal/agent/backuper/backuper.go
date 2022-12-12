package backuper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func genUuid() string {
	u := uuid.New()
	return u.String()
}

// var HashType = sha256.New()
var HashType = md5.New()

type Block struct {
	io.WriteCloser
	hash hash.Hash
}

func NewWriteBlock(writer io.WriteCloser) *Block {
	return &Block{
		WriteCloser: writer,
		hash:        HashType,
	}
}

func (w *Block) Write(p []byte) (n int, err error) {
	n, err = w.WriteCloser.Write(p)
	if n > 0 {
		w.hash.Write(p)
	}
	return n, err
}

func (w *Block) Hash() string {
	b := w.hash.Sum(nil)
	return hex.EncodeToString(b)
}

type DataStore interface {
	NewWriteBlock() (*Block, error)
}

type LocalFileSystem struct {
	MountPoint string
}

func NewLocalFileSystem(mountPoint string) *LocalFileSystem {
	return &LocalFileSystem{
		MountPoint: mountPoint,
	}
}

func (l *LocalFileSystem) NewWriteBlock() (*Block, error) {
	_, err := os.Stat(l.MountPoint)
	if os.IsNotExist(err) {
		os.MkdirAll(l.MountPoint, 0766)
	}

	p := filepath.Join(l.MountPoint, genUuid())

	fs, err := os.Create(p)
	if err != nil {
		return nil, err
	}

	return NewWriteBlock(fs), nil
}

type Backuper struct {
	dataStore  DataStore
	IsCompress bool
}

func NewBackuper(isCompress bool, dataStore DataStore) *Backuper {
	return &Backuper{
		dataStore:  dataStore,
		IsCompress: isCompress,
	}
}

func (b *Backuper) Backup(path string) (err error) {
	bt := time.Now().Format("20060102150405")

	if b.IsCompress {
		path, err = Compress(path, bt)
		if err != nil {
			return err
		}

		defer os.Remove(path)
	}

	fp, err := os.Open(path)
	if err != nil {
		return err
	}

	blk, err := b.dataStore.NewWriteBlock()
	if err != nil {
		return err
	}

	_, err = io.Copy(blk, fp)
	if err != nil {
		blk.Close()
		return err
	}

	blk.Close()

	fmt.Printf("md5: %s\n", blk.Hash())

	return nil
}

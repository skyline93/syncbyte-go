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
	"github.com/skyline93/syncbyte-go/pkg/cache"
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

func (b *Backuper) backup(path string) (err error) {
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

type Jobs struct {
	cache *cache.Cache
}

type JobDetail struct {
	Status string
	Type   string
}

func (j *Jobs) Add(jd *JobDetail) (id string) {
	return j.cache.SetDefaultWithUuidKey(&jd)
}

func (j *Jobs) Get(id string) *JobDetail {
	result := j.cache.Get(id)

	jd, ok := result.(JobDetail)
	if !ok {
		return nil
	}

	return &jd
}

func (j *Jobs) Update(id string, jd *JobDetail) {
	jobs.cache.SetDefault(id, jd)
}

func NewJobs() *Jobs {
	return &Jobs{
		cache: cache.New(1024, cache.DefaultDuration*60*30, time.Second*60*30),
	}
}

var jobs *Jobs = NewJobs()

func (b *Backuper) Backup(path string) (err error) {
	jobID := jobs.Add(&JobDetail{Status: "running", Type: "backup"})

	if err = b.backup(path); err != nil {
		jobs.Update(jobID, &JobDetail{Status: "failed", Type: "backup"})
		return err
	}

	jobs.Update(jobID, &JobDetail{Status: "succeeded", Type: "backup"})

	return nil
}

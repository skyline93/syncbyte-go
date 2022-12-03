package lz4

import (
	"fmt"
	"io"
	"os"

	"code.cloudfoundry.org/bytefmt"
	"github.com/pierrec/lz4/v4"
)

const (
	lz4Extension   = ".lz4"
	blockMaxSize   = "4M"  //block max size [64K,256K,1M,4M]
	blockChecksum  = false //enable block checksum
	streamChecksum = false //disable stream checksum
	level          = 0     //compression level (0=fastest)
	concurrency    = -1    //concurrency (default=all CPUs)
)

func Compress(filePath string) (string, error) {
	sz, err := bytefmt.ToBytes(blockMaxSize)
	if err != nil {
		return "", err
	}

	zw := lz4.NewWriter(nil)
	options := []lz4.Option{
		lz4.BlockChecksumOption(blockChecksum),
		lz4.BlockSizeOption(lz4.BlockSize(sz)),
		lz4.ChecksumOption(streamChecksum),
		lz4.CompressionLevelOption(lz4.CompressionLevel(level)),
		lz4.ConcurrencyOption(concurrency),
	}
	if err := zw.Apply(options...); err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	finfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	mode := finfo.Mode() // use the same mode for the output file

	zfilename := fmt.Sprintf("%s%s", filePath, lz4Extension)
	zfile, err := os.OpenFile(zfilename, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return "", err
	}
	zw.Reset(zfile)

	_, err = io.Copy(zw, file)
	if err != nil {
		return "", err
	}
	for _, c := range []io.Closer{zw, zfile} {
		err := c.Close()
		if err != nil {
			return "", err
		}
	}

	return zfilename, nil
}

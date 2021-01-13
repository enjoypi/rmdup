package cmd

import (
	"crypto/sha256"
	"hash/adler32"
	"io"
	"os"

	"go.uber.org/zap"
)

const sizeForAdler32 = 4096

type fileHash struct {
	absPath string
	info    os.FileInfo
	adler32 uint32 // to compare fast
	sha256  []byte
}

func (f *fileHash) SumAdler32() {
	if f.adler32 > 0 {
		return
	}

	file, err := os.Open(f.absPath)
	if err != nil {
		logger.Error("open file", zap.Error(err))
		return
	}
	defer file.Close()

	var data []byte
	if f.info.Size() > sizeForAdler32 {
		data = make([]byte, sizeForAdler32)
	} else {
		data = make([]byte, f.info.Size())
	}

	_, err = file.Read(data)
	if err != nil {
		logger.Error("read file", zap.Error(err))
		return
	}

	f.adler32 = adler32.Checksum(data)
}

func (f *fileHash) SumSHA256() {
	file, err := os.Open(f.absPath)
	if err != nil {
		logger.Error("open file", zap.Error(err))
		return
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		logger.Error("read file", zap.Error(err))
		return
	}

	f.sha256 = h.Sum(nil)
}

func (f *fileHash) Same(other *fileHash) bool {
	if f.info.Size() != other.info.Size() {
		return false
	}

	f.SumAdler32()
	other.SumAdler32()

	if f.adler32 != other.adler32 {
		return false
	}

	f.SumSHA256()
	other.SumSHA256()
	if len(f.sha256) != len(other.sha256) {
		return false
	}

	for i := 0; i < len(f.sha256); i++ {
		if f.sha256[i] != other.sha256[i] {
			return false
		}
	}
	return true
}

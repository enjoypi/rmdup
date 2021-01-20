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

func (f *fileHash) PHash() {
	//file1, _ := os.Open("sample1.jpg")
	//file2, _ := os.Open("sample2.jpg")
	//defer file1.Close()
	//defer file2.Close()
	//
	//img1, _ := jpeg.Decode(file1)
	//img2, _ := jpeg.Decode(file2)
	//hash1, _ := goimagehash.AverageHash(img1)
	//hash2, _ := goimagehash.AverageHash(img2)
	//distance, _ := hash1.Distance(hash2)
	//fmt.Printf("Distance between images: %v\n", distance)
	//
	//hash1, _ = goimagehash.DifferenceHash(img1)
	//hash2, _ = goimagehash.DifferenceHash(img2)
	//distance, _ = hash1.Distance(hash2)
	//fmt.Printf("Distance between images: %v\n", distance)
	//width, height := 8, 8
	//hash3, _ = goimagehash.ExtAverageHash(img1, width, height)
	//hash4, _ = goimagehash.ExtAverageHash(img2, width, height)
	//distance, _ = hash3.Distance(hash4)
	//fmt.Printf("Distance between images: %v\n", distance)
	//fmt.Printf("hash3 bit size: %v\n", hash3.Bits())
	//fmt.Printf("hash4 bit size: %v\n", hash4.Bits())
	//
	//var b bytes.Buffer
	//foo := bufio.NewWriter(&b)
	//_ = hash4.Dump(foo)
	//foo.Flush()
	//bar := bufio.NewReader(&b)
	//hash5, _ := goimagehash.LoadExtImageHash(bar)
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

package cmd

import (
	"crypto/sha256"
	"hash/adler32"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/corona10/goimagehash"
	"go.uber.org/zap"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const sizeForAdler32 = 4096

type fileHash struct {
	absPath   string
	info      os.FileInfo
	adler32   uint32 // to compare fast
	sha256    []byte
	imageHash *goimagehash.ImageHash
	mw        *imagick.MagickWand
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

func (f *fileHash) MagickWand() {
	if f.mw != nil {
		return
	}

	mw := imagick.NewMagickWand()
	if err := mw.ReadImage(f.absPath); err != nil {
		logger.Debug(err.Error())
		return
	}
	f.mw = mw
}

func (f *fileHash) PHash() error {
	if f.imageHash != nil {
		return nil
	}

	file, err := os.Open(f.absPath)
	if err != nil {
		logger.Error("open file", zap.Error(err))
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		logger.Debug("image.Decode", zap.String("file", f.absPath), zap.Error(err))
		return err
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		logger.Debug("goimagehash.PerceptionHash", zap.String("file", f.absPath), zap.Error(err))
		return err
	}

	f.imageHash = hash

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

	//var b bytes.Buffer
	//foo := bufio.NewWriter(&b)
	//_ = hash4.Dump(foo)
	//foo.Flush()
	//bar := bufio.NewReader(&b)
	//hash5, _ := goimagehash.LoadExtImageHash(bar)
	return nil
}

func (f *fileHash) Same(other *fileHash) bool {

	if f.info.Size() != other.info.Size() {
		return false
	}

	f.SumAdler32()
	other.SumAdler32()

	if f.adler32 != other.adler32 {
		return f.SameImage(other)
	}

	f.SumSHA256()
	other.SumSHA256()
	if len(f.sha256) != len(other.sha256) {
		return f.SameImage(other)
	}

	for i := 0; i < len(f.sha256); i++ {
		if f.sha256[i] != other.sha256[i] {
			return f.SameImage(other)
		}
	}
	return true
}

func (f *fileHash) SameImage(other *fileHash) bool {
	//f.MagickWand()
	//other.MagickWand()
	//
	//if f.mw != nil && other.mw != nil {
	//	_, distortion := f.mw.CompareImages(other.mw, imagick.METRIC_PERCEPTUAL_HASH_ERROR)
	//	logger.Info("different image", zap.String("file1", f.absPath), zap.String("file2", other.absPath), zap.Float64("distortion", distortion))
	//}

	if f.PHash() == nil && other.PHash() == nil {
		if dis, err := f.imageHash.Distance(other.imageHash); err == nil {
			//if dis != 0 {
			//	logger.Info("different image", zap.String("file1", f.absPath), zap.String("file2", other.absPath), zap.Int("distance", dis))
			//}
			return dis == 0
		}
	}
	return false
}

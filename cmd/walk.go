package cmd

import (
	"os"
	"path/filepath"
)

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || info.Size() <= 0 {
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	hashes, ok := files[info.Size()]
	if !ok {
		hashes = make([]*fileHash, 0)
	}

	hashes = append(hashes, &fileHash{absPath: absPath, info: info})
	files[info.Size()] = hashes

	return nil
}

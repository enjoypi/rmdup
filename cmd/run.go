package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var files = make(map[int64][]*fileHash)

func run(cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		if err := filepath.Walk(arg, walk); err != nil {
			logger.Error("invalid directory", zap.Error(err), zap.String("directory", arg))
		}
	}

	keys := make([]int, 0)
	for k, hashes := range files {
		if len(hashes) > 1 {
			keys = append(keys, int(k))
		}
	}

	sort.Ints(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		hashes := files[int64(keys[i])]
		first := hashes[0]
		same := false
		for j := 1; j < len(hashes); j++ {
			other := hashes[j]
			if first.Same(hashes[j]) {
				if !same {
					fmt.Printf("#rm %s\n", first.absPath)
					same = true
				}
				fmt.Printf("#rm %s\n", other.absPath)
			}
		}
		if same {
			fmt.Println()
		}
	}

	return nil
}

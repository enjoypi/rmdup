package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"go.uber.org/zap"

	"github.com/enjoypi/gojob"
	"github.com/spf13/cobra"
)

var files = make(map[int64][]*fileHash)

func run(cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		if err := filepath.Walk(arg, walk); err != nil {
			return err
		}
	}

	keys := make([]int, 0)
	for k, hashes := range files {
		if len(hashes) > 1 {
			keys = append(keys, int(k))
		}
	}
	sort.Ints(keys)

	m := gojob.NewManager(8)
	for i := len(keys) - 1; i >= 0; i-- {
		values := context.WithValue(m.Context, "size", int64(keys[i]))
		m.Go(func(ctx context.Context, id gojob.TaskID) error {
			size := ctx.Value("size").(int64)
			logger.Debug("started", zap.Int32("taskID", id), zap.Int64("size", size))

			hashes := files[size]
			//logDupFiles(hashes)

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
			return nil
		}, values, nil)
	}

	m.Wait()

	return nil
}

func logDupFiles(hashes []*fileHash) {
	fmt.Print("#")
	for _, v := range hashes {
		fmt.Print(v.absPath, " ")
	}
	fmt.Println()
}

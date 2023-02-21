package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var (
	humanReadable = flag.Bool("h", false, "print sizes in human-readable format (e.g., 1K 234M 2G)")
	blockSize     = flag.Int64("b", 1, "block size in bytes")
	summarize     = flag.Bool("s", false, "display only a total for each argument")
)

func main() {
	flag.Parse()

	var roots []string
	if flag.NArg() == 0 {
		roots = []string{"."}
	} else {
		roots = flag.Args()
	}

	var totalSize int64
	var totalInodes int64
	var wg sync.WaitGroup

	for _, root := range roots {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return nil
			}
			if !info.IsDir() {
				atomic.AddInt64(&totalSize, info.Size())
				atomic.AddInt64(&totalInodes, 1)
				return nil
			}
			if *summarize {
				// only compute the size of the directory itself
				if path != root {
					return filepath.SkipDir
				}
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				size, inodes := getDirInfo(path)
				atomic.AddInt64(&totalSize, size)
				atomic.AddInt64(&totalInodes, inodes)
				if !*summarize {
					fmt.Printf("%d\t%s\n", size/(*blockSize), path)
				}
			}()
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	wg.Wait()

	if *summarize {
		fmt.Printf("%d\t%s\n", totalSize/(*blockSize), "total")
	} else {
		if *humanReadable {
			fmt.Printf("total size: %s\n", humanize(totalSize))
		} else {
			fmt.Printf("total size: %d\n", totalSize/(*blockSize))
		}
		fmt.Printf("total inodes: %d\n", totalInodes)
	}
}

func getDirInfo(path string) (int64, int64) {
	var size int64
	var inodes int64

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return size, inodes
	}

	for _, file := range files {
		if !file.IsDir() {
			size += file.Size()
			inodes++
			continue
		}
		subdir := filepath.Join(path, file.Name())
		subdirSize, subdirInodes := getDirInfo(subdir)
		size += subdirSize
		inodes += subdirInodes
	}

	inodes++ // count directory itself

	return size, inodes
}

func humanize(n int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
	)
	switch {
	case n >= TB:
		return fmt.Sprintf("%.2fT", float64(n)/TB)
	case n >= GB:
		return fmt.Sprintf("%.2fG", float64(n)/GB)
	case n >= MB:
		return fmt.Sprintf("%.2fM", float64(n)/MB)
	case n >= KB:
		return fmt.Sprintf("%.2fK", float64(n)/KB)
	default:
		return fmt.Sprintf("%dB", n)
	}
}

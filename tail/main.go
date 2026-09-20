package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

type TailOptions struct {
	Lines     int
	BlockSize int
}

func main() {
	f, err := os.Open("./test/test.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "file open failed", err)
		os.Exit(1)
	}
	defer f.Close()

	options := TailOptions{
		Lines:     10,
		BlockSize: 1024,
	}
	fInfo, err := f.Stat()
	fSize := fInfo.Size()
	buf := make([]byte, options.BlockSize)
	start := math.Max(0, float64(fSize)-float64(len(buf)))

	// ファイル内容の読み込み位置を指定
	_, err = f.Seek(int64(start), io.SeekStart)

	rbyte, err := f.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, "read failed:", err)
		os.Exit(1)
	}
	// 逆から走査する
	for index := rbyte - 1; index >= 0; index-- {
		if buf[index] == '\n' {
			options.Lines--
		}
		if options.Lines < 0 {
			off := int64(start) + int64(index) + 1
			_, err = f.Seek(off, io.SeekStart)
			if err != nil {
				fmt.Fprintln(os.Stderr, "set seek failed:", err)
				os.Exit(1)
			}
			_, err = io.Copy(os.Stdout, f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "copy std output failed:", err)
				os.Exit(1)
			}
			return
		}
	}
}

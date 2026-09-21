package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

type TailOptions struct {
	readLines int
	BlockSize int
}

func main() {
	options := TailOptions{
		readLines: 10,
		BlockSize: 1024,
	}
	args := os.Args
	files := args[1:]
	for _, fileName := range files {
		fmt.Printf("==> %s <==\n", fileName)
		targetPath, err := filepath.Abs(fileName)
		err = checkError(os.Stderr, "could not resolve file path:", err)
		if err != nil {
			continue
		}

		f, err := os.Open(targetPath)
		err = checkError(os.Stderr, "file open failed:", err)
		if err != nil {
			continue
		}

		linesRemain := options.readLines
		fInfo, err := f.Stat()
		err = checkError(os.Stderr, "get fileInfo failed:", err)
		if err != nil {
			continue
		}
		fSize := fInfo.Size()
		buf := make([]byte, options.BlockSize)
		start := computeStart(fSize, buf)

		// ファイル内容の読み込み位置を指定
		_, err = f.Seek(int64(start), io.SeekStart)

		rbyte, err := f.Read(buf)
		if err != nil && err != io.EOF {
			err = checkError(os.Stderr, "read failed:", err)
			continue
		}
		// 逆から走査する
		for index := rbyte - 1; index >= 0; index-- {
			if buf[index] == '\n' {
				linesRemain--
			}
			if linesRemain < 0 {
				off := int64(start) + int64(index) + 1
				_, err = f.Seek(off, io.SeekStart)
				err = checkError(os.Stderr, "set seek failed:", err)
				if err != nil {
					break
				}
				_, err = io.Copy(os.Stdout, f)
				err = checkError(os.Stderr, "copy std output failed:", err)
				if err != nil {
					break
				}
				break
			}
		}
		if linesRemain >= 0 {
			_, err = f.Seek(0, io.SeekStart)
			err = checkError(os.Stderr, "set seek failed:", err)
			if err != nil {
				break
			}
			_, err = io.Copy(os.Stdout, f)
			err = checkError(os.Stderr, "copy std output failed:", err)
			if err != nil {
				break
			}
		}
		f.Close()
	}
}

func computeStart(size int64, buff []byte) int64 {
	start := math.Max(0, float64(size)-float64(len(buff)))
	return int64(start)
}

func checkError(w io.Writer, msg string, err error) error {
	if err != nil {
		fmt.Fprintln(w, msg, err)
	}

	return err
}

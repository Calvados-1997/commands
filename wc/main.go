package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "wc:", e)
	}
}

func main() {
	f, err := os.Open("test.txt")
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	fmt.Println(scanner.Text())
	// buffer := make([]byte, 4096)
	// count1, err := f.Read(buffer)
	// check(err)

	// TODO:
	// 改行文字をカウント
	// 文字数を読み込むたびにカウント
	// 読み込んだ総バイト数をカウント
	// fmt.Printf("%d bytes: %s\n", count1, string(buffer[:count1]))
	// count2, err := f.Read(buffer)
	// fmt.Printf("%d bytes: %s\n", count2, string(buffer[:count2]))
	// check(err)
}

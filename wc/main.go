package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode"
)

func main() {
	args := os.Args[1:]
	files := slices.DeleteFunc(args, func(f string) bool {
		return strings.Contains(f, "-")
	})

	lines, words, bytes := 0, 0, 0
	buff := make([]byte, 4096)
	for _, fileName := range files {
		inWord := false

		f, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc: ", err)
			os.Exit(1)
		}

		for {
			bytesRead, err := f.Read(buff)
			if bytesRead > 0 {
				lines += countLines(buff[:bytesRead])
				wordCount, isInWord := countWords(buff[:bytesRead], inWord)
				inWord = isInWord
				words += wordCount
				bytes += bytesRead
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "wc: ", err)
				os.Exit(1)
			}
		}
		fmt.Printf("%d %d %d %s\n", lines, words, bytes, fileName)
		f.Close()
	}

	if len(files) > 1 {
		fmt.Printf("%d %d %d total\n", lines, words, bytes)
	}
}

func countLines(data []byte) int {
	words := string(data)
	result := strings.Count(words, "\n")
	return result
}

func countWords(data []byte, inWord bool) (int, bool) {
	wordCount := 0
	for _, b := range data {
		letter := rune(b)
		if unicode.IsSpace(letter) {
			inWord = false
			continue
		}
		if !inWord {
			wordCount++
			inWord = true
		}
	}

	return wordCount, inWord
}

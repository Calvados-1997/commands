package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	countBytesOnly    bool
	countNewLinesOnly bool
	countCharsOnly    bool
	countWordsOnly    bool
)

func init() {
	flag.BoolVar(&countBytesOnly, "c", false, "Write to the standard output the number of bytes in each input file.")
	flag.BoolVar(&countNewLinesOnly, "l", false, "Write to the standard output the number of <newline> characters in each input file.")
	flag.BoolVar(&countCharsOnly, "m", false, "Write to the standard output the number of characters in each input file.")
	flag.BoolVar(&countWordsOnly, "w", false, "Write to the standard output the number of words in each input file.")
	flag.Parse()
}

func main() {
	files := flag.Args()

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
				switch {
				case countBytesOnly:
					bytes += bytesRead
				case countNewLinesOnly:
					lines += countLines(buff[:bytesRead])
				case countCharsOnly:
					chars += countCharacters(buff[:bytesRead])
				case countWordsOnly:
					wordCount, isInWord := countWords(buff[:bytesRead], inWord)
					inWord = isInWord
					words += wordCount
				default:
					lines += countLines(buff[:bytesRead])
					wordCount, isInWord := countWords(buff[:bytesRead], inWord)
					inWord = isInWord
					words += wordCount
					bytes += bytesRead
				}
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

func countCharacters(data []byte) int {
	words := string(data)
	return len([]rune(words))
}

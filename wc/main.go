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

	totalLines, totalWords, totalBytes, totalChars := 0, 0, 0, 0
	buff := make([]byte, 4096)
	for _, fileName := range files {
		lines, words, bytes, chars := 0, 0, 0, 0

		f, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc: ", err)
			continue
		}

		inWord := false

		for {
			bytesRead, err := f.Read(buff)
			if bytesRead > 0 {
				lines += countLines(buff[:bytesRead])
				wordCount, isInWord := countWords(buff[:bytesRead], inWord)
				inWord = isInWord
				words += wordCount
				chars += countCharacters(buff[:bytesRead])
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
		fmt.Println(formatResult(lines, words, bytes, chars, fileName))
		f.Close()

		totalLines += lines
		totalWords += words
		totalBytes += bytes
		totalChars += chars
	}

	if len(files) > 1 {
		fmt.Println(formatResult(totalLines, totalWords, totalBytes, totalChars, "total"))
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

func formatResult(lines, words, bytes, chars int, lastInfo string) string {
	cols := []string{}

	noFlags := !countNewLinesOnly && !countWordsOnly && !countBytesOnly && !countCharsOnly

	if noFlags || countNewLinesOnly {
		cols = append(cols, fmt.Sprintf("%d", lines))
	}
	if noFlags || countWordsOnly {
		cols = append(cols, fmt.Sprintf("%d", words))
	}
	if noFlags || countBytesOnly {
		cols = append(cols, fmt.Sprintf("%d", bytes))
	}
	if countCharsOnly {
		cols = append(cols, fmt.Sprintf("%d", chars))
	}

	cols = append(cols, lastInfo)
	return strings.Join(cols, " ")
}

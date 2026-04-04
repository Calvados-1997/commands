package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type counts struct {
	lines int
	words int
	bytes int
	chars int
}

func (c *counts) add(other counts) {
	c.lines += other.lines
	c.words += other.words
	c.bytes += other.bytes
	c.chars += other.chars
}

type options struct {
	showLines bool
	showWords bool
	showBytes bool
	showChars bool
}

func (o options) isDefault() bool {
	return !o.showLines && !o.showWords && !o.showBytes && !o.showChars
}

func main() {
	opts := parseFlags()
	files := flag.Args()

	var total counts
	for _, fileName := range files {
		c, err := processFile(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			continue
		}
		fmt.Println(formatResult(c, opts, fileName))
		total.add(c)
	}

	if len(files) > 1 {
		fmt.Println(formatResult(total, opts, "total"))
	}
}

func parseFlags() options {
	var opts options
	flag.BoolVar(&opts.showBytes, "c", false, "Write to the standard output the number of bytes in each input file.")
	flag.BoolVar(&opts.showLines, "l", false, "Write to the standard output the number of <newline> characters in each input file.")
	flag.BoolVar(&opts.showChars, "m", false, "Write to the standard output the number of characters in each input file.")
	flag.BoolVar(&opts.showWords, "w", false, "Write to the standard output the number of words in each input file.")
	flag.Parse()
	return opts
}

func processFile(fileName string) (counts, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return counts{}, err
	}
	defer f.Close()

	var c counts
	var inWord bool
	buf := make([]byte, 4096)

	for {
		n, err := f.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			c.lines += countLines(chunk)
			wc, iw := countWords(chunk, inWord)
			inWord = iw
			c.words += wc
			c.chars += countChars(chunk)
			c.bytes += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func countLines(data []byte) int {
	return strings.Count(string(data), "\n")
}

func countWords(data []byte, inWord bool) (int, bool) {
	n := 0
	for _, b := range data {
		if unicode.IsSpace(rune(b)) {
			inWord = false
			continue
		}
		if !inWord {
			n++
			inWord = true
		}
	}
	return n, inWord
}

func countChars(data []byte) int {
	return len([]rune(string(data)))
}

func formatResult(c counts, opts options, name string) string {
	var cols []string

	if opts.isDefault() || opts.showLines {
		cols = append(cols, fmt.Sprintf("%7d", c.lines))
	}
	if opts.isDefault() || opts.showWords {
		cols = append(cols, fmt.Sprintf("%7d", c.words))
	}
	if opts.isDefault() || opts.showBytes {
		cols = append(cols, fmt.Sprintf("%7d", c.bytes))
	}
	if opts.showChars {
		cols = append(cols, fmt.Sprintf("%7d", c.chars))
	}

	cols = append(cols, name)
	return strings.Join(cols, " ")
}

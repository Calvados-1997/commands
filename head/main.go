package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Flags struct {
	c int
	n int
}

func parseFlags(args []string) (Flags, []string, int, error) {
	var flagOptions Flags
	fs := flag.NewFlagSet("head", flag.ContinueOnError)
	fs.IntVar(&flagOptions.c, "c", 0, "option specify bytes.")
	fs.IntVar(&flagOptions.n, "n", 10, "option specify lines.")
	err := fs.Parse(args)
	return flagOptions, fs.Args(), fs.NFlag(), err
}

func main() {
	flags, args, nFlag, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "head: ", err)
		os.Exit(1)
	}
	if nFlag >= 2 {
		fmt.Fprintln(os.Stderr, "head: cannot combine line and bytes")
		os.Exit(1)
	}

	fileNameList := args

	for _, fileName := range fileNameList {
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "head: ", err)
			os.Exit(1)
		}

		err = processFile(f, flags)
		if err != nil {
			fmt.Fprintln(os.Stderr, "head: ", err)
			os.Exit(1)
		}
		f.Close()
	}
}

func processFile(rd io.Reader, flags Flags) error {
	if flags.c != 0 {
		reader := bufio.NewReader(rd)
		result, err := readBytes(reader, flags.c)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	}

	scanner := bufio.NewScanner(rd)
	resultLines := readLines(scanner, flags.n)
	for _, line := range resultLines {
		fmt.Println(line)
	}

	return nil
}

func readBytes(reader io.Reader, bytes int) (string, error) {
	rdBytes := make([]byte, bytes)
	data, err := io.ReadFull(reader, rdBytes)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", err
	}

	return string(rdBytes[:data]), nil
}

func readLines(scanner *bufio.Scanner, lineLimit int) []string {
	result := []string{}
	for i := 0; i < lineLimit; i++ {
		if !scanner.Scan() {
			break
		}
		result = append(result, scanner.Text())
	}
	return result
}

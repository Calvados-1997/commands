package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	for _, arg := range args {
		ctn, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: ", err)
			os.Exit(1)
		}
		os.Stdout.Write(ctn)
	}
}

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if _, err := fmt.Println(strings.Join(args, " ")); err != nil {
		fmt.Fprintln(os.Stderr, "echo:", err)
		os.Exit(1)
	}
}

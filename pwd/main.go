package main

import (
	"fmt"
	"os"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		os.Exit(1)
	}
	fmt.Println(currDir)
}

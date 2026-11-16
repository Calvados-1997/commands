package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	/*
		カレントディレクトリにあるファイルやディレクトリ名を出力する。
	*/
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	files, err := os.ReadDir(currDir)
	if err != nil {
		fmt.Println(err)
	}
	res := []string{}
	for _, file := range files {
		res = append(res, file.Name())
	}
	fmt.Print(strings.Join(res, " "))
}

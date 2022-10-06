package main

import (
	"edholm.dev/repofinder"
	"fmt"
	"os"
)

func main() {
	pathsToSearch := os.Args[1:]
	if len(pathsToSearch) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: failed to get current working directory: %v\n", err)
			os.Exit(1)
		}
		pathsToSearch = append(pathsToSearch, currentDir)
	}
	if err := repofinder.Search(pathsToSearch); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

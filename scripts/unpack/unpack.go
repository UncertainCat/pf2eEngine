package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"pf2eEngine/scripts"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: unpack <packed-file>")
		return
	}

	packedFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening packed file:", err)
		return
	}
	defer packedFile.Close()

	scanner := bufio.NewScanner(packedFile)
	var currentFile *os.File
	var currentPath string

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, scripts.StartDelimiter):
			currentPath = extractPath(line, scripts.StartDelimiter)
			if err := createFile(currentPath); err != nil {
				fmt.Printf("Error creating %s: %v\n", currentPath, err)
				currentPath = ""
				continue
			}
			currentFile, err = os.Create(currentPath)
			if err != nil {
				fmt.Printf("Error opening %s: %v\n", currentPath, err)
				currentPath = ""
			}

		case strings.HasPrefix(line, scripts.EndDelimiter):
			if currentFile != nil {
				currentFile.Close()
				currentFile = nil
			}
			currentPath = ""

		case currentFile != nil:
			_, err := fmt.Fprintln(currentFile, line)
			if err != nil {
				fmt.Printf("Error writing to %s: %v\n", currentPath, err)
				currentFile.Close()
				currentPath = ""
			}
		}
	}

	if currentFile != nil {
		currentFile.Close()
	}
}

func extractPath(line, prefix string) string {
	return strings.TrimPrefix(line, prefix)
}

func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

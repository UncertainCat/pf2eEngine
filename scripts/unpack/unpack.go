package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

		// Check if the line indicates the start of a new file
		if strings.HasPrefix(line, "====== FILE-START: ") {
			// Extract file path
			currentPath = extractPath(line, "====== FILE-START: ", " ======")
			if currentPath == "" {
				fmt.Println("Error: Could not extract file path from:", line)
				continue
			}

			// Ensure directories exist
			if err := createFile(currentPath); err != nil {
				fmt.Printf("Error creating %s: %v\n", currentPath, err)
				currentPath = ""
				continue
			}

			// Open file for writing
			currentFile, err = os.Create(currentPath)
			if err != nil {
				fmt.Printf("Error opening %s: %v\n", currentPath, err)
				currentPath = ""
			}
			continue
		}

		// Check for end of a file
		if strings.HasPrefix(line, "====== FILE-END: ") {
			if currentFile != nil {
				currentFile.Close()
				currentFile = nil
			}
			currentPath = ""
			continue
		}

		// Write contents to the current file
		if currentFile != nil {
			_, err := currentFile.WriteString(line + "\n")
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

	fmt.Println("Unpacking completed successfully.")
}

// Extracts the file path from the delimiter line
func extractPath(line, start, end string) string {
	line = strings.TrimPrefix(line, start)
	line = strings.TrimSuffix(line, end)
	return strings.TrimSpace(line)
}

// Ensures the directory for the file exists before creating it
func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

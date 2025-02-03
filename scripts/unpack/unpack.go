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
	// Open packed file
	packedFile, err := os.Open(scripts.OutputFileName)
	if err != nil {
		fmt.Println("❌ Error opening packed file:", err)
		return
	}
	defer packedFile.Close()

	// **Step 1: Remove existing .go files**
	err = clearGoFiles(".")
	if err != nil {
		fmt.Println("❌ Error clearing existing .go files:", err)
		return
	}
	fmt.Println("✅ Cleared old .go files.")

	// **Step 2: Unpacking process**
	scanner := bufio.NewScanner(packedFile)
	var currentFile *os.File
	var currentPath string

	for scanner.Scan() {
		line := scanner.Text()

		// **Detect start of a new file**
		if strings.HasPrefix(line, "====== FILE-START: ") {
			if currentFile != nil {
				currentFile.Close()
			}

			// Extract filename
			currentPath = extractPath(line, scripts.StartDelimiter)
			if currentPath == "" {
				fmt.Println("❌ Error: Could not extract file path from:", line)
				continue
			}

			// Ensure directories exist
			if err := createFile(currentPath); err != nil {
				fmt.Printf("❌ Error creating %s: %v\n", currentPath, err)
				currentPath = ""
				continue
			}

			// Open new file for writing
			currentFile, err = os.Create(currentPath)
			if err != nil {
				fmt.Printf("❌ Error opening %s: %v\n", currentPath, err)
				currentPath = ""
				continue
			}

			fmt.Println("📄 Creating file:", currentPath)
			continue
		}

		// **Detect end of a file**
		if strings.HasPrefix(line, "====== FILE-END: ") {
			if currentFile != nil {
				currentFile.Close()
				fmt.Println("✅ Successfully unpacked:", currentPath)
				currentFile = nil
			}
			currentPath = ""
			continue
		}

		// **Write contents to the current file**
		if currentFile != nil {
			_, err := currentFile.WriteString(line + "\n")
			if err != nil {
				fmt.Printf("❌ Error writing to %s: %v\n", currentPath, err)
				currentFile.Close()
				currentPath = ""
			}
		}
	}

	if currentFile != nil {
		currentFile.Close()
	}

	fmt.Println("🎉 Unpacking completed successfully!")
}

// **Extracts the file path from the delimiter line**
func extractPath(line, delimiter string) string {
	trimmed := strings.TrimPrefix(line, fmt.Sprintf(delimiter, ""))
	return strings.TrimSpace(trimmed)
}

// **Ensures the directory exists before creating the file**
func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

// **Removes all existing .go files before extraction**
func clearGoFiles(rootDir string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fmt.Println("🗑️ Deleting:", path)
			return os.Remove(path)
		}
		return nil
	})
}

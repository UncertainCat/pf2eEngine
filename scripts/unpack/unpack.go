package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"pf2eEngine/scripts"
	"strings"
)

type PackedFile struct {
	Path    string
	Content strings.Builder
}

func main() {
	packedFile, err := os.Open(scripts.OutputFileName)
	if err != nil {
		fmt.Println("❌ Error opening packed file:", err)
		return
	}
	defer packedFile.Close()

	// First pass: Parse all files from packed file
	packedFiles, err := parsePackedFiles(packedFile)
	if err != nil {
		fmt.Println("❌ Error parsing packed files:", err)
		return
	}

	// Delete files that exist but aren't in the packed files
	err = deleteOrphanedGoFiles(".", packedFiles)
	if err != nil {
		fmt.Println("❌ Error cleaning orphaned files:", err)
		return
	}

	// Second pass: Write/update all packed files
	for _, pf := range packedFiles {
		err := writeFile(pf.Path, pf.Content.String())
		if err != nil {
			fmt.Printf("❌ Error writing %s: %v\n", pf.Path, err)
		} else {
			fmt.Printf("📄 Updated/Created: %s\n", pf.Path)
		}
	}

	fmt.Println("🎉 Unpacking completed successfully!")
}

func parsePackedFiles(packedFile *os.File) ([]PackedFile, error) {
	var files []PackedFile
	var currentFile *PackedFile
	scanner := bufio.NewScanner(packedFile)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "====== FILE-START: ") {
			// Close previous file if any
			if currentFile != nil {
				files = append(files, *currentFile)
			}

			path := extractPath(line)
			if path == "" {
				return nil, fmt.Errorf("invalid file start line: %s", line)
			}

			currentFile = &PackedFile{
				Path: path,
			}
			continue
		}

		if strings.HasPrefix(line, "====== FILE-END: ") {
			if currentFile != nil {
				files = append(files, *currentFile)
				currentFile = nil
			}
			continue
		}

		if currentFile != nil {
			currentFile.Content.WriteString(line + "\n")
		}
	}

	if currentFile != nil {
		files = append(files, *currentFile)
	}

	return files, scanner.Err()
}

func deleteOrphanedGoFiles(rootDir string, packedFiles []PackedFile) error {
	// Create set of packed file paths
	packedPaths := make(map[string]struct{})
	for _, pf := range packedFiles {
		packedPaths[pf.Path] = struct{}{}
	}

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}

			if _, exists := packedPaths[relPath]; !exists {
				fmt.Printf("🗑️ Deleting orphaned file: %s\n", relPath)
				return os.Remove(path)
			}
		}
		return nil
	})
}

func writeFile(path string, content string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Check if content matches existing file
	existingContent, err := os.ReadFile(path)
	if err == nil && string(existingContent) == content {
		fmt.Printf("✅ No changes: %s\n", path)
		return nil
	}

	// Write file with new content
	return os.WriteFile(path, []byte(content), 0644)
}

func extractPath(line string) string {
	startMarker := "====== FILE-START: "
	endMarker := " ======"

	if !strings.HasPrefix(line, startMarker) || !strings.HasSuffix(line, endMarker) {
		return ""
	}

	path := strings.TrimPrefix(line, startMarker)
	path = strings.TrimSuffix(path, endMarker)
	return strings.TrimSpace(path)
}

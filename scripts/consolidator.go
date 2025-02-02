package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const outputFileName = "merged_go_files.txt"

func main() {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			if err := appendFileContents(outputFile, path); err != nil {
				fmt.Println("Error appending file contents:", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func appendFileContents(outputFile *os.File, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(outputFile, "\n==== File: %s ====\n\n", filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(outputFile, "\n==== End of", filePath, "====\n")
	return err
}

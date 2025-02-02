package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pf2eEngine/scripts"
	"strings"
)

func main() {
	outputFile, err := os.Create(scripts.OutputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return err
		}

		return processFile(outputFile, path)
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
	}
}

func processFile(outputFile *os.File, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write start delimiter
	if _, err := fmt.Fprintf(outputFile, scripts.StartDelimiter, path); err != nil {
		return err
	}

	// Copy file contents
	if _, err := io.Copy(outputFile, file); err != nil {
		return err
	}

	// Ensure content ends with newline before end delimiter
	if !endsWithNewline(file) {
		if _, err := fmt.Fprintln(outputFile); err != nil {
			return err
		}
	}

	// Write end delimiter
	_, err = fmt.Fprintf(outputFile, scripts.EndDelimiter, path)
	return err
}

func endsWithNewline(file *os.File) bool {
	stat, _ := file.Stat()
	if stat.Size() == 0 {
		return true
	}

	_, _ = file.Seek(-1, io.SeekEnd)
	var buf [1]byte
	_, err := file.Read(buf[:])
	_, err = file.Seek(0, io.SeekStart) // Reset file pointer
	return err == nil && buf[0] == '\n'
}

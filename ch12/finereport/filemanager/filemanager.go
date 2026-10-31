// Package filemanager reads lines from one file and
// writes JSON to another.
package filemanager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// FileManager knows its input and output paths.
type FileManager struct {
	InputPath  string
	OutputPath string
}

// New returns a FileManager for the two paths.
func New(in, out string) FileManager {
	return FileManager{InputPath: in, OutputPath: out}
}

// ReadLines returns every non-empty line of the input.
func (fm FileManager) ReadLines() ([]string, error) {
	file, err := os.Open(fm.InputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	return lines, nil
}

// WriteResult writes data as indented JSON.
func (fm FileManager) WriteResult(data any) error {
	file, err := os.Create(fm.OutputPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}
	defer file.Close()

	time.Sleep(time.Second) // stand-in for a slow disk

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

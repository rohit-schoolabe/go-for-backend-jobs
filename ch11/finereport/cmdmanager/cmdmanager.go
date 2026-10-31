// Package cmdmanager reads from and writes to the terminal.
package cmdmanager

import (
	"encoding/json"
	"fmt"
	"os"
)

// CMDManager talks to the person at the desk.
type CMDManager struct{}

// New returns a CMDManager.
func New() CMDManager { return CMDManager{} }

// ReadLines asks for days late until the user types done.
func (CMDManager) ReadLines() ([]string, error) {
	fmt.Println("Days late, one per line. Type done to finish.")
	var lines []string
	for {
		var entry string
		fmt.Print("Days: ")
		if _, err := fmt.Scanln(&entry); err != nil {
			return nil, err
		}
		if entry == "done" {
			return lines, nil
		}
		lines = append(lines, entry)
	}
}

// WriteResult prints data as indented JSON.
func (CMDManager) WriteResult(data any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

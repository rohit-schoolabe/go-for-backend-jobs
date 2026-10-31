// Package fileops stores small values in text files.
package fileops

import (
	"errors"
	"os"
	"strconv"
)

// WriteInt saves value as decimal text in the named file.
func WriteInt(fileName string, value int) error {
	text := strconv.Itoa(value)
	return os.WriteFile(fileName, []byte(text), 0644)
}

// ReadInt reads a whole number saved by WriteInt.
// If the file is missing or broken it returns defaultValue
// together with an error that says what went wrong.
func ReadInt(fileName string, defaultValue int) (int, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return defaultValue, errors.New("no saved value found")
	}
	value, err := strconv.Atoi(string(data))
	if err != nil {
		return defaultValue, errors.New("value is not a number")
	}
	return value, nil
}

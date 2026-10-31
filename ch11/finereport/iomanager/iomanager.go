// Package iomanager says where a job reads and writes.
package iomanager

// IOManager supplies input lines and stores a result.
type IOManager interface {
	ReadLines() ([]string, error)
	WriteResult(data any) error
}

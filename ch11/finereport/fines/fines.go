// Package fines works out overdue fines for one daily rate.
package fines

import (
	"fmt"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/conversion"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch11/finereport/iomanager"
)

// Job applies one daily rate to a list of days late.
type Job struct {
	IO       iomanager.IOManager `json:"-"`
	Rate     float64             `json:"rate_per_day"`
	DaysLate []int               `json:"days_late"`
	Fines    map[string]string   `json:"fines"`
}

// NewJob prepares a job; data is loaded by Process.
func NewJob(io iomanager.IOManager, rate float64) *Job {
	return &Job{IO: io, Rate: rate}
}

// LoadData reads and parses the days late.
func (job *Job) LoadData() error {
	lines, err := job.IO.ReadLines()
	if err != nil {
		return err
	}
	days, err := conversion.StringsToInts(lines)
	if err != nil {
		return err
	}
	job.DaysLate = days
	return nil
}

// Process loads data, works out fines and writes them.
func (job *Job) Process() error {
	if err := job.LoadData(); err != nil {
		return fmt.Errorf("rate %.0f: %w", job.Rate, err)
	}
	result := make(map[string]string, len(job.DaysLate))
	for _, d := range job.DaysLate {
		fine := min(float64(d)*job.Rate, 50)
		result[fmt.Sprint(d)] = fmt.Sprintf("%.2f", fine)
	}
	job.Fines = result
	return job.IO.WriteResult(job)
}

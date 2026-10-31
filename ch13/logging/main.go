// Command logging shows structured logs with log/slog.
package main

import (
	"errors"
	"log/slog"
	"os"
	"time"
)

func main() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	slog.Info("book issued",
		"code", "CS-101", "member", "Ananya", "days", 14)

	// With returns a logger that adds fields to every line.
	job := logger.With("job", "fine-report")
	job.Debug("starting", "files", 3)
	job.Error("cannot read input",
		"path", "overdue.txt",
		"err", errors.New("permission denied"))

	slog.Info("done", "took", 1500*time.Millisecond)
}

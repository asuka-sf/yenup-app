package usecase

import (
	"context"
	"errors"
	"fmt"
	"yenup/internal/domain/notifier"
	"yenup/internal/domain/storage"
)

type WeeklyReportUsecase interface {
	GenerateReport(ctx context.Context) error
}

type WeeklyReporter struct {
	StorageClient storage.Client
	Notifier      notifier.Notifier
}

// NewWeeklyReporter creates a new WeeklyReporter with the given storage client and notifier.
func NewWeeklyReporter(storageClient storage.Client, notifier notifier.Notifier) *WeeklyReporter {
	return &WeeklyReporter{
		StorageClient: storageClient,
		Notifier:      notifier,
	}
}

// GenerateReport reads rates from GCS, calucurates weekly summery, and notifies via Slack.
func (w *WeeklyReporter) GenerateReport(ctx context.Context) error {
	// read rates
	rates, err := w.StorageClient.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read rates: %w", err)
	}
	// validate rates
	if len(rates) == 0 {
		return errors.New("no rates found")
	}

	var total float64
	dateMap := make(map[string]bool)
	baseBase := rates[0].Base
	baseTarget := rates[0].Target
	max := rates[0].Value
	min := rates[0].Value
	// calcurate max, min, average
	for _, r := range rates {
		if dateMap[r.Date] {
			return errors.New("duplicate dates")
		}
		dateMap[r.Date] = true
		if r.Base != baseBase || r.Target != baseTarget {
			return errors.New("inconsistent base/target")
		}

		if max < r.Value {
			max = r.Value
		}
		if min > r.Value {
			min = r.Value
		}

		total += r.Value
	}
	average := total / float64(len(rates))

	// Notify
	msg := fmt.Sprintf("This week report. Average: %.2f, Max: %.2f, Min: %.2f", average, max, min)

	if err := w.Notifier.Notify(msg); err != nil {
		return fmt.Errorf("failed to notify: %w", err)
	}

	return nil
}

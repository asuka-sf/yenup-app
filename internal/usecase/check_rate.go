package usecase

import (
	"context"
	"fmt"
	"time"

	"yenup/internal/domain/notifier"
	"yenup/internal/domain/rate"
	"yenup/internal/domain/storage"
)

// RateCheckUsecase is the interface for the rate check usecase
type RateCheckUsecase interface {
	CheckRates(ctx context.Context, base, target string, forceNotify bool) (*CheckRateResult, error)
}

// CheckRateResult is the result of checking the rate
type CheckRateResult struct {
	TodayRate     float64
	YesterdayRate float64
	IsNotified    bool
}

// RateChecker is the usecase for checking the rate
type RateChecker struct {
	StorageClient storage.Client
	Fetcher       rate.RateFetcher
	Notifier      notifier.Notifier
}

func NewRateChecker(storageClient storage.Client, fetcher rate.RateFetcher, notifier notifier.Notifier) *RateChecker {
	return &RateChecker{
		StorageClient: storageClient,
		Fetcher:       fetcher,
		Notifier:      notifier,
	}
}

func (r *RateChecker) CheckRates(ctx context.Context, base, target string, forceNotify bool) (*CheckRateResult, error) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	todayStr := today.Format("2006-01-02")
	yesterdayStr := yesterday.Format("2006-01-02")

	// Get rates from repository
	todayRate, err := r.Fetcher.FetchRate(todayStr, base, target)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's rate: %w", err)
	}
	yesterdayRate, err := r.Fetcher.FetchRate(yesterdayStr, base, target)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch yesterday's rate: %w", err)
	}

	// read the JSON file to get rates information
	rates, err := r.StorageClient.Read(ctx)
	if err != nil {
		return nil, err
	}

	// update the JSON file
	if err := r.saveRates(ctx, &todayRate, rates); err != nil {
		return nil, err
	}

	result := &CheckRateResult{
		TodayRate:     todayRate.Value,
		YesterdayRate: yesterdayRate.Value,
		IsNotified:    false,
	}

	shouldNotify := forceNotify || (todayRate.Value < yesterdayRate.Value)
	if !shouldNotify {
		return result, nil
	}

	msg := fmt.Sprintf(
		"JPY Stronger Alert! %s/%s: Yesterday %.4f -> Today %.4f",
		base,
		target,
		yesterdayRate.Value,
		todayRate.Value,
	)
	if forceNotify && !(todayRate.Value < yesterdayRate.Value) {
		msg = fmt.Sprintf(
			"Test Notification (forced). %s/%s: Yesterday %.4f -> Today %.4f",
			base,
			target,
			yesterdayRate.Value,
			todayRate.Value,
		)
	}

	if err := r.Notifier.Notify(msg); err != nil {
		return nil, fmt.Errorf("failed to notify: %w", err)
	}

	result.IsNotified = true
	return result, nil
}

func (r *RateChecker) saveRates(ctx context.Context, rate *rate.Rate, rates []*rate.Rate) error {

	// if the json file has more than 7 days' rates, remove the early days' ones
	if len(rates) >= 7 {
		n := len(rates) - 6
		rates = rates[n:]
	}
	rates = append(rates, rate)

	if err := r.StorageClient.Write(ctx, rates); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}
	return nil
}

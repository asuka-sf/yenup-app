package usecase

import (
	"fmt"
	"time"
	"yenup/internal/domain/rate"
)

// RateCheckUsecase is the interface for the rate check usecase
type RateCheckUsecase interface {
	CheckRates(base, target string, forceNotify bool) (*CheckRateResult, error)
}

// CheckRateResult is the result of checking the rate
type CheckRateResult struct {
	TodayRate     float64
	YesterdayRate float64
	IsNotified    bool
}

// RateChecker is the usecase for checking the rate
type RateChecker struct {
	Fetcher  rate.RateFetcher
	Notifier rate.Notifier
}

func NewRateChecker(repo rate.RateFetcher, notifier rate.Notifier) *RateChecker {
	return &RateChecker{
		Fetcher:  repo,
		Notifier: notifier,
	}
}

func (r *RateChecker) CheckRates(base, target string, forceNotify bool) (*CheckRateResult, error) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	// Get rates from repository
	todayRate, err := r.Fetcher.FetchRate(today, base, target)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's rate: %w", err)
	}
	yesterdayRate, err := r.Fetcher.FetchRate(yesterday, base, target)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch yesterday's rate: %w", err)
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

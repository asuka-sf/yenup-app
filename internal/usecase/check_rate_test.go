package usecase

import (
	"context"
	"errors"
	"testing"

	"yenup/internal/domain/rate"

	"github.com/stretchr/testify/assert"
)

func TestCheckRates(t *testing.T) {
	tests := []struct {
		name             string
		mockRates        []*rate.Rate
		mockFetcher      []rate.Rate
		expected         *CheckRateResult
		wantWrittenRates []*rate.Rate
		mockFetchErr     error
		mockReadErr      error
		mockWriteErr     error
		mockNotifyErr    error
		forceNotify      bool
		wantErr          bool
	}{
		{
			name:        "success: JPY is stronger",
			mockRates:   []*rate.Rate{},
			mockFetcher: []rate.Rate{todayRate, yesterdayRate},
			expected:    &CheckRateResult{TodayRate: todayRate.Value, YesterdayRate: yesterdayRate.Value, IsNotified: true},
		},
		{
			name:      "success: JPY is weaker but forceNotify is true",
			mockRates: []*rate.Rate{},
			mockFetcher: []rate.Rate{
				{Date: "2026-03-19", Base: "CAD", Target: "JPY", Value: 113.20},
				yesterdayRate,
			},
			expected:    &CheckRateResult{TodayRate: 113.20, YesterdayRate: yesterdayRate.Value, IsNotified: true},
			forceNotify: true,
		},
		{
			name:      "success: JPY is weaker, no notification",
			mockRates: []*rate.Rate{},
			mockFetcher: []rate.Rate{
				{Date: "2026-03-19", Base: "CAD", Target: "JPY", Value: 113.20},
				yesterdayRate,
			},
			expected: &CheckRateResult{TodayRate: 113.20, YesterdayRate: yesterdayRate.Value, IsNotified: false},
		},
		{
			name:        "success: return 7 rates",
			mockRates:   testValidRates,
			mockFetcher: []rate.Rate{todayRate, yesterdayRate},
			expected:    &CheckRateResult{TodayRate: todayRate.Value, YesterdayRate: yesterdayRate.Value, IsNotified: true},
		},
		{
			name: "success: replace existing entry for the same date",
			mockRates: []*rate.Rate{
				{Date: "2026-03-18", Base: "CAD", Target: "JPY", Value: 111.00},
				{Date: "2026-03-19", Base: "CAD", Target: "JPY", Value: 999.99},
			},
			mockFetcher: []rate.Rate{todayRate, yesterdayRate},
			expected:    &CheckRateResult{TodayRate: todayRate.Value, YesterdayRate: yesterdayRate.Value, IsNotified: true},
			wantWrittenRates: []*rate.Rate{
				{Date: "2026-03-18", Base: "CAD", Target: "JPY", Value: 111.00},
				{Date: "2026-03-19", Base: "CAD", Target: "JPY", Value: 110.22},
			},
		},
		{
			name:         "error: fail to fetch rate",
			mockRates:    []*rate.Rate{},
			mockFetchErr: errors.New("failed to fetch rate"),
			wantErr:      true,
		},
		{
			name:        "error: fail to load a JSON file",
			mockRates:   []*rate.Rate{},
			mockFetcher: []rate.Rate{todayRate, yesterdayRate},
			mockReadErr: errors.New("failed to load a JSON file"),
			wantErr:     true,
		},
		{
			name:         "error: fail to write a JSON file",
			mockRates:    []*rate.Rate{},
			mockFetcher:  []rate.Rate{todayRate, yesterdayRate},
			mockWriteErr: errors.New("failed to write a JSON file"),
			wantErr:      true,
		},
		{
			name:          "error: fail to notify",
			mockRates:     []*rate.Rate{},
			mockFetcher:   []rate.Rate{todayRate, yesterdayRate},
			mockNotifyErr: errors.New("failed to notify"),
			wantErr:       true,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockStorageClient{
				rates:    tt.mockRates,
				readErr:  tt.mockReadErr,
				writeErr: tt.mockWriteErr,
			}
			fetcher := &MockFetcher{
				rates: tt.mockFetcher,
				err:   tt.mockFetchErr,
			}
			notifier := &MockNotifier{err: tt.mockNotifyErr}
			uc := NewRateChecker(storage, fetcher, notifier)
			result, err := uc.CheckRates(ctx, "CAD", "JPY", tt.forceNotify)

			if tt.wantErr {
				// if an error expected
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.wantWrittenRates != nil {
					assert.Equal(t, tt.wantWrittenRates, storage.writtenRates)
				} else if len(tt.mockRates) >= 7 {
					assert.Len(t, storage.writtenRates, 7)
				} else {
					assert.Len(t, storage.writtenRates, len(tt.mockRates)+1)
				}

				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

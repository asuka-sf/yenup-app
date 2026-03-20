package usecase

import (
	"context"

	"yenup/internal/domain/rate"
)

type MockStorageClient struct {
	rates        []*rate.Rate
	readErr      error
	writeErr     error
	writtenRates []*rate.Rate
}

// ----------------------------------------------------------------------------------------------
// weekly_report_test helpers
// ----------------------------------------------------------------------------------------------

func (m *MockStorageClient) Read(ctx context.Context) ([]*rate.Rate, error) {
	return m.rates, m.readErr
}

func (m *MockStorageClient) Write(ctx context.Context, rates []*rate.Rate) error {
	m.writtenRates = rates
	return m.writeErr
}

type MockNotifier struct {
	msg string
	err error
}

func (m *MockNotifier) Notify(message string) error {
	m.msg = message
	return m.err
}

var testValidRates = []*rate.Rate{
	{Date: "2026-01-01", Base: "CAD", Target: "JPY", Value: 113.2207},
	{Date: "2026-01-02", Base: "CAD", Target: "JPY", Value: 112.5783},
	{Date: "2026-01-03", Base: "CAD", Target: "JPY", Value: 113.3475},
	{Date: "2026-01-04", Base: "CAD", Target: "JPY", Value: 110.4845},
	{Date: "2026-01-05", Base: "CAD", Target: "JPY", Value: 115.0950},
	{Date: "2026-01-06", Base: "CAD", Target: "JPY", Value: 114.2859},
	{Date: "2026-01-07", Base: "CAD", Target: "JPY", Value: 111.6880},
}

var testDuplicatedDate = []*rate.Rate{
	{Date: "2026-01-01", Base: "CAD", Target: "JPY", Value: 113.2207},
	{Date: "2026-01-07", Base: "CAD", Target: "JPY", Value: 111.6880},
	{Date: "2026-01-07", Base: "CAD", Target: "JPY", Value: 110.6090},
}
var testInconsistentBase = []*rate.Rate{
	{Date: "2026-01-01", Base: "CAD", Target: "JPY", Value: 113.2207},
	{Date: "2026-01-08", Base: "USD", Target: "JPY", Value: 113.8581},
}
var testInconsistentTarget = []*rate.Rate{
	{Date: "2026-01-01", Base: "CAD", Target: "JPY", Value: 113.2207},
	{Date: "2026-01-09", Base: "CAD", Target: "USD", Value: 112.0795},
}

// ----------------------------------------------------------------------------------------------
// check_rate_test helpers
// ----------------------------------------------------------------------------------------------

var todayRate = rate.Rate{Date: "2026-03-19", Base: "CAD", Target: "JPY", Value: 110.22}
var yesterdayRate = rate.Rate{Date: "2026-03-18", Base: "CAD", Target: "JPY", Value: 112.50}

type MockFetcher struct {
	rates []rate.Rate
	idx   int
	err   error
}

func (m *MockFetcher) FetchRate(date, base, target string) (rate.Rate, error) {
	if m.err != nil {
		return rate.Rate{}, m.err
	}
	r := m.rates[m.idx]
	m.idx++
	return r, m.err
}

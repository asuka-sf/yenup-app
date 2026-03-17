package rate

import (
	"encoding/json"
	"fmt"
	neturl "net/url"
	"strings"

	domain "yenup/internal/domain/rate"
)

// FrankfurterResponse is the response structure for Frankfurter API
type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// FrankfurterFetcher fetches rates from the Frankfurter API (free, no API key required)
type FrankfurterFetcher struct {
	URL string
}

// NewFrankfurterFetcher creates a new FrankfurterFetcher
func NewFrankfurterFetcher(url string) *FrankfurterFetcher {
	return &FrankfurterFetcher{
		URL: url,
	}
}

// FetchRate fetches the exchange rate for the given date, base, and target currencies
// If the specific date is not available, it falls back to the latest available date
func (f *FrankfurterFetcher) FetchRate(date, base, target string) (domain.Rate, error) {
	url := fmt.Sprintf("%s%s?from=%s&to=%s",
		f.URL,
		neturl.PathEscape(date),
		neturl.QueryEscape(base),
		neturl.QueryEscape(target),
	)

	rate, err := f.fetchFromURL(url, base, target)
	if err != nil {
		// If 404 (not found), try fetching the latest available data
		if strings.Contains(err.Error(), "not found") {
			latestUrl := fmt.Sprintf("%slatest?from=%s&to=%s",
				f.URL,
				neturl.QueryEscape(base),
				neturl.QueryEscape(target),
			)
			return f.fetchFromURL(latestUrl, base, target)
		}
		return domain.Rate{}, err
	}

	return rate, nil
}

// fetchFromURL fetches rate data from a given URL and parses the response
func (f *FrankfurterFetcher) fetchFromURL(url string, base string, target string) (domain.Rate, error) {
	body, err := doGet(url)
	if err != nil {
		return domain.Rate{}, err
	}

	var data FrankfurterResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return domain.Rate{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	rateValue := data.Rates[target]
	if rateValue == 0 {
		return domain.Rate{}, fmt.Errorf("rate for %s not found in response", target)
	}

	return domain.Rate{
		Base:   base,
		Target: target,
		Value:  rateValue,
		Date:   data.Date,
	}, nil
}

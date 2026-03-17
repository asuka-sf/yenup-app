package rate

type RateFetcher interface {
	// Get the rate for a given base and target currency
	FetchRate(date, base, target string) (Rate, error)
}

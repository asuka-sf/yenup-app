package registry

import (
	"yenup/internal/config"
	domainRate "yenup/internal/domain/rate"
	"yenup/internal/handler"
	rateHandler "yenup/internal/handler/rate"
	notifierRepo "yenup/internal/infrastructure/repository/notifier"
	rateRepo "yenup/internal/infrastructure/repository/rate"
	"yenup/internal/usecase"
)

type Registry struct {
	config     *config.Config
	AppHandler *handler.Handler
}

func NewRegistry(cfg *config.Config) (*Registry, error) {
	// Select rate fetcher based on API_PROVIDER config
	var rateFetcher domainRate.RateFetcher
	if cfg.APIProvider == "frankfurter" {
		rateFetcher = rateRepo.NewFrankfurterFetcher(cfg.FrankfurterAPIURL)
	} else {
		rateFetcher = rateRepo.NewExchangeRatesFetcher(cfg.ExchangeRateAPIKey, cfg.ExchangeRateAPIURL)
	}

	slackNotifier := notifierRepo.NewSlackNotifier(cfg.SlackWebhookURL)

	// usecase
	rateUsecase := usecase.NewRateChecker(rateFetcher, slackNotifier)

	// handler
	rateHandler := rateHandler.NewRateHandler(rateUsecase)

	// app handler
	appHandler := handler.NewHandler(rateHandler)

	return &Registry{
		config:     cfg,
		AppHandler: appHandler,
	}, nil
}

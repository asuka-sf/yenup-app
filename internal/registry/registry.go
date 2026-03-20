package registry

import (
	"yenup/internal/config"
	domainRate "yenup/internal/domain/rate"
	"yenup/internal/handler"
	rateHandler "yenup/internal/handler/rate"
	reportHandler "yenup/internal/handler/report"
	notifierRepo "yenup/internal/infrastructure/repository/notifier"
	rateRepo "yenup/internal/infrastructure/repository/rate"
	storageRepo "yenup/internal/infrastructure/repository/storage"
	"yenup/internal/usecase"

	"cloud.google.com/go/storage"
)

type Registry struct {
	config     *config.Config
	AppHandler *handler.Handler
}

func NewRegistry(cfg *config.Config, gcsClient *storage.Client) (*Registry, error) {

	// storageClient provides read/write access to rate data stored in GCS.
	storageClient := storageRepo.NewGCSClient(gcsClient, cfg.GCSBucketName, cfg.GCSObjectName)

	// Select rate fetcher based on API_PROVIDER config
	var rateFetcher domainRate.RateFetcher
	if cfg.APIProvider == "frankfurter" {
		rateFetcher = rateRepo.NewFrankfurterFetcher(cfg.FrankfurterAPIURL)
	} else {
		rateFetcher = rateRepo.NewExchangeRatesFetcher(cfg.ExchangeRateAPIKey, cfg.ExchangeRateAPIURL)
	}

	slackNotifier := notifierRepo.NewSlackNotifier(cfg.SlackWebhookURL)

	// usecase
	rateUsecase := usecase.NewRateChecker(storageClient, rateFetcher, slackNotifier)
	reportUsecase := usecase.NewWeeklyReporter(storageClient, slackNotifier)

	// handler
	rateHandler := rateHandler.NewRateHandler(rateUsecase)
	reportHandler := reportHandler.NewReportHandler(reportUsecase)

	// app handler
	appHandler := handler.NewHandler(rateHandler, reportHandler)

	return &Registry{
		config:     cfg,
		AppHandler: appHandler,
	}, nil
}

package handler

import (
	"yenup/internal/handler/rate"
	"yenup/internal/handler/report"
)

type Handler struct {
	RateHandler   *rate.RateHandler
	ReportHandler *report.ReportHandler
}

func NewHandler(rateHandler *rate.RateHandler, reportHandler *report.ReportHandler) *Handler {
	return &Handler{
		RateHandler:   rateHandler,
		ReportHandler: reportHandler,
	}
}

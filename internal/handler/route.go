package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/check-rate", h.RateHandler.CheckRate)
	r.GET("/weekly-report", h.ReportHandler.GenerateReport)
}

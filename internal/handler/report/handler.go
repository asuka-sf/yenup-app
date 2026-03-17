package report

import (
	"log"
	"net/http"

	"yenup/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Usecase usecase.WeeklyReportUsecase
}

func NewReportHandler(u usecase.WeeklyReportUsecase) *ReportHandler {
	return &ReportHandler{
		Usecase: u,
	}
}

func (h *ReportHandler) GenerateReport(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.Usecase.GenerateReport(ctx); err != nil {
		log.Printf("GenerateReport failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate weekly report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "weekly report generated"})
}

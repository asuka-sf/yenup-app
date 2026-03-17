package report

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "weekly report generated"})
}

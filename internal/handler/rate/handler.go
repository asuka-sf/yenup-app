package rate

import (
	"net/http"
	"strconv"

	"yenup/internal/usecase"

	"github.com/gin-gonic/gin"
)

// RateHandler is the handler for the rate route
type RateHandler struct {
	Usecase usecase.RateCheckUsecase // Changed to interface
}

// NewRateHandler creates a new RateHandler
func NewRateHandler(u usecase.RateCheckUsecase) *RateHandler { // Changed to interface
	return &RateHandler{
		Usecase: u,
	}
}

// Response is the response for the rate route
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RateData is the data returned in the response
type RateData struct {
	Base          string  `json:"base"`
	Target        string  `json:"target"`
	TodayRate     float64 `json:"today_rate"`
	YesterdayRate float64 `json:"yesterday_rate"`
	Change        string  `json:"change"`
	IsNotified    bool    `json:"is_notified"`
}

// CheckRate checks the rate of the base and target currencies
func (h *RateHandler) CheckRate(c *gin.Context) {
	ctx := c.Request.Context()
	base := c.Query("base")
	target := c.Query("target")
	// If notification=true, force sending a Slack message for testing/verification.
	notificationRaw := c.Query("notification")

	if base == "" || target == "" {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "base and target are required",
			Data:    nil,
		})
		return
	}

	forceNotify := false
	if notificationRaw != "" {
		parsed, err := strconv.ParseBool(notificationRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "notification must be a boolean (true/false)",
				Data:    nil,
			})
			return
		}
		forceNotify = parsed
	}

	// use usecase to check the rate
	result, err := h.Usecase.CheckRates(ctx, base, target, forceNotify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Determine change direction
	change := "unchanged"
	if result.TodayRate < result.YesterdayRate {
		change = "down (JPY stronger)"
	} else if result.TodayRate > result.YesterdayRate {
		change = "up (JPY weaker)"
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Rate check executed successfully",
		Data: RateData{
			Base:          base,
			Target:        target,
			TodayRate:     result.TodayRate,
			YesterdayRate: result.YesterdayRate,
			Change:        change,
			IsNotified:    result.IsNotified,
		},
	})
}

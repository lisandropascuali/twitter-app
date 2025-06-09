package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lisandro/timeline-service/internal/usecase"
)

type TimelineHandler struct {
	timelineUseCase usecase.TimelineUseCase
}

func NewTimelineHandler(timelineUseCase usecase.TimelineUseCase) *TimelineHandler {
	return &TimelineHandler{
		timelineUseCase: timelineUseCase,
	}
}

// GetTimeline godoc
// @Summary Get user timeline
// @Description Get timeline of tweets from users that the authenticated user follows
// @Tags timeline
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} domain.Timeline
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /timeline [get]
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		log.Println("Timeline request failed: X-User-ID header is required")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "X-User-ID header is required"})
		return
	}

	log.Printf("Getting timeline for user %s", userID)
	timeline, err := h.timelineUseCase.GetTimeline(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Failed to get timeline for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get timeline"})
		return
	}

	log.Printf("Successfully retrieved timeline for user %s with %d tweets", userID, len(timeline.Tweets))
	c.JSON(http.StatusOK, timeline)
}

type ErrorResponse struct {
	Error string `json:"error"`
} 
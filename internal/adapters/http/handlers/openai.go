package handlers

import (
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for film requests
type OpenAi struct {
	svc ports.OpenAIService
}

// NewUserHandler creates a new UserHandler instance
func NewOpenAi(svc ports.OpenAIService) *OpenAi {
	return &OpenAi{
		svc,
	}
}

// requestAi represents the request body for getting an answer from AI
type requestAi struct {
	Content string `json:"content" binding:"omitempty,required" example:"hi ai who are you?"`
}

// AiRequest godoc
//
//	@Summary		Send a request to AI
//	@Description	Create a request to get a response from the AI
//	@Tags			AI
//	@Accept			json
//	@Produce		json
//	@Param			requestAi	body		requestAi	true	"AI request payload"
//	@Success		200			{object}	aiResponse		"AI response returned successfully"
//	@Failure		400			{object}	errorResponse	"Validation error: Invalid request payload"
//	@Failure		401			{object}	errorResponse	"Unauthorized error: Invalid or missing authentication token"
//	@Failure		404			{object}	errorResponse	"Data not found error"
//	@Failure		409			{object}	errorResponse	"Data conflict error"
//	@Failure		500			{object}	errorResponse	"Internal server error: Failed to process AI request"
//	@Router			/openai [post]
//	@Security		BearerAuth
func (h OpenAi) AiRequest(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)

	var req requestAi
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationError(ctx, err)
		return
	}

	text, err := h.svc.GenerateText(ctx, req.Content, payload.UserID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, NewResponse(text))
}

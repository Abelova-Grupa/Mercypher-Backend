package handlers

import (
	"net/http"
	"time"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	service *services.SessionService
}

func NewSessionHandler(service *services.SessionService) *SessionHandler {
	return &SessionHandler{service: service}
}

func (h *SessionHandler) CreateToken(ctx *gin.Context) {
	var req struct {
		UserID    string          `json:"user_id"`
		Role      string          `json:"role"`
		TokenType token.TokenType `json:"token_type"`
		Duration  time.Duration   `json:"duration"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	token, _, err := h.service.CreateToken(ctx, req.UserID, req.Duration, req.TokenType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"access_token": token})
}

func (h *SessionHandler) VerifyToken(ctx *gin.Context) {
	var req struct {
		TokenToVerify string          `json:"token"`
		TokenType     token.TokenType `json:"token_type"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	payload, err := h.service.VerifyToken(ctx, req.TokenToVerify, req.TokenType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	if payload != nil {
		ctx.JSON(http.StatusOK, gin.H{"verifed": true})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"verified": false})
	}

}


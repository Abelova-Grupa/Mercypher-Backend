package handlers

import (
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// func (h *UserHandler) Register(c *gin.Context) {
// 	var req struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	user, err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Password)
// 	if err != nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})
// }

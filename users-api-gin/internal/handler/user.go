package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"users-api-gin/internal/service"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}


type UserHandler struct {
	service service.UserService
}


func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}



func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.CreateUser(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, resp)
}


func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	resp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}


func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, resp)
}

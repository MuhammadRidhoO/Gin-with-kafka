package controllers

import (
	"go-kafka/internal/dto/request"
	"go-kafka/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ uc service.UserUsecase }

func NewUserHandler(userUC service.UserUsecase) *UserHandler {
	return &UserHandler{uc: userUC}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.Request_User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload", "error": err.Error()})
		return
	}

	user, err := h.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "email already registered") {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.uc.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Get(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user, err := h.uc.Get(c, uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

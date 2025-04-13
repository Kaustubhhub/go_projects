package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaustubhhub/authentication-gin-gonic/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "details": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if err := h.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "User created Successfully",
		"user":    user,
	})
}

func (h *handler) GetUsers(c *gin.Context) {
	var users []models.User

	if err := h.Db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"details": err.Error(),
		})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

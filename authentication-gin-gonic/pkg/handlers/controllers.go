package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaustubhhub/authentication-gin-gonic/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) CheckServerHealth(c *gin.Context) {
	fmt.Println("In the api")
	c.JSON(http.StatusOK, gin.H{"message": "Server is healthy"})
}

func (h *handler) DeleteUser(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	result := h.Db.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted Successfully"})
}

func (h *handler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "details": err.Error()})
		return
	}
	fmt.Println("user : ", user)

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

package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaustubhhub/authentication-gin-gonic/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

func (h *handler) SignIn(c *gin.Context) {
	var user models.User     
	var isUser models.User    

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed"})
		return
	}

	if err := h.Db.Where("username = ?", user.Username).First(&isUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(isUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	// Create JWT token
	claims := jwt.MapClaims{
		"username": isUser.Username,
		"user_type": isUser.UserType,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Secret key for signing (store in env variable in production)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultSecret" // fallback (not safe for production)
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"token":     tokenString,
		"username":  isUser.Username,
		"user_type": isUser.UserType,
	})
}


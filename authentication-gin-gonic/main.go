package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaustubhhub/authentication-gin-gonic/pkg/db"
	"github.com/kaustubhhub/authentication-gin-gonic/pkg/handlers"
)

func main() {
	DB := db.Init()
	h := handlers.New(DB)
	router := gin.Default()

	router.GET("/api/v1/health", h.CheckServerHealth)
	router.GET("/api/v1/users", h.GetUsers)
	router.POST("/api/v1/user", h.CreateUser)
	router.DELETE("/api/v1/user/:id", h.DeleteUser)
	router.POST("/api/v1/signin", h.SignIn)

	router.Run()
}

package handlers

import (
	"net/http"

	"github.com/example/golang-postgres-crud/auth"
	"github.com/example/golang-postgres-crud/models"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if u.Username == "Chek" && u.Password == "123456" {
		tokenString, err := auth.CreateToken(u.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

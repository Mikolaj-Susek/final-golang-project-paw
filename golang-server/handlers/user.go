package handlers

import (
	"net/http"

	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/models"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User
	if err := db.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := db.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := db.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	db.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

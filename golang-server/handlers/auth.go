package handlers

import (
	"net/http"

	"github.com/example/golang-postgres-crud/auth"
	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Creates a new user account with a hashed password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Registration Info"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /register [post]
func RegisterHandler(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("username = ?", u.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this name already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}
	u.Password = string(hashedPassword)

	if err := db.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginHandler godoc
// @Summary      Logs in a user
// @Description  Authenticates a user and returns a JWT token upon successful login.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Login Credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	var u models.User
	var foundUser models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.DB.Where("username = ?", u.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := auth.CreateToken(foundUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

package main

import (
	"net/http"
	"time"

	"github.com/bcantrell1/pro-motocross-api/internal/database"
	"github.com/bcantrell1/pro-motocross-api/internal/env"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
	Secret 	 string `json:"secret" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// RegisterUser registers a new user
// @Summary Register a new user
// @Description Create a new user account with the provided details and the secret
// @Tags auth
// @Accept json
// @Produce json
// @Param register body registerRequest true "User registration data"
// @Success 201 {object} database.User
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Failed to create the user"
// @Router /api/v1/auth/register [post]
func (app *application) registerUser(c *gin.Context) {
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expectedSecret := env.GetEnvString("REGISTER_SECRET", "secret_to_register")
	if register.Secret != expectedSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid registration secret"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong..."})
		return
	}

	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "We ran into an issue creating the user."})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login logs in a user and returns a JWT token
// @Summary Log in a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body loginRequest true "User login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Invalid password"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Error generating token"
// @Router /api/v1/auth/login [post]
func (app *application) login(c *gin.Context) {

	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // This sets our token to expire in 24 hours
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

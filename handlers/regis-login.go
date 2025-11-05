package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mhandyalf/go-passmanager/models"
	"github.com/mhandyalf/go-passmanager/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserName string `json:"username"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{ID: uuid.New(), UserName: input.UserName, Email: input.Email, Password: string(hashed)}
	if err := userRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"` // Changed from UserName
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user, err := userRepo.GetByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Step 1: Request reset
func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userRepo.GetByEmail(input.Email)
	if err != nil {
		// Tetap balas sukses biar aman
		c.JSON(http.StatusOK, gin.H{"message": "If this email exists, a reset link has been sent"})
		return
	}

	// Generate reset token (JWT exp 15 menit)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"type":    "reset",
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err := utils.SendResetEmail(user.Email, tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If this email exists, a reset link has been sent"})
}

// Step 2: Reset password pakai token
func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	userID := claims["user_id"].(string)

	user, err := userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	user.Password = string(hashed)
	if err := userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}

package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mhandyalf/go-passmanager/database"
	"github.com/mhandyalf/go-passmanager/models"

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
	database.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	database.DB.Where("email = ?", input.UserName).First(&user)

	if user.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}) // Ganti ke StatusUnauthorized
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}) // Ganti ke StatusUnauthorized
		return
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	// Tandatangani token dengan secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Kirim token ke klien
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

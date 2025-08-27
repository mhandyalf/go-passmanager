// handlers/password.go
package handlers

import (
	"github.com/mhandyalf/go-passmanager/database"
	"github.com/mhandyalf/go-passmanager/models"
	"github.com/mhandyalf/go-passmanager/utils" // Buat file baru untuk fungsi enkripsi
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePassword ...
func CreatePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input struct {
		Title    string `json:"title" binding:"required"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
		Tags     string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := utils.EncryptAES(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}

	password := models.Password{
		Title:             input.Title,
		Username:          input.Username,
		EncryptedPassword: encryptedPassword,
		Tags:              input.Tags,
		UserID:            userID,
	}

	database.DB.Create(&password)
	c.JSON(http.StatusOK, gin.H{"data": password})
}

// GetPasswords ...
func GetPasswords(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var passwords []models.Password
	database.DB.Where("user_id = ?", userID).Find(&passwords)
	c.JSON(http.StatusOK, gin.H{"data": passwords})
}

// UpdatePassword ...
func UpdatePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input models.Password
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var password models.Password
	if err := database.DB.First(&password, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
		return
	}

	if password.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Perbarui password jika ada di input
	if input.EncryptedPassword != "" {
		encryptedPassword, err := utils.EncryptAES(input.EncryptedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		input.EncryptedPassword = encryptedPassword
	}

	database.DB.Model(&password).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": password})
}

// DeletePassword ...
func DeletePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var password models.Password

	if err := database.DB.First(&password, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
		return
	}

	if password.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	database.DB.Delete(&password)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

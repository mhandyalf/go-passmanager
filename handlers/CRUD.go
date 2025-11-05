package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mhandyalf/go-passmanager/models"
	"github.com/mhandyalf/go-passmanager/repository"
	"github.com/mhandyalf/go-passmanager/utils" // Buat file baru untuk fungsi enkripsi

	"github.com/gin-gonic/gin"
)

var pwdRepo repository.PasswordRepository

// InitHandlers initializes handler-level dependencies (repository, etc.)
func InitHandlers(r repository.PasswordRepository) {
	pwdRepo = r
}

// CreatePassword ...
func CreatePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	log.Printf("[CreatePassword] user_id=%s", userID)

	var input struct {
		Title    string `json:"title" binding:"required"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
		Tags     string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[CreatePassword][ERROR] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[CreatePassword] Input parsed: title=%s username=%s tags=%s", input.Title, input.Username, input.Tags)

	encryptedPassword, err := utils.EncryptAES(input.Password)
	if err != nil {
		log.Printf("[CreatePassword][ERROR] Failed to encrypt password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	log.Println("[CreatePassword] Password encrypted successfully")

	password := models.Password{
		Title:             input.Title,
		Username:          input.Username,
		EncryptedPassword: encryptedPassword,
		Tags:              input.Tags,
		UserID:            userID,
	}

	if err := pwdRepo.Create(&password); err != nil {
		log.Printf("[CreatePassword][ERROR] Failed to insert into DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": password})
}

// GetPasswords ...
func GetPasswords(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	passwords, err := pwdRepo.GetByUserID(userID)
	if err != nil {
		log.Printf("[GetPasswords][ERROR] Failed to fetch passwords: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch passwords"})
		return
	}

	for i := range passwords {
		decrypted, err := utils.DecryptAES(passwords[i].EncryptedPassword)
		if err != nil {
			fmt.Printf("Failed to decrypt password: %v\n", err)
			continue
		}
		passwords[i].DecryptedPassword = decrypted
	}

	c.JSON(http.StatusOK, gin.H{"data": passwords})
}

// UpdatePassword ...
func UpdatePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input models.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var password *models.Password
	p, err := pwdRepo.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
		return
	}
	password = p

	if password.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	updates := map[string]interface{}{
		"title":    input.Title,
		"username": input.Username,
		"tags":     input.Tags,
	}

	if input.Password != "" {
		encryptedPassword, err := utils.EncryptAES(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		updates["encrypted_password"] = encryptedPassword
	}

	if err := pwdRepo.Update(password, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	// reload to return fresh data
	updated, err := pwdRepo.GetByID(fmt.Sprint(password.ID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": password})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeletePassword ...
func DeletePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	p, err := pwdRepo.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
		return
	}

	if p.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	if err := pwdRepo.Delete(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

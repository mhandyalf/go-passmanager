package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendResetEmail(to string, token string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")

	// Gmail SMTP server
	host := "smtp.gmail.com"
	port := "587"
	addr := host + ":" + port

	// Reset password link
	link := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", token)

	// Email message (minimalis, bisa kamu tambahin HTML)
	message := []byte(fmt.Sprintf(
		"Subject: Password Reset Request\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"Hi,\n\nPlease click the link below to reset your password:\n%s\n\n"+
			"If you didn’t request this, please ignore.\n", link))

	// Auth ke Gmail
	auth := smtp.PlainAuth("", from, password, host)

	// Kirim
	err := smtp.SendMail(addr, auth, from, []string{to}, message)
	return err
}

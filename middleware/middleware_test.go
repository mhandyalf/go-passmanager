package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthMiddlewareUsesIdentityFromAuthService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	auth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Fatalf("authorization header was not forwarded")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"valid":true,"user_id":"` + userID.String() + `"}`))
	}))
	defer auth.Close()

	r := gin.New()
	r.Use(AuthMiddleware(auth.URL, auth.Client()))
	r.GET("/private", func(c *gin.Context) {
		got := c.MustGet("user_id").(uuid.UUID)
		c.JSON(http.StatusOK, gin.H{"user_id": got})
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// NewAuthProxy keeps the public API contract stable while authentication is
// served by handy-auth. It contains no credential or token business logic.
func NewAuthProxy(authServiceURL string) (gin.HandlerFunc, error) {
	target, err := url.Parse(authServiceURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		log.Printf("auth service proxy error: %v", err)
		http.Error(w, `{"error":"authentication service unavailable"}`, http.StatusBadGateway)
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}

package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthProxyPreservesPathAndBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/login" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"username":"handy"}` {
			t.Fatalf("body = %q", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"forwarded"}`))
	}))
	defer upstream.Close()

	proxy, err := NewAuthProxy(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.POST("/api/login", proxy)
	proxyServer := httptest.NewServer(r)
	defer proxyServer.Close()

	response, err := http.Post(proxyServer.URL+"/api/login", "application/json", strings.NewReader(`{"username":"handy"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK || !strings.Contains(string(body), "forwarded") {
		t.Fatalf("status = %d, body = %s", response.StatusCode, body)
	}
}

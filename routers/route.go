package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mhandyalf/go-passmanager/database"
	handlers "github.com/mhandyalf/go-passmanager/handlers"
	middleware "github.com/mhandyalf/go-passmanager/middleware"
)

func SetupRouter() *gin.Engine {
	database.ConnectDB()
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // React default
			"http://localhost:8081", // Vue default
			"http://localhost:5173", // Vite default
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:5173",
			"http://38.47.176.19",    // akses lewat IP
			"http://38.47.176.19:80", // kalau spesifik port
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
			"X-Requested-With", "Accept", "Accept-Encoding", "Accept-Language",
			"Connection", "Host", "Referer", "User-Agent",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rute publik (Register, Login)
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	// Rute yang butuh otentikasi
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/passwords", handlers.CreatePassword)
		api.GET("/passwords", handlers.GetPasswords)
		api.PUT("/passwords/:id", handlers.UpdatePassword)
		api.DELETE("/passwords/:id", handlers.DeletePassword)
	}

	return r
}

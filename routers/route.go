package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mhandyalf/go-passmanager/database"
	handlers "github.com/mhandyalf/go-passmanager/handlers"
	"github.com/mhandyalf/go-passmanager/repository"
	middleware "github.com/mhandyalf/go-passmanager/middleware"
)

func SetupRouter() *gin.Engine {
	database.ConnectDB()
	// initialize repository and wire to handlers
	repo := repository.NewPasswordRepository(database.DB)
	handlers.InitHandlers(repo)
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			// local dev
			"http://localhost:3000",
			"http://localhost:8081",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:5173",

			// akses lewat IP VPS (optional, kalau dipakai langsung)
			"http://38.47.176.19",
			"http://38.47.176.19:8081",
			"http://38.47.176.19:8080",

			// domain frontend
			"http://gembolspwmanager.online",
			"https://gembolspwmanager.online",
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

	// Rute publik (ForgotPassword, ResetPassword)
	r.POST("/api/forgot-password", handlers.ForgotPassword)
	r.POST("/api/reset-password", handlers.ResetPassword)

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

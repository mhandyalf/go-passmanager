package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mhandyalf/go-passmanager/database"
	handlers "github.com/mhandyalf/go-passmanager/handlers"
	middleware "github.com/mhandyalf/go-passmanager/middleware"
)

func SetupRouter() *gin.Engine {
	database.ConnectDB()
	r := gin.Default()

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

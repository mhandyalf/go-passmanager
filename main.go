package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	routers "github.com/mhandyalf/go-passmanager/routers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := routers.SetupRouter()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

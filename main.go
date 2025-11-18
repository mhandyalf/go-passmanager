package main

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	routers "github.com/mhandyalf/go-passmanager/routers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	routers := routers.SetupRouter()
	http.Handle("/metrics", promhttp.Handler())
	routers.Run(":8080")
}

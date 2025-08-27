package main

import (
	_ "github.com/joho/godotenv/autoload"
	routers "github.com/mhandyalf/go-passmanager/routers"
)

func main() {
	routers := routers.SetupRouter()
	routers.Run(":8080")
}

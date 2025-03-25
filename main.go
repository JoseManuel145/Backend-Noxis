package main

import (
	infrastructure "Backend/src/Admin/Infrastructure"
	"Backend/src/Admin/Infrastructure/adapters"
	"Backend/src/core"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router := gin.Default()

	router.Use(cors.New(config))

	db := adapters.NewMySQL()

	infrastructure.InitUser(db, router)

	go core.StartConsumer()

	router.Run(":8080")
}

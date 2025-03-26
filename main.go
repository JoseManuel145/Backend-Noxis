package main

import (
	infrastructure "Backend/src/Admin/Infrastructure"
	"Backend/src/Admin/Infrastructure/adapters"
	"Backend/src/core"
	"os"
	"os/signal"
	"syscall"
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		core.StartConsumer()
	}()

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic("Error iniciando el servidor: " + err.Error())
		}
	}()

	<-sigChan
	println("\nApagando el servidor y cerrando conexiones...")
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	infraestructureAdmin "Backend/src/Admin/Infrastructure"
	infraestructureSensor "Backend/src/Alerts/Infrastructure"
	infraestructureKits "Backend/src/Kits/infrastructure"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración de CORS
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Configuración de Gin Router
	router := gin.Default()
	router.Use(cors.New(config))

	// Inicializar las dependencias del Administrador
	infraestructureAdmin.InitUser(router)

	// Inicializar dependencias de la aplicación de alertas
	infraestructureSensor.InitAlerts(router)

	// Inicializar dependencias de los kits
	infraestructureKits.InitKits(router)

	// Configuración para recibir señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic("Error iniciando el servidor: " + err.Error())
		}
	}()
	<-sigChan
}

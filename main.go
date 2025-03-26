package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	infraestructureAdmin "Backend/src/Admin/Infrastructure"
	"Backend/src/Admin/Infrastructure/adapters"
	infraestructureSensor "Backend/src/Alerts/infrastructure"
	"Backend/src/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración de CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Configuración de Gin Router
	router := gin.Default()
	router.Use(cors.New(config))

	// Inicializar base de datos
	db := adapters.NewMySQL()

	// Inicializar rutas de usuarios
	infraestructureAdmin.InitUser(db, router)

	// Conexión a RabbitMQ
	rabbitConn := core.GetRabbitMQ()
	if rabbitConn.Err != "" {
		panic("Error al conectar RabbitMQ: " + rabbitConn.Err)
	}

	// Inicializar dependencias de la aplicación de alertas
	if err := infraestructureSensor.NewDependencies(router, rabbitConn); err != nil {
		panic("Error al inicializar dependencias: " + err.Error())
	}

	// Configuración para recibir señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic("Error iniciando el servidor: " + err.Error())
		}
	}()

	go func() {
		core.StartConsumer()
	}()

	<-sigChan
	println("\nApagando el servidor y cerrando conexiones...")
}

package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnMongo struct {
	Client *mongo.Client
	Err    string
}

func GetMongoDB() *ConnMongo {
	error := ""
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")

	// Crear y conectar el cliente MongoDB en una sola operación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		error = fmt.Sprintf("Error al conectar con MongoDB: %v", err)
		return &ConnMongo{Client: nil, Err: error}
	}

	// Verificar la conexión con un ping
	err = client.Ping(ctx, nil)
	if err != nil {
		error = fmt.Sprintf("Error de conexión con MongoDB: %v", err)
		return &ConnMongo{Client: nil, Err: error}
	}

	return &ConnMongo{Client: client, Err: error}
}

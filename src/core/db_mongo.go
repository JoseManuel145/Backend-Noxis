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
	Client     *mongo.Client
	Database   string
	Collection string
	Err        string
}

func GetMongoDB() *ConnMongo {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	database := os.Getenv("MONGO_DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return &ConnMongo{Err: fmt.Sprintf("Error al conectar con MongoDB: %v", err)}
	}

	// Verificar la conexión con un ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return &ConnMongo{Err: fmt.Sprintf("Error de conexión con MongoDB: %v", err)}
	}

	return &ConnMongo{
		Client:     client,
		Database:   database,
		Collection: collection,
		Err:        "",
	}
}

package initialializers

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectToDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Get MongoDB URI from environment variable
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI environment variable is not set")
    }

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }

    // Ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Set the database
    DB = client.Database(os.Getenv("DB_NAME"))
    
    log.Println("Connected to MongoDB!")
}
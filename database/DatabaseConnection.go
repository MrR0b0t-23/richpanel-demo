package database

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

func ConnectDB() *mongo.Client{
	//check if environment variables is correctly loaded
	err := godotenv.Load(".env") 
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//configures the client and check for errors
	MongoDB := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))
	if err != nil{
		log.Fatal(err.Error())
	}

	//cancel the connection if conneting period exceeds 10 sec
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil{
		log.Fatal(err.Error())
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil{
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to MongoDB")

	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    collection := client.Database("richPanel-webapp").Collection(collectionName)
    return collection
}
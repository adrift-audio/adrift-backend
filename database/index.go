package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Instance MongoInstance

var Collections CollectionsStruct

func ConnectMongo() error {
	databaseConnection := os.Getenv("DATABASE_CONNECTION_STRING")
	databaseName := os.Getenv("DATABASE_NAME")

	if databaseConnection == "" {
		log.Fatal("Missing DATABASE_CONNECTION_STRING")
	}

	if databaseName == "" {
		log.Fatal("Missing DATABASE_NAME")
	}

	client, clientError := mongo.NewClient(options.Client().ApplyURI(databaseConnection))
	if clientError != nil {
		return clientError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionError := client.Connect(ctx)
	db := client.Database(databaseName)

	if connectionError != nil {
		return connectionError
	}

	Instance = MongoInstance{
		Client:   client,
		Database: db,
	}

	Collections = CollectionsStruct{
		Password:   "Password",
		User:       "User",
		UserSecret: "UserSecret",
	}

	fmt.Println("-- database: connected")

	return nil
}

package database

import "go.mongodb.org/mongo-driver/mongo"

type CollectionsStruct struct {
	Password   string
	User       string
	UserSecret string
}

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

package db

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var dbURL = os.Getenv("port")

func GetDBCollection() (*mongo.Collection, error) {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("xoxo1").Collection("users")
	return collection, nil
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
	var url = dbURL

	if url == "" {
		url = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), url)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("xoxo1").Collection(collectionName)
	return collection, nil
}

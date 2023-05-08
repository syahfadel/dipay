package main

import (
	"context"
	"log"
	"tehcTest/routers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db  *mongo.Database
	ctx *context.Context
	err error
)

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	db = client.Database("dipayDB")
}

func main() {
	var PORT = ":4000"
	routers.StartService(db, ctx).Run(PORT)
}

package main

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main(){
	_, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:example@localhost:27017/?authSource=admin"))
	if err != nil {
		panic(err)
	}
}
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Deck struct {
	Name string `bson:"name"`
}

func main() {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:example@localhost:27017/?authSource=admin"))
	if err != nil {
		panic(err)
	}

	collection := client.Database("deck").Collection("deckData")
	
	deck := &Deck{
		Name: "Orzhov Eriette Auras",
	}

	// INSERT
	_, err = collection.InsertOne(context.Background(), deck)
	if err != nil {
		panic(err)
	}

	_, err = collection.InsertMany(context.Background(), bson.D{{"name","Teste"}})
	if err != nil {
		panic(err)
	}

	// Find one
	filter := bson.D{{"name",deck.Name}}
	searchResult := Deck{}
	err = collection.FindOne(context.Background(), filter).Decode(&searchResult)
	if err != nil {
		panic(err.Error())
	}

	// Find
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		deckResult := Deck{}

		err := cur.Decode(&deckResult)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(deckResult)
	}

	// UPDATE
	filter = bson.D{{"name", deck.Name}}
	update := bson.D{{"$set", bson.D{{"name", "Dimir Yuriko Ninjas"}}}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err.Error())
	}

	// DELETE
	filter = bson.D{{"name",deck.Name}}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(updateResult)
}

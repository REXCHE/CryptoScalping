package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoConnection() *mongo.Client {

	password := "APuXI7kPYRKNhaYA"

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://vicarisiventures:" + password + "@cluster0.ing0x.mongodb.net/cluster0?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return client

}

type MarketMakingData struct {
	// Coinbase
	CoinbaseMidpoint float64
	CoinbaseWeighted float64
	CoinbaseBook     []float64

	// Kraken
	KrakenMidpoint float64
	KrakenWeighted float64
	KrakenBook     []float64

	// Gemini
	GeminiMidpoint float64
	GeminiWeighted float64
	GeminiBook     []float64

	// Crypto
	CryptoMidpoint float64
	CryptoWeighted float64
	CryptoBook     []float64

	// FTX US
	FTXMidpoint float64
	FTXWeighted float64
	FTXBook     []float64

	// Additional Shit
	IsSkewed bool
}

func appendMongo(client *mongo.Client, class MarketMakingData, capacity int64, coll_name string) {

	collection := client.Database("MarketMaking").Collection(coll_name)

	// Fetch length of collection
	length, err := collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	if length > capacity {

		for length > capacity {

			// Delete first document
			delete, err := collection.DeleteOne(context.Background(), bson.M{})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(delete)

			length, err = collection.CountDocuments(context.Background(), bson.M{})

			if err != nil {
				log.Fatal(err)
			}

		}

	}

	// Insert newest document
	insert, err := collection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "coinbaseMidpoint", Value: class.CoinbaseMidpoint},
			{Key: "coinbaseWeighted", Value: class.CoinbaseWeighted},
			{Key: "coinbaseBook", Value: class.CoinbaseBook},

			{Key: "krakenMidpoint", Value: class.KrakenMidpoint},
			{Key: "krakenWeighted", Value: class.KrakenWeighted},
			{Key: "krakenBook", Value: class.KrakenBook},

			{Key: "geminiMidpoint", Value: class.GeminiWeighted},
			{Key: "geminiWeighted", Value: class.GeminiWeighted},
			{Key: "geminibook", Value: class.GeminiBook},

			{Key: "cryptoMidpoint", Value: class.CryptoMidpoint},
			{Key: "cryptoWeighted", Value: class.CryptoWeighted},
			{Key: "cryptoBook", Value: class.CryptoBook},

			{Key: "ftxMidpoint", Value: class.FTXMidpoint},
			{Key: "ftxWeighted", Value: class.FTXWeighted},
			{Key: "ftxBook", Value: class.FTXBook},

			{Key: "isSkewed", Value: class.IsSkewed},
		})

	if err != nil {
		log.Fatal("Err with database")
	}

	fmt.Println(insert)

}

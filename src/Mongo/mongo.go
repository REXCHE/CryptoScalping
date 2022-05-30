package Mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoConnection() *mongo.Client {

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

func AppendMongo(client *mongo.Client, class MarketMakingData, capacity int64, coll_name string) {

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
			{Key: "spread", Value: class.Spread},

			{Key: "open", Value: class.Open},
			{Key: "high", Value: class.High},
			{Key: "low", Value: class.Low},
			{Key: "close", Value: class.Close},

			{Key: "recentTrades", Value: class.RecentTrades},
			{Key: "recentVolatility", Value: class.Volatility},

			{Key: "gamma", Value: class.Gamma},
			{Key: "kappa", Value: class.Kappa},
			{Key: "tau", Value: class.Tau},
			{Key: "sigma", Value: class.Sigma},

			// TODO

		})

	if err != nil {
		log.Fatal("Err with database")
	}

	fmt.Println(insert)

}

func FetchMongoDB(client *mongo.Client, coll_name string) []MarketMakingData {

	var MMD []MarketMakingData

	collection := client.Database("MarketMaking").Collection("OrderBooks")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var iterations []bson.M
	err = cursor.All(ctx, &iterations)

	if err != nil {
		log.Fatal(err)
	}

	var mmd MarketMakingData

	for _, itr := range iterations {

		// Coinbase
		mmd.CoinbaseMidpoint = itr["coinbaseMidpoint"].(float64)
		mmd.CoinbaseWeighted = itr["coinbaseWeighted"].(float64)

		// Kraken
		mmd.KrakenMidpoint = itr["krakenMidpoint"].(float64)
		mmd.KrakenWeighted = itr["krakenWeighted"].(float64)

		// Gemini
		mmd.GeminiMidpoint = itr["geminiMidpoint"].(float64)
		mmd.GeminiWeighted = itr["geminiWeighted"].(float64)

		// Crypto
		mmd.CryptoMidpoint = itr["cryptoMidpoint"].(float64)
		mmd.CryptoWeighted = itr["cryptoWeighted"].(float64)

		// FTX US
		mmd.FTXMidpoint = itr["ftxMidpoint"].(float64)
		mmd.FTXWeighted = itr["ftxWeighted"].(float64)

		// Additional Shit
		for i := 0; i < len(itr["coinbaseBook"].(primitive.A)); i++ {
			mmd.CoinbaseBook = append(mmd.CoinbaseBook, itr["coinbaseBook"].(primitive.A)[i].(float64))
			mmd.KrakenBook = append(mmd.KrakenBook, itr["krakenBook"].(primitive.A)[i].(float64))
			mmd.GeminiBook = append(mmd.GeminiBook, itr["geminibook"].(primitive.A)[i].(float64))
			mmd.CryptoBook = append(mmd.CryptoBook, itr["cryptoBook"].(primitive.A)[i].(float64))
			mmd.FTXBook = append(mmd.FTXBook, itr["ftxBook"].(primitive.A)[i].(float64))
		}

		// More Additional Shit
		mmd.IsSkewed = itr["isSkewed"].(bool)

		MMD = append(MMD, mmd)

	}

	return MMD

}

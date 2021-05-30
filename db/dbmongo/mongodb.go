package dbmongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var err error

// collections
var ChatUserCol *mongo.Collection
var ChatContact *mongo.Collection
var ChatCommunity *mongo.Collection

func init() {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("connect error ", err.Error())
		return
	}

	/*defer func() {
		if err = MongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()*/
	//ping

	if err := MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	ChatUserCol = MongoClient.Database("chat").Collection("user")
	ChatContact = MongoClient.Database("chat").Collection("contact")
	ChatCommunity = MongoClient.Database("chat").Collection("community")
}
func CloseClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = MongoClient.Disconnect(ctx); err != nil {
		panic(err)
	}
}

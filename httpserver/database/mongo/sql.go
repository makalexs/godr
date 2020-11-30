package postgres

import (
	"context"
	"github.com/makalexs/godr/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func CommonInsert(request interface{}) string {
	uri := "mongodb://"+config.GetConfiguration().DatabaseMongo.Url+":"+config.GetConfiguration().DatabaseMongo.Port
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer client.Disconnect(ctx)

	collection := client.Database("drotus").Collection("objects")
	result, err := collection.InsertOne(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return result.InsertedID.(primitive.ObjectID).Hex()
}

func CommonFind(request bson.M, limit int64) []interface{} {
	uri := "mongodb://"+config.GetConfiguration().DatabaseMongo.Url+":"+config.GetConfiguration().DatabaseMongo.Port
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer client.Disconnect(ctx)

	collection := client.Database("drotus").Collection("objects")
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	result, err := collection.Find(context.TODO(), request, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var resultArray []interface{}
	var resultOne bson.M
	for result.Next(ctx) {
		err = result.Decode(&resultOne)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		resultArray = append(resultArray, resultOne)
	}

	return resultArray
}
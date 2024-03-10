package database

import (
	"context"
	"fmt"
	"lxdexplorer-collector/config"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conf = config.Conf

var database = "lxd-dev"

func connect() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.GetDatabaseURI()))
	if err != nil {
		panic(err)
	}

	ping(client)

	return client
}

func ping(c *mongo.Client) {
	err := c.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

func InsertOne(collection string, document interface{}) {
	c := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database(database).Collection(collection).InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
	log.Printf("Database: Inserted document with ID %v\n", i.InsertedID)

}

func InsertMany(collection string, documents []interface{}) {
	c := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database(database).Collection(collection).InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Database: Inserted %v documents with IDs %v\n", len(i.InsertedIDs), i.InsertedIDs)

}

func FindOne(collection string, filter interface{}) *mongo.SingleResult {
	c := connect()
	defer c.Disconnect(context.Background())

	return c.Database(database).Collection(collection).FindOne(context.Background(), filter)
}

func FindAll(collection string) ([]primitive.M, error) {
	c := connect()
	defer c.Disconnect(context.Background())

	cur, err := c.Database(database).Collection(collection).Find(context.Background(), bson.D{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil

}

func ReplaceOne(collection string, filter interface{}, replacement interface{}) (*mongo.UpdateResult, error) {
	c := connect()
	defer c.Disconnect(context.Background())

	return c.Database(database).Collection(collection).ReplaceOne(context.Background(), filter, replacement)
}

func AddTTL(collection string, field string, seconds int32) {
	c := connect()
	defer c.Disconnect(context.Background())

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}

	indexView := c.Database(database).Collection(collection).Indexes()

	// Create a new index
	var err error
	_, err = indexView.CreateOne(context.Background(), indexModel)
	if err != nil {
		// Drop the existing index
		log.Println(err)
		_, err := indexView.DropOne(context.Background(), string(field+"_1"))
		if err != nil {
			// Handle error
			fmt.Println(err)
		}
		// log.Printf("Database: Dropped existing TTL index on %s\n", field)
		// Create a new index
		_, err = indexView.CreateOne(context.Background(), indexModel)
		if err != nil {
			// Handle error
			log.Println(err)
		}
	}

}

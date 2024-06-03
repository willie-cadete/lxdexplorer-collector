package database

import (
	"context"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database = "lxd"

func (s *Database) InsertOne(collection string, document interface{}) {
	c := s.connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database(database).Collection(collection).InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
	log.Debugf("Database: Inserted document with ID %v\n", i.InsertedID)

}

// TODO: Return Error
func (s *Database) InsertMany(collection string, documents []interface{}) {
	c := s.connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database(database).Collection(collection).InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	log.Debugf("Database: Inserted %v documents with IDs %v\n", len(i.InsertedIDs), i.InsertedIDs)

}

func (s *Database) FindOne(collection string, filter interface{}) *mongo.SingleResult {
	c := s.connect()
	defer c.Disconnect(context.Background())

	return c.Database(database).Collection(collection).FindOne(context.Background(), filter)
}

func (s *Database) FindAll(collection string) ([]primitive.M, error) {
	c := s.connect()
	defer c.Disconnect(context.Background())

	cur, err := c.Database(database).Collection(collection).Find(context.Background(), bson.D{})

	if err != nil {
		log.Info(err)
		return nil, err
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Info(err)
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil

}

func (s *Database) ReplaceOne(collection string, filter interface{}, replacement interface{}) (*mongo.UpdateResult, error) {
	c := s.connect()
	defer c.Disconnect(context.Background())

	return c.Database(database).Collection(collection).ReplaceOne(context.Background(), filter, replacement)
}

// TODO: Return Error
func (s *Database) AddTTL(collection string, field string, seconds int32) error {
	c := s.connect()
	defer c.Disconnect(context.Background())

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}

	indexView := c.Database(database).Collection(collection).Indexes()

	// Create a new index
	_, err := indexView.CreateOne(context.Background(), indexModel)
	if err != nil {
		// Drop the existing indexx``
		_, err := indexView.DropOne(context.Background(), string(field+"_1"))
		if err != nil {
			return err
		}
		log.Debugln("Database: Dropped existing TTL index on", collection, "field", field)
		// Create a new index
		_, err = indexView.CreateOne(context.Background(), indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Database) connect() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(s.config.GetDatabaseURI()))
	if err != nil {
		panic(err)
	}

	s.ping(client)

	return client
}

func (s *Database) ping(c *mongo.Client) {
	err := c.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

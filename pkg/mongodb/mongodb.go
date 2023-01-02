package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	db          *mongo.Database
	mongoClient *mongo.Client
}

func NewClient(uri string) (client *Client, err error) {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := c.Database("syncbyte")

	return &Client{db: db, mongoClient: c}, nil
}

func (c *Client) Close() error {
	return c.mongoClient.Disconnect(context.TODO())
}

func (c *Client) GetCollection(col string) *mongo.Collection {
	return c.db.Collection(col)
}

// func (c *Client) InsertOne(col string, document interface{}) error {
// 	coll := c.db.Collection(col)

// 	if _, err := coll.InsertOne(context.Background(), document); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (c *Client) Ping() error {
	return c.mongoClient.Ping(context.TODO(), readpref.Primary())
}

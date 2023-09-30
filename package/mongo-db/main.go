package mongo_db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DBInterface interface {
	ConnectDB(ctx context.Context) error
	InsertRecord(ctx context.Context, collectionName string, record any) error
	UpdateRecord(ctx context.Context, collectionName string, filter, update interface{}) error
	DeleteRecord(ctx context.Context, collectionName string, filter interface{}) error
	GetCount(ctx context.Context, collectionName string, filter interface{}) (int64, error)
	GetByID(ctx context.Context, collectionName string, id interface{}, result interface{}) error
	GetAllRecords(ctx context.Context, collectionName string) (*mongo.Cursor, error)
	CreateCollection(ctx context.Context, collectionName string) error
}

type DB struct {
	connectionString string
	database         string
	client           *mongo.Client
}

func (d *DB) ConnectDB(_ context.Context) error {
	clientOptions := options.Client().ApplyURI(d.connectionString)

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.client = client

	fmt.Println("Connected to MongoDB!")
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		fmt.Println("if")
	//		log.Fatal(err)
	//	}
	//}()
	return nil
}

func (d *DB) InsertRecord(ctx context.Context, collectionName string, record any) error {
	_, err := d.client.Database(d.database).Collection(collectionName).InsertOne(ctx, record)
	return err
}

func (d *DB) UpdateRecord(ctx context.Context, collectionName string, filter, update interface{}) error {
	_, err := d.client.Database(d.database).Collection(collectionName).UpdateOne(ctx, filter, update)
	return err
}

func (d *DB) DeleteRecord(ctx context.Context, collectionName string, filter interface{}) error {
	_, err := d.client.Database(d.database).Collection(collectionName).DeleteOne(ctx, filter)
	return err
}

func (d *DB) GetCount(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	count, err := d.client.Database(d.database).Collection(collectionName).CountDocuments(ctx, filter)
	return count, err
}

func (d *DB) GetByID(ctx context.Context, collectionName string, id interface{}, result interface{}) error {
	filter := bson.M{"_id": id}
	err := d.client.Database(d.database).Collection(collectionName).FindOne(ctx, filter).Decode(result)
	return err
}

func (d *DB) GetAllRecords(ctx context.Context, collectionName string) (*mongo.Cursor, error) {
	cursor, err := d.client.Database(d.database).Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (d *DB) CreateCollection(ctx context.Context, collectionName string) error {
	options := options.CreateCollection().SetCapped(false) // You can adjust the options as needed
	err := d.client.Database(d.database).CreateCollection(ctx, collectionName, options)
	return err
}

func NewDB(database string, connectionUrl string) DBInterface {
	return &DB{
		database:         database,
		connectionString: connectionUrl,
		client:           nil,
	}
}

package mongo_db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DB struct {
	connectionString string
	databaseName     string
	client           *mongo.Client
	db               *mongo.Database
}

func cursorToStruct[T any](ctx context.Context, curr *mongo.Cursor) ([]T, error) {
	var m []map[string]any
	var results []T

	if err := curr.All(ctx, &m); err != nil {
		return results, err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return results, err
	}

	if err := json.Unmarshal(b, &results); err != nil {
		return results, err
	}

	return results, nil
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
	d.db = client.Database(d.databaseName)

	fmt.Println("Connected to MongoDB!")
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		fmt.Println("if")
	//		log.Fatal(err)
	//	}
	//}()
	return nil
}

func (d *DB) InsertRecord(ctx context.Context, collectionName string, record any) (any, error) {
	//_, err := d.db.Collection(collectionName).InsertOne(ctx, record)
	record, err := d.db.Collection(collectionName).InsertOne(ctx, record)
	return record, err
}

func (d *DB) InsertMany(ctx context.Context, collectionName string, records []any) error {
	_, err := d.db.Collection(collectionName).InsertMany(ctx, records)
	return err
}

func (d *DB) UpdateRecord(ctx context.Context, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	result, err := d.db.Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) UpdateByID(ctx context.Context, collectionName string, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	result, err := d.db.Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) DeleteRecord(ctx context.Context, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := d.db.Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) DeleteByID(ctx context.Context, collectionName string, id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}
	result, err := d.db.Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) GetCount(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	count, err := d.db.Collection(collectionName).CountDocuments(ctx, filter)
	return count, err
}

func (d *DB) Find(ctx context.Context, collectionName string) (*mongo.Cursor, error) {
	cursor, err := d.db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
	//return cursorToStruct(ctx, cursor)
}

func (d *DB) FindOne(ctx context.Context, collectionName string, filter interface{}) *mongo.SingleResult {
	return d.db.Collection(collectionName).FindOne(ctx, filter)
}

func (d *DB) FindByID(ctx context.Context, collectionName string, id string) *mongo.SingleResult {
	//filter := bson.M{"_id": id}
	filter := Filter{"id": id}
	//err := d.db.Collection(collectionName).FindOne(ctx, filter).Decode(result)
	return d.db.Collection(collectionName).FindOne(ctx, filter)
}

func (d *DB) CreateCollection(ctx context.Context, collectionName string) error {
	options := options.CreateCollection().SetCapped(false) // You can adjust the options as needed
	err := d.db.CreateCollection(ctx, collectionName, options)
	return err
}

func NewDB(databaseName string, connectionUrl string) DBRepo {
	return &DB{
		databaseName:     databaseName,
		connectionString: connectionUrl,
		client:           nil,
		db:               nil,
	}
}

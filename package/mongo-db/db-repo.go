package mongo_db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Entity interface {
	GetPrimitiveID() ID
	GetId() ID
	SetId(id ID)
	GetCreationTime() time.Time
	GetUpdateTime() time.Time
	SetCreationTime(time.Time)
	SetUpdateTime(time.Time)
	IsZero() bool

	//IncrementRecordVersion()
	//GetRecordVersion() int
	//IsMarkedForDeletion() bool
}

type Filter map[string]interface{}

type Query struct {
	Filter Filter
	Sort   map[string]interface{}
}

type ID string

type DBRepo[T Entity] interface {
	ConnectDB(ctx context.Context) error

	InsertRecord(ctx context.Context, collectionName string, record any) (any, error)
	InsertMany(ctx context.Context, collectionName string, record []any) error

	UpdateRecord(ctx context.Context, collectionName string, filter any, update interface{}) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, collectionName string, id interface{}, update interface{}) (*mongo.UpdateResult, error)

	DeleteRecord(ctx context.Context, collectionName string, filter interface{}) (*mongo.DeleteResult, error)
	DeleteByID(ctx context.Context, collectionName string, id string) (*mongo.DeleteResult, error)

	GetCount(ctx context.Context, collectionName string, filter interface{}) (int64, error)

	Find(ctx context.Context, collectionName string) ([]T, error)
	FindOne(ctx context.Context, collectionName string, filter interface{}) *mongo.SingleResult
	FindByID(ctx context.Context, collectionName string, id string) *mongo.SingleResult
	CreateCollection(ctx context.Context, collectionName string) error
}

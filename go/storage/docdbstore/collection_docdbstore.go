package docdbstore

import (
	"context"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionDocumentDBStore struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func CreateCollectionDocumentDBStore(connectionUri string) (*CollectionDocumentDBStore, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	dbStore := CollectionDocumentDBStore{
		Client: client,
		DB:     client.Database("dwn_collections_db"),
	}

	return &dbStore, nil

}

func (store *CollectionDocumentDBStore) GetCollectionItem(identifier string) (*storage.CollectionItem, error) {

	// Get Collection
	var collectionItem storage.CollectionItem
	result := store.DB.Collection("dwn_collections").FindOne(context.Background(), bson.D{{"_id", identifier}})
	err := result.Decode(&collectionItem)
	if err != nil {
		return nil, err
	}

	return &collectionItem, nil

}

func (store *CollectionDocumentDBStore) PutCollectionItem(collectionItem *storage.CollectionItem) error {

	if collectionItem.ID.IsZero() {
		// this is a new item
		collectionItem.ID = primitive.NewObjectID()
	}

	_, err := store.DB.Collection("dwn_collections").InsertOne(context.Background(), collectionItem)
	if err != nil {
		return err
	}

	return nil

}

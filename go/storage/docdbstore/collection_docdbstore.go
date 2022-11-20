package docdbstore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionDocumentDBStore struct {
	SchemaURI string
	Client    *mongo.Client
	DB        *mongo.Database
}

func CreateCollectionDocumentDBStore(schemaUri string, connectionUri string) (*CollectionDocumentDBStore, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	dbStore := CollectionDocumentDBStore{
		SchemaURI: schemaUri,
		Client:    client,
		DB:        client.Database("dwn_collections"),
	}

	return &dbStore, nil

}

func (store *CollectionDocumentDBStore) GetSchemaURI() string {
	return store.SchemaURI
}

func (store *CollectionDocumentDBStore) GetCollectionItem(identifier any) any {
	//TODO implement me
	panic("implement me")
}

func (store *CollectionDocumentDBStore) PutCollectionItem(collectionItem any) {
	//TODO implement me
	panic("implement me")
}

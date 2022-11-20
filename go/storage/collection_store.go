package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type CollectionStore interface {
	GetSchemaURI() string
	GetCollectionItem(identifier string) (*CollectionItem, error)
	PutCollectionItem(collectionItem *CollectionItem) error
}

type CollectionItem struct {
	ID        primitive.ObjectID `bson:"_id"`
	SchemaURI string             `bson:"schema_uri"`
}

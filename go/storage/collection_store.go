package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type CollectionStore interface {
	GetCollectionItem(identifier string) (*CollectionItem, error)
	PutCollectionItem(collectionItem *CollectionItem) error
}

type CollectionItem struct {
	ID           primitive.ObjectID `bson:"_id"`
	OwnerDID     string             `bson:"owner_did"`
	AuthorDID    string             `bson:"author_did"`
	RecipientDID string             `bson:"recipient_did"`
	SchemaURI    string             `bson:"schema_uri"`
	Content      []byte             `bson:"content"`
}

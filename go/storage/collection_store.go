package storage

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type CollectionStore interface {
	GetSchemaURI() string
	GetCollectionItem(identifier any) any
	PutCollectionItem(collectionItem any)
}

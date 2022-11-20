package collection

import (
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StoreCollection(collectionStore storage.CollectionStore, schemaURI string, collectionItemId string, collectionItemBytes []byte) (string, error) {

	var collectionItem storage.CollectionItem
	if collectionItemId != "" {
		objectId, err := primitive.ObjectIDFromHex(collectionItemId)
		if err != nil {
			return "", err
		}
		collectionItem.ID = objectId
	}

	collectionItem.SchemaURI = schemaURI
	collectionItem.Content = collectionItemBytes

	err := collectionStore.PutCollectionItem(&collectionItem)
	if err != nil {
		return "", err
	}

	return collectionItem.ID.String(), nil

}

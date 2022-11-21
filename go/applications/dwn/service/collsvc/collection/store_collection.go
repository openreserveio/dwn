package collection

import (
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StoreCollection(collectionStore storage.CollectionStore, schemaURI string, collectionItemId string, collectionItemBytes []byte, authorDID string, recipientDID string) (string, string, error) {

	var collectionItem storage.CollectionItem
	if collectionItemId != "" {
		objectId, err := primitive.ObjectIDFromHex(collectionItemId)
		if err != nil {
			return "", "", err
		}
		collectionItem.ID = objectId
	}

	collectionItem.SchemaURI = schemaURI
	collectionItem.Content = collectionItemBytes
	collectionItem.AuthorDID = authorDID
	collectionItem.RecipientDID = recipientDID
	collectionItem.OwnerDID = recipientDID

	err := collectionStore.PutCollectionItem(&collectionItem)
	if err != nil {
		return "", "", err
	}

	return collectionItem.ID.Hex(), collectionItem.OwnerDID, nil

}

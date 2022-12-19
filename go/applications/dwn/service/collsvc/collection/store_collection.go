package collection

import (
	"github.com/openreserveio/dwn/go/storage"
)

type StoreCollectionResult struct {
	RecordID string
}

func StoreCollection(collectionStore storage.CollectionStore, schemaURI string, recordId string, parentId string, descriptorCID string, processingCID string, dataBytes []byte, authorDID string, recipientDID string) (*StoreCollectionResult, error) {

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing

	return &StoreCollectionResult{}, nil
}

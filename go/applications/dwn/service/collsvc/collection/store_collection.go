package collection

import (
	"github.com/openreserveio/dwn/go/storage"
)

func StoreCollection(collectionStore storage.CollectionStore, schemaURI string, recordId string, parentId string, descriptorCID string, processingCID string, dataBytes []byte, authorDID string, recipientDID string) (string, string, error) {

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing

	return "", "", nil
}

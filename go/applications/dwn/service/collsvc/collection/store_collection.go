package collection

import (
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
)

type StoreCollectionResult struct {
	RecordID string
}

func StoreCollection(collectionStore storage.CollectionStore, methodName string, schemaURI string, recordId string, parentId string, descriptorCID string, processingCID string, dataBytes []byte, authorDID string, recipientDID string) (*StoreCollectionResult, error) {

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing
	switch methodName {

	case model.METHOD_COLLECTIONS_WRITE:

	case model.METHOD_COLLECTIONS_COMMIT:

	case model.METHOD_COLLECTIONS_DELETE:

	default:

	}

	return &StoreCollectionResult{}, nil
}

func collectionsWrite() {

	/*
		Generate the message’s Entry ID by performing the Record ID Generation Process.
		IF the generated Entry ID matches the recordId value of the message it is the Initial Entry for a record, store the entry as the Initial Entry for the record if no Initial Entry exists and cease any further processing.
		ELSE the message may be an overwriting entry for the record; continue processing.
	*/

}

package collection

import (
	"errors"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	ERR_DUPLICATE_INITIAL_ENTRY = "Trying to write an Initial Entry to a record that already exists"
)

type StoreCollectionResult struct {
	Status   string
	RecordID string
	Error    error
}

type CollectionsWriteParams struct {
	RecordID  string `json:"recordId,omitempty"`
	ContextID string `json:"contextId,omitempty"`
	Data      string `json:"data,omitempty"`

	ProcessingNonce string `json:"nonce"`
	AuthorDID       string `json:"author"`
	RecipientDID    string `json:"recipient"`

	Method          string    `json:"method"`
	DataCID         string    `json:"dataCid,omitempty"`
	DataFormat      string    `json:"dataFormat,omitempty"`
	ParentID        string    `json:"parentId,omitempty"`
	Protocol        string    `json:"protocol,omitempty"`
	ProtocolVersion string    `json:"protocolVersion,omitempty"`
	Schema          string    `json:"schema,omitempty"`
	CommitStrategy  string    `json:"commitStrategy,omitempty"`
	Published       bool      `json:"published,omitempty"`
	DateCreated     time.Time `json:"dateCreated,omitempty"`
	DatePublished   time.Time `json:"datePublished,omitempty"`
}

func StoreCollection(collectionStore storage.CollectionStore, collectionsWriteParams *CollectionsWriteParams) (*StoreCollectionResult, error) {

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing
	result := StoreCollectionResult{}
	switch collectionsWriteParams.Method {

	case model.METHOD_COLLECTIONS_WRITE:
		err := collectionsWrite(collectionStore, collectionsWriteParams)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err
			return &result, nil
		}
		result.Status = "OK"
		result.RecordID = collectionsWriteParams.RecordID

	case model.METHOD_COLLECTIONS_COMMIT:

	case model.METHOD_COLLECTIONS_DELETE:

	default:
		result.Status = "UNSUPPORTED_METHOD"

	}

	return &result, nil
}

func collectionsWrite(collectionStore storage.CollectionStore, params *CollectionsWriteParams) error {

	/*
			Generate the message’s Entry ID by performing the Record ID Generation Process.
			- IF the generated Entry ID matches the recordId value of the message it is the Initial Entry
		      for a record, store the entry as the Initial Entry for the record if no Initial Entry exists
		      and cease any further processing.
			- ELSE the message may be an overwriting entry for the record; continue processing.
	*/
	descriptor := model.Descriptor{
		Method:          model.METHOD_COLLECTIONS_WRITE,
		DataCID:         model.CreateDataCID(params.Data),
		DataFormat:      params.DataFormat,
		ParentID:        params.ParentID,
		Protocol:        params.Protocol,
		ProtocolVersion: params.ProtocolVersion,
		Schema:          params.Schema,
		CommitStrategy:  params.CommitStrategy,
		Published:       params.Published,
		DateCreated:     params.DateCreated,
		DatePublished:   params.DatePublished,
	}
	entryId := model.CreateDescriptorCID(descriptor)

	if entryId == params.RecordID {

		// This is the first entry of the record.  Create it and return
		// If there is an existing record id, there's a problem and return an error
		existingRecord := collectionStore.GetCollectionRecord(entryId)
		if existingRecord != nil {
			return errors.New(ERR_DUPLICATE_INITIAL_ENTRY)
		}

		record := storage.CollectionRecord{
			ID:                      primitive.NewObjectID(),
			RecordID:                entryId,
			CreatorDID:              params.AuthorDID,
			OwnerDID:                params.RecipientDID,
			WriterDIDs:              []string{params.AuthorDID},
			ReaderDIDs:              []string{params.AuthorDID, params.RecipientDID},
			InitialEntryID:          entryId,
			LatestEntryID:           entryId,
			LatestCheckpointEntryID: entryId,
		}

		entry := storage.MessageEntry{
			ID:             primitive.NewObjectID(),
			MessageEntryID: entryId,
			ParentEntryID:  "",
			RecordID:       entryId,
			Schema:         params.Schema,
			Method:         params.Method,
			Data:           []byte(params.Data),
			DataCID:        params.DataCID,
			Protocol:       params.Protocol,
			CreatedDate:    params.DateCreated,
			PublishedDate:  params.DatePublished,
		}

		err := collectionStore.CreateCollectionRecord(&record, &entry)
		return err

	} else if true {

		// This is an attempt to overwrite a previous version.
		// So, let's get the parent version
		parentCollRec := collectionStore.GetCollectionRecord(params.ParentID)
		if parentCollRec == nil {
			// We dont have the parent.  Reject with err
			return errors.New("Unable to find Parent Record for Overwrite")
		}
	}

	return errors.New("Supports initial entry only for now")

}

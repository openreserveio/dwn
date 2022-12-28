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
	messageProcessing := model.MessageProcessing{
		Nonce:        params.ProcessingNonce,
		AuthorDID:    params.AuthorDID,
		RecipientDID: params.RecipientDID,
	}
	descriptorId := model.CreateDescriptorCID(descriptor)
	processingId := model.CreateProcessingCID(messageProcessing)
	entryId := model.CreateRecordCID(descriptorId, processingId)

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
			ID:              primitive.NewObjectID(),
			MessageEntryID:  entryId,
			ParentEntryID:   "",
			RecordID:        entryId,
			Schema:          params.Schema,
			Method:          params.Method,
			Data:            []byte(params.Data),
			DataCID:         params.DataCID,
			Protocol:        params.Protocol,
			ProtocolVersion: params.ProtocolVersion,
			CreatedDate:     params.DateCreated,
			PublishedDate:   params.DatePublished,
		}

		err := collectionStore.CreateCollectionRecord(&record, &entry)
		return err

	} else if true {

		// the message may be an overwriting entry for the record; continue processing.
		// This is an attempt to overwrite a previous version.
		// So, let's get the parent version
		parentCollRec := collectionStore.GetCollectionRecord(params.ParentID)
		if parentCollRec == nil {
			// If a message is not the Initial Entry, its descriptor MUST contain a parentId to
			// determine the entry’s position in the record’s lineage. If a parentId is present
			// proceed with processing, else discard the record and cease processing.
			// We dont have the parent.  Reject with err
			return errors.New("Unable to find Parent Record for Overwrite")
		}

		// Ensure all immutable values from the Initial Entry remained unchanged if present in the
		// inbound message. If any have been mutated, discard the message and cease processing.
		initialMessageEntry := collectionStore.GetMessageEntryByID(parentCollRec.InitialEntryID)
		if initialMessageEntry == nil {
			return errors.New("Unable to find an initial entry")
		}
		if initialMessageEntry.Protocol != params.Protocol ||
			initialMessageEntry.ProtocolVersion != params.ProtocolVersion ||
			initialMessageEntry.Schema != params.Schema {
			return errors.New("Attempt to mutate an immutable value")
		}

		// Retrieve the Latest Checkpoint Entry, which will be either the Initial Entry or the latest CollectionsDelete,
		// and compare the parentId value of the inbound message to the Entry ID of the
		// Latest Checkpoint Entry derived from running the Record ID Generation Process on it.
		// If the values match, proceed with processing, if the values do not match discard the message and cease processing.
		latestCheckpointEntry := collectionStore.GetMessageEntryByID(parentCollRec.LatestCheckpointEntryID)
		if latestCheckpointEntry == nil {
			return errors.New("Unable to find the latest checkpoint entry")
		}

		if params.ParentID != latestCheckpointEntry.MessageEntryID {
			return errors.New("The parent ID of the inbound message must match the latest checkpoint ID.")
		}

		// If an existing CollectionsWrite entry linked to the Latest Checkpoint Entry IS NOT present and
		// the dateCreated value of the inbound message is greater than the Latest Checkpoint Entry,
		// store the message as the Latest Entry and cease processing, else discard the inbound message
		// and cease processing.
		if latestCheckpointEntry.Method != model.METHOD_COLLECTIONS_WRITE &&
			params.DateCreated.After(latestCheckpointEntry.CreatedDate) {
			latestEntry := storage.MessageEntry{
				ID:              primitive.NewObjectID(),
				MessageEntryID:  entryId,
				ParentEntryID:   params.ParentID,
				RecordID:        params.RecordID,
				Schema:          params.Schema,
				Method:          params.Method,
				Data:            []byte(params.Data),
				DataCID:         params.DataCID,
				Protocol:        params.Protocol,
				ProtocolVersion: params.ProtocolVersion,
				CreatedDate:     params.DateCreated,
				PublishedDate:   params.DatePublished,
			}

			err := collectionStore.AddCollectionMessageEntry(&latestEntry)
			if err != nil {
				return err
			}

			parentCollRec.LatestEntryID = entryId
			err = collectionStore.SaveCollectionRecord(parentCollRec)
			if err != nil {
				return err
			}

		}

		// If an exiting CollectionsWrite entry linked to the Latest Checkpoint Entry IS present
		// all of the following conditions MUST be true:
		//   - The dateCreated value of the inbound message is greater than the existing CollectionsWrite,
		//     or if the dateCreated values are the same, the Entry ID of the inbound message is greater
		//     than the existing entry when the Entry IDs of the two are compared lexicographically.
		// If all of the following conditions for Step 6 are true, store the inbound message as the Latest Entry
		// and discard the existing CollectionsWrite entry that was attached to the Latest Checkpoint Entry.
		if latestCheckpointEntry.Method == model.METHOD_COLLECTIONS_WRITE {

			if latestCheckpointEntry.CreatedDate.Equal(params.DateCreated) ||
				params.DateCreated.After(latestCheckpointEntry.CreatedDate) {

				// TODO:  How to compare lexicographically?  Will come back to this

				latestEntry := storage.MessageEntry{
					ID:              primitive.NewObjectID(),
					MessageEntryID:  entryId,
					ParentEntryID:   params.ParentID,
					RecordID:        params.RecordID,
					Schema:          params.Schema,
					Method:          params.Method,
					Data:            []byte(params.Data),
					DataCID:         params.DataCID,
					Protocol:        params.Protocol,
					ProtocolVersion: params.ProtocolVersion,
					CreatedDate:     params.DateCreated,
					PublishedDate:   params.DatePublished,
				}

				err := collectionStore.AddCollectionMessageEntry(&latestEntry)
				if err != nil {
					return err
				}

				err = collectionStore.DeleteCollectionMessageEntry(latestCheckpointEntry)
				if err != nil {
					return err
				}

				parentCollRec.LatestCheckpointEntryID = entryId
				err = collectionStore.SaveCollectionRecord(parentCollRec)
				if err != nil {
					return err
				}

			}

		}
	}

	return errors.New("Supports initial entry only for now")

}

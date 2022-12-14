package collection

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ERR_DUPLICATE_INITIAL_ENTRY                     = "Trying to write an Initial Entry to a record that already exists"
	ERR_COMMIT_TO_RECORD_NOT_FOUND                  = "Trying to commit to a record that does not yet exist"
	ERR_COMMIT_TO_RECORD_CHECKPOINT_ENTRY_NOT_FOUND = "Trying to commit to a record that does not have a latest checkpoint"

	ERR_MUTATE_UMMUTABLE_VALUE     = "Attempt to mutate an immutable value"
	ERR_MISMATCHED_COMMIT_STRATEGY = "Commit Strategy value in existing checkpoint record does not match the commitStrategy value specified in the inbound message,"

	ERR_COMMIT_MESSAGE_CREATE_DATE_BEFORE_WRITE = "Commit message's created date is before the latest write message."
)

type StoreCollectionResult struct {
	Status   string
	RecordID string
	Error    error
}

func StoreCollection(collectionStore storage.CollectionStore, collectionMessage *model.Message) (*StoreCollectionResult, error) {

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing
	result := StoreCollectionResult{}
	switch collectionMessage.Descriptor.Method {

	case model.METHOD_COLLECTIONS_WRITE:
		err := collectionsWrite(collectionStore, collectionMessage)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err
			return &result, nil
		}
		result.Status = "OK"
		result.RecordID = collectionMessage.RecordID

	case model.METHOD_COLLECTIONS_COMMIT:
		err := collectionsCommit(collectionStore, collectionMessage)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err
			return &result, nil
		}
		result.Status = "OK"
		result.RecordID = collectionMessage.RecordID

	case model.METHOD_COLLECTIONS_DELETE:

	default:
		result.Status = "UNSUPPORTED_METHOD"

	}

	return &result, nil
}

func collectionsCommit(collectionsStore storage.CollectionStore, collectionsCommitMessage *model.Message) error {

	// Retrieve the currently active CollectionsWrite entry for the recordId specified in the inbound CollectionsCommit
	// message. If there is no currently active CollectionsWrite entry, discard the inbound message and cease processing.
	existingCollRecord := collectionsStore.GetCollectionRecord(collectionsCommitMessage.Descriptor.ParentID)
	if existingCollRecord == nil {
		return errors.New(ERR_COMMIT_TO_RECORD_NOT_FOUND)
	}

	latestCheckpointEntry := collectionsStore.GetMessageEntryByID(existingCollRecord.LatestCheckpointEntryID)
	if latestCheckpointEntry == nil {
		return errors.New(ERR_COMMIT_TO_RECORD_CHECKPOINT_ENTRY_NOT_FOUND)
	}

	// Ensure all immutable values from the Initial Entry remained unchanged if present in the inbound message.
	// If any have been mutated, discard the message and cease processing.
	// Ensure all immutable values from the Initial Entry remained unchanged if present in the
	// inbound message. If any have been mutated, discard the message and cease processing.
	initialMessageEntry := collectionsStore.GetMessageEntryByID(existingCollRecord.InitialEntryID)
	if initialMessageEntry == nil {
		return errors.New("Unable to find an initial entry")
	}
	if initialMessageEntry.Descriptor.Protocol != collectionsCommitMessage.Descriptor.Protocol ||
		initialMessageEntry.Descriptor.ProtocolVersion != collectionsCommitMessage.Descriptor.ProtocolVersion ||
		initialMessageEntry.Descriptor.Schema != collectionsCommitMessage.Descriptor.Schema {
		return errors.New("Attempt to mutate an immutable value")
	}

	// If the currently active CollectionsWrite does not have a commitStrategy value, or the value does not
	// match the commitStrategy value specified in the inbound message, discard the message and cease processing.
	if latestCheckpointEntry.Descriptor.CommitStrategy != collectionsCommitMessage.Descriptor.CommitStrategy {
		return errors.New(ERR_MISMATCHED_COMMIT_STRATEGY)
	}

	// he parentId of the message MUST match the currently active CollectionsWrite message???s Entry ID or that of
	// another CollectionsCommit that descends from it. If the parentId does not match any of the messages in the
	// commit tree, discard the inbound message and cease processing.
	// This is done by way of searching for the parent ID

	// The inbound message???s entry dateCreated value is less than the dateCreated value of the message in the commit
	// tree its parentId references, discard the message and cease processing.
	if collectionsCommitMessage.Descriptor.DateCreated.Before(latestCheckpointEntry.Descriptor.DateCreated) {
		return errors.New(ERR_COMMIT_MESSAGE_CREATE_DATE_BEFORE_WRITE)
	}

	// If all of the above steps are successful, store the message in relation to the record.
	// I think this means we set the latest checkpoint to the latest entry
	existingCollRecord.LatestEntryID = existingCollRecord.LatestCheckpointEntryID
	collectionsStore.SaveCollectionRecord(existingCollRecord)

	return nil

}

func collectionsWrite(collectionStore storage.CollectionStore, collectionsWriteMessage *model.Message) error {

	/*
			Generate the message???s Entry ID by performing the Record ID Generation Process.
			- IF the generated Entry ID matches the recordId value of the message it is the Initial Entry
		      for a record, store the entry as the Initial Entry for the record if no Initial Entry exists
		      and cease any further processing.
			- ELSE the message may be an overwriting entry for the record; continue processing.
	*/
	descriptorId := model.CreateDescriptorCID(collectionsWriteMessage.Descriptor)
	processingId := model.CreateProcessingCID(collectionsWriteMessage.Processing)
	entryId := model.CreateRecordCID(descriptorId, processingId)

	//log.Info(" Entry ID: %s", entryId)
	//log.Info("Record ID: %s", collectionsWriteMessage.RecordID)

	// TODO:  Come back to this.  This should match
	// if entryId == collectionsWriteMessage.RecordID {
	// For now:  If there is no parent ID, ASSUME first entry
	if collectionsWriteMessage.Descriptor.ParentID == "" {

		// This is the first entry of the record.  Create it and return
		// If there is an existing record id, there's a problem and return an error
		existingRecord := collectionStore.GetCollectionRecord(entryId)
		if existingRecord != nil {
			return errors.New(ERR_DUPLICATE_INITIAL_ENTRY)
		}

		record := storage.CollectionRecord{
			ID:         primitive.NewObjectID(),
			RecordID:   collectionsWriteMessage.RecordID,
			CreatorDID: collectionsWriteMessage.Processing.AuthorDID,
			OwnerDID:   collectionsWriteMessage.Processing.RecipientDID,
			WriterDIDs: []string{collectionsWriteMessage.Processing.AuthorDID},
			ReaderDIDs: []string{collectionsWriteMessage.Processing.AuthorDID, collectionsWriteMessage.Processing.RecipientDID},
		}

		entry := storage.MessageEntry{
			ID:             primitive.NewObjectID(),
			MessageEntryID: entryId,
			Message:        *collectionsWriteMessage,
		}

		err := collectionStore.CreateCollectionRecord(&record, &entry)
		return err

	} else {

		// the message may be an overwriting entry for the record; continue processing.
		// This is an attempt to overwrite a previous version.
		// So, let's get the parent version
		parentCollRec := collectionStore.GetCollectionRecord(collectionsWriteMessage.Descriptor.ParentID)
		if parentCollRec == nil {
			// If a message is not the Initial Entry, its descriptor MUST contain a parentId to
			// determine the entry???s position in the record???s lineage. If a parentId is present
			// proceed with processing, else discard the record and cease processing.
			// We dont have the parent.  Reject with err
			return fmt.Errorf("Unable to find Parent Record for Overwrite using Parent ID:  %s", collectionsWriteMessage.Descriptor.ParentID)
		}

		// Ensure all immutable values from the Initial Entry remained unchanged if present in the
		// inbound message. If any have been mutated, discard the message and cease processing.
		initialMessageEntry := collectionStore.GetMessageEntryByID(parentCollRec.InitialEntryID)
		if initialMessageEntry == nil {
			return errors.New("Unable to find an initial entry")
		}
		if initialMessageEntry.Descriptor.Protocol != collectionsWriteMessage.Descriptor.Protocol ||
			initialMessageEntry.Descriptor.ProtocolVersion != collectionsWriteMessage.Descriptor.ProtocolVersion ||
			initialMessageEntry.Descriptor.Schema != collectionsWriteMessage.Descriptor.Schema {
			return errors.New(ERR_MUTATE_UMMUTABLE_VALUE)
		}

		// Retrieve the Latest Checkpoint Entry, which will be either the Initial Entry or the latest CollectionsDelete,
		// and compare the parentId value of the inbound message to the Entry ID of the
		// Latest Checkpoint Entry derived from running the Record ID Generation Process on it.
		// If the values match, proceed with processing, if the values do not match discard the message and cease processing.
		latestCheckpointEntry := collectionStore.GetMessageEntryByID(parentCollRec.LatestCheckpointEntryID)
		if latestCheckpointEntry == nil {
			return errors.New("Unable to find the latest checkpoint entry")
		}

		if collectionsWriteMessage.Descriptor.ParentID != latestCheckpointEntry.RecordID {
			return errors.New("The parent ID of the inbound message must match the latest checkpoint record ID.")
		}

		// If an existing CollectionsWrite entry linked to the Latest Checkpoint Entry IS NOT present and
		// the dateCreated value of the inbound message is greater than the Latest Checkpoint Entry,
		// store the message as the Latest Entry and cease processing, else discard the inbound message
		// and cease processing.
		if latestCheckpointEntry.Descriptor.Method != model.METHOD_COLLECTIONS_WRITE &&
			collectionsWriteMessage.Descriptor.DateCreated.After(latestCheckpointEntry.Descriptor.DateCreated) {

			latestEntry := storage.MessageEntry{
				ID:             primitive.NewObjectID(),
				MessageEntryID: uuid.NewString(),
				Message:        *collectionsWriteMessage,
			}

			err := collectionStore.AddCollectionMessageEntry(&latestEntry)
			if err != nil {
				return err
			}

			parentCollRec.LatestEntryID = latestEntry.MessageEntryID
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
		if latestCheckpointEntry.Descriptor.Method == model.METHOD_COLLECTIONS_WRITE {

			if latestCheckpointEntry.Descriptor.DateCreated.Equal(collectionsWriteMessage.Descriptor.DateCreated) ||
				collectionsWriteMessage.Descriptor.DateCreated.After(latestCheckpointEntry.Descriptor.DateCreated) {

				// TODO:  How to compare lexicographically?  Will come back to this

				latestEntry := storage.MessageEntry{
					ID:             primitive.NewObjectID(),
					MessageEntryID: uuid.NewString(),
					Message:        *collectionsWriteMessage,
				}

				err := collectionStore.AddCollectionMessageEntry(&latestEntry)
				if err != nil {
					return err
				}

				// I don't believe ww want to delete the message entry
				//err = collectionStore.DeleteCollectionMessageEntry(latestCheckpointEntry)
				//if err != nil {
				//	return err
				//}

				parentCollRec.LatestCheckpointEntryID = latestEntry.MessageEntryID
				err = collectionStore.SaveCollectionRecord(parentCollRec)
				if err != nil {
					return err
				}

			}

		}
	}

	return nil

}

package record

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	ERR_DUPLICATE_INITIAL_ENTRY                     = "Trying to write an Initial Entry to a record that already exists"
	ERR_COMMIT_TO_RECORD_NOT_FOUND                  = "Trying to commit to a record that does not yet exist"
	ERR_COMMIT_TO_RECORD_CHECKPOINT_ENTRY_NOT_FOUND = "Trying to commit to a record that does not have a latest checkpoint"

	ERR_MUTATE_UMMUTABLE_VALUE     = "Attempt to mutate an immutable value"
	ERR_MISMATCHED_COMMIT_STRATEGY = "Commit Strategy value in existing checkpoint record does not match the commitStrategy value specified in the inbound message,"

	ERR_COMMIT_MESSAGE_CREATE_DATE_BEFORE_WRITE = "Commit message's created date is before the latest write message."
	ERR_INVALID_DATE_FORMAT                     = "Invalid date format, you must use RFC3339 format"
)

type StoreRecordResult struct {
	Status       string
	RecordID     string
	InitialEntry bool
	Error        error
}

func StoreRecord(ctx context.Context, recordStore storage.RecordStore, recordMessage *model.Message) (*StoreRecordResult, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "StoreRecord")
	defer sp.End()

	// Need to implement this message process flow per spec:
	// https://identity.foundation/decentralized-web-node/spec/#retained-message-processing
	result := StoreRecordResult{}
	switch recordMessage.Descriptor.Method {

	case model.METHOD_RECORDS_WRITE:
		initialEntry, err := recordWrite(ctx, recordStore, recordMessage)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err
			return &result, nil
		}
		result.Status = "OK"
		result.RecordID = recordMessage.RecordID
		result.InitialEntry = initialEntry
		sp.AddEvent(fmt.Sprintf("Initial Entry?  %v", initialEntry))

	case model.METHOD_RECORDS_COMMIT:
		err := recordCommit(ctx, recordStore, recordMessage)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err
			return &result, nil
		}
		result.Status = "OK"
		result.RecordID = recordMessage.RecordID
		result.InitialEntry = false
		sp.AddEvent(fmt.Sprintf("Initial Entry?  %v", false))

	case model.METHOD_RECORDS_DELETE:

	default:
		result.Status = "UNSUPPORTED_METHOD"

	}

	return &result, nil
}

func recordCommit(ctx context.Context, recordStore storage.RecordStore, recordCommitMessage *model.Message) error {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordCommit")
	defer sp.End()

	// Retrieve the currently active CollectionsWrite entry for the recordId specified in the inbound CollectionsCommit
	// message. If there is no currently active CollectionsWrite entry, discard the inbound message and cease processing.
	// The parentID in the case of a COMMIT is the specific record ID in the chain of messages.
	// We need to retrieve THAT entry, ensure its the latest checkpoint entry and then "commit"
	sp.AddEvent("Get record for commit from parent ID")
	existingRecord, parentWrite := recordStore.GetRecordForCommit(recordCommitMessage.Descriptor.ParentID)
	if existingRecord == nil || parentWrite == nil {
		return errors.New(ERR_COMMIT_TO_RECORD_NOT_FOUND)
	}

	sp.AddEvent("Get Message Entry from existing record")
	latestCheckpointEntry := recordStore.GetMessageEntryByID(existingRecord.LatestCheckpointEntryID)
	if latestCheckpointEntry == nil {
		return errors.New(ERR_COMMIT_TO_RECORD_CHECKPOINT_ENTRY_NOT_FOUND)
	}

	// Ensure all immutable values from the Initial Entry remained unchanged if present in the inbound message.
	// If any have been mutated, discard the message and cease processing.
	// Ensure all immutable values from the Initial Entry remained unchanged if present in the
	// inbound message. If any have been mutated, discard the message and cease processing.
	sp.AddEvent("Get the initial entry from the existing record")
	initialMessageEntry := recordStore.GetMessageEntryByID(existingRecord.InitialEntryID)
	if initialMessageEntry == nil {
		return errors.New("Unable to find an initial entry")
	}

	sp.AddEvent("Checking for immutables")
	if initialMessageEntry.Descriptor.Protocol != recordCommitMessage.Descriptor.Protocol ||
		initialMessageEntry.Descriptor.ProtocolVersion != recordCommitMessage.Descriptor.ProtocolVersion ||
		initialMessageEntry.Descriptor.Schema != recordCommitMessage.Descriptor.Schema {
		return errors.New("Attempt to mutate an immutable value")
	}

	// If the currently active CollectionsWrite does not have a commitStrategy value, or the value does not
	// match the commitStrategy value specified in the inbound message, discard the message and cease processing.
	if latestCheckpointEntry.Descriptor.CommitStrategy != recordCommitMessage.Descriptor.CommitStrategy {
		return errors.New(ERR_MISMATCHED_COMMIT_STRATEGY)
	}

	// The parentId of the message MUST match the currently active CollectionsWrite message’s Entry ID or that of
	// another CollectionsCommit that descends from it. If the parentId does not match any of the messages in the
	// commit tree, discard the inbound message and cease processing.
	// This is done by way of searching for the parent ID

	// The inbound message’s entry dateCreated value is less than the dateCreated value of the message in the commit
	// tree its parentId references, discard the message and cease processing.
	recordCommitMessageDateCreated, err := time.Parse(time.RFC3339, recordCommitMessage.Descriptor.DateCreated)
	latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointEntry.Descriptor.DateCreated)
	if err != nil {
		return errors.New(ERR_INVALID_DATE_FORMAT)
	}

	if recordCommitMessageDateCreated.Before(latestCheckpointEntryDateCreated) {
		return errors.New(ERR_COMMIT_MESSAGE_CREATE_DATE_BEFORE_WRITE)
	}

	sp.AddEvent("Adding Message Entry")
	commitMessageEntry := storage.MessageEntry{
		Message:        *recordCommitMessage,
		MessageEntryID: uuid.NewString(),
	}
	err = recordStore.AddMessageEntry(&commitMessageEntry)
	if err != nil {
		sp.RecordError(err)
		return err
	}

	// If all of the above steps are successful, store the message in relation to the record.
	// I think this means we set the latest checkpoint to the latest entry
	sp.AddEvent("Updating existing record")
	existingRecord.LatestEntryID = existingRecord.LatestCheckpointEntryID
	existingRecord.LatestCheckpointEntryID = commitMessageEntry.MessageEntryID
	recordStore.SaveRecord(existingRecord)

	return nil

}

func recordWrite(ctx context.Context, recordStore storage.RecordStore, recordMessage *model.Message) (bool, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordWrite")
	defer sp.End()

	/*
			Generate the message’s Entry ID by performing the Record ID Generation Process.
			- IF the generated Entry ID matches the recordId value of the message it is the Initial Entry
		      for a record, store the entry as the Initial Entry for the record if no Initial Entry exists
		      and cease any further processing.
			- ELSE the message may be an overwriting entry for the record; continue processing.
	*/
	descriptorId := model.CreateDescriptorCID(recordMessage.Descriptor)
	processingId := model.CreateProcessingCID(recordMessage.Processing)
	entryId := model.CreateRecordCID(descriptorId, processingId)

	//log.Info(" Entry ID: %s", entryId)
	//log.Info("Record ID: %s", recordMessage.RecordID)

	// TODO:  Come back to this.  This should match
	// if entryId == recordMessage.RecordID {
	// For now:  If there is no parent ID, ASSUME first entry
	sp.AddEvent("Is this the first entry?")
	if recordMessage.Descriptor.ParentID == "" {

		sp.AddEvent("THIS IS THE FIRST ENTRY")
		// This is the first entry of the record.  Create it and return
		// If there is an existing record id, there's a problem and return an error
		// Let caller know this is an intial entry via boolean
		existingRecord := recordStore.GetRecord(entryId)
		if existingRecord != nil {
			sp.RecordError(errors.New(ERR_DUPLICATE_INITIAL_ENTRY))
			return false, errors.New(ERR_DUPLICATE_INITIAL_ENTRY)
		}

		record := storage.Record{
			ID:         primitive.NewObjectID(),
			RecordID:   recordMessage.RecordID,
			CreatorDID: recordMessage.Processing.AuthorDID,
			OwnerDID:   recordMessage.Processing.RecipientDID,
			WriterDIDs: []string{recordMessage.Processing.AuthorDID},
			ReaderDIDs: []string{recordMessage.Processing.AuthorDID, recordMessage.Processing.RecipientDID},
		}

		entry := storage.MessageEntry{
			ID:             primitive.NewObjectID(),
			MessageEntryID: entryId,
			Message:        *recordMessage,
		}

		sp.AddEvent("Creating record")
		err := recordStore.CreateRecord(&record, &entry)
		if err != nil {
			sp.RecordError(err)
			return false, err
		}

		// This was an initial entry!
		sp.AddEvent("CREATED AN INITIAL ENTRY")
		return true, nil

	} else {

		sp.AddEvent("NOT an initial entry")
		// the message may be an overwriting entry for the record; continue processing.
		// This is an attempt to overwrite a previous version.
		// So, let's get the parent version
		sp.AddEvent("Get the parent  record for write referenced in the parentID")
		parentRecord := recordStore.GetRecord(recordMessage.Descriptor.ParentID)

		if parentRecord == nil {
			// If a message is not the Initial Entry, its descriptor MUST contain a parentId to
			// determine the entry’s position in the record’s lineage. If a parentId is present
			// proceed with processing, else discard the record and cease processing.
			// We dont have the parent.  Reject with err
			sp.RecordError(errors.New("Unable to find Parent Records for overwrite using parent ID"))
			return false, fmt.Errorf("Unable to find Parent Record for Overwrite using Parent ID:  %s", recordMessage.Descriptor.ParentID)
		}

		// Ensure all immutable values from the Initial Entry remained unchanged if present in the
		// inbound message. If any have been mutated, discard the message and cease processing.
		sp.AddEvent("Get the parent's Initial Entry ID")
		initialMessageEntry := recordStore.GetMessageEntryByID(parentRecord.InitialEntryID)
		if initialMessageEntry == nil {
			return false, errors.New("Unable to find an initial entry")
		}

		sp.AddEvent("Ensuring immutable values stay immutable")
		if initialMessageEntry.Descriptor.Protocol != recordMessage.Descriptor.Protocol ||
			initialMessageEntry.Descriptor.ProtocolVersion != recordMessage.Descriptor.ProtocolVersion ||
			initialMessageEntry.Descriptor.Schema != recordMessage.Descriptor.Schema {
			return false, errors.New(ERR_MUTATE_UMMUTABLE_VALUE)
		}

		// Retrieve the Latest Checkpoint Entry, which will be either the Initial Entry or the latest CollectionsDelete,
		// and compare the parentId value of the inbound message to the Entry ID of the
		// Latest Checkpoint Entry derived from running the Record ID Generation Process on it.
		// If the values match, proceed with processing, if the values do not match discard the message and cease processing.
		sp.AddEvent("Getting the latest checkpoint entry")
		latestCheckpointEntry := recordStore.GetMessageEntryByID(parentRecord.LatestCheckpointEntryID)
		if latestCheckpointEntry == nil {
			return false, errors.New("Unable to find the latest checkpoint entry")
		}

		if recordMessage.Descriptor.ParentID != latestCheckpointEntry.RecordID {
			return false, errors.New("The parent ID of the inbound message must match the latest checkpoint record ID.")
		}

		// If an existing CollectionsWrite entry linked to the Latest Checkpoint Entry IS NOT present and
		// the dateCreated value of the inbound message is greater than the Latest Checkpoint Entry,
		// store the message as the Latest Entry and cease processing, else discard the inbound message
		// and cease processing.
		recordMessageDateCreated, err := time.Parse(time.RFC3339, recordMessage.Descriptor.DateCreated)
		latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointEntry.Descriptor.DateCreated)
		if err != nil {
			return false, errors.New(ERR_INVALID_DATE_FORMAT)
		}
		if latestCheckpointEntry.Descriptor.Method != model.METHOD_RECORDS_WRITE &&
			recordMessageDateCreated.After(latestCheckpointEntryDateCreated) {

			sp.AddEvent("Storing latest entry, adjusting latest entry ID to this WRITE")
			latestEntry := storage.MessageEntry{
				ID:             primitive.NewObjectID(),
				MessageEntryID: uuid.NewString(),
				Message:        *recordMessage,
			}

			err := recordStore.AddMessageEntry(&latestEntry)
			if err != nil {
				return false, err
			}

			parentRecord.LatestEntryID = latestEntry.MessageEntryID
			err = recordStore.SaveRecord(parentRecord)
			if err != nil {
				return false, err
			}

		}

		// If an exiting CollectionsWrite entry linked to the Latest Checkpoint Entry IS present
		// all of the following conditions MUST be true:
		//   - The dateCreated value of the inbound message is greater than the existing CollectionsWrite,
		//     or if the dateCreated values are the same, the Entry ID of the inbound message is greater
		//     than the existing entry when the Entry IDs of the two are compared lexicographically.
		// If all of the following conditions for Step 6 are true, store the inbound message as the Latest Entry
		// and discard the existing CollectionsWrite entry that was attached to the Latest Checkpoint Entry.
		if latestCheckpointEntry.Descriptor.Method == model.METHOD_RECORDS_WRITE {

			recordMessageDateCreated, err := time.Parse(time.RFC3339, recordMessage.Descriptor.DateCreated)
			latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointEntry.Descriptor.DateCreated)
			if err != nil {
				return false, errors.New(ERR_INVALID_DATE_FORMAT)
			}
			if latestCheckpointEntryDateCreated.Equal(recordMessageDateCreated) ||
				recordMessageDateCreated.After(latestCheckpointEntryDateCreated) {

				// TODO:  How to compare lexicographically?  Will come back to this
				sp.AddEvent("Storing latest entry OVERRIDE previous WRITE, adjusting latest entry ID to this WRITE")
				latestEntry := storage.MessageEntry{
					ID:             primitive.NewObjectID(),
					MessageEntryID: uuid.NewString(),
					Message:        *recordMessage,
				}

				err := recordStore.AddMessageEntry(&latestEntry)
				if err != nil {
					return false, err
				}

				// I don't believe ww want to delete the message entry
				//err = recordStore.DeleteCollectionMessageEntry(latestCheckpointEntry)
				//if err != nil {
				//	return err
				//}

				parentRecord.LatestCheckpointEntryID = latestEntry.MessageEntryID
				err = recordStore.SaveRecord(parentRecord)
				if err != nil {
					return false, err
				}

			}

		}
	}

	return false, nil

}

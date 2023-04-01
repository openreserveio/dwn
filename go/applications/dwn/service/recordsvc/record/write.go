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

func RecordWrite(ctx context.Context, recordStore storage.RecordStore, recordMessage *model.Message) (string, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordsvc.record.RecordWrite")
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
			return "", errors.New(ERR_DUPLICATE_INITIAL_ENTRY)
		}

		record := storage.Record{
			ID:         primitive.NewObjectID(),
			RecordID:   recordMessage.RecordID,
			CreatorDID: recordMessage.Processing.AuthorDID,
			OwnerDID:   recordMessage.Processing.RecipientDID,
			WriterDIDs: []string{recordMessage.Processing.RecipientDID},
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
			return "", err
		}

		// This was an initial entry!
		sp.AddEvent("CREATED AN INITIAL ENTRY")
		return record.RecordID, nil

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
			return "", fmt.Errorf("Unable to find Parent Record for Overwrite using Parent ID:  %s", recordMessage.Descriptor.ParentID)
		}

		// Ensure all immutable values from the Initial Entry remained unchanged if present in the
		// inbound message. If any have been mutated, discard the message and cease processing.
		sp.AddEvent("Get the parent's Initial Entry ID")
		initialMessageEntry := recordStore.GetMessageEntryByID(parentRecord.InitialEntryID)
		if initialMessageEntry == nil {
			return "", errors.New("Unable to find an initial entry")
		}

		sp.AddEvent("Ensuring immutable values stay immutable")
		if initialMessageEntry.Descriptor.Protocol != recordMessage.Descriptor.Protocol ||
			initialMessageEntry.Descriptor.ProtocolVersion != recordMessage.Descriptor.ProtocolVersion ||
			initialMessageEntry.Descriptor.Schema != recordMessage.Descriptor.Schema {
			return "", errors.New(ERR_MUTATE_UMMUTABLE_VALUE)
		}

		// Retrieve the Latest Checkpoint Entry, which will be either the Initial Entry or the latest CollectionsDelete,
		// and compare the parentId value of the inbound message to the Entry ID of the
		// Latest Checkpoint Entry derived from running the Record ID Generation Process on it.
		// If the values match, proceed with processing, if the values do not match discard the message and cease processing.
		sp.AddEvent("Getting the latest checkpoint entry")
		latestCheckpointEntry := recordStore.GetMessageEntryByID(parentRecord.LatestCheckpointEntryID)
		if latestCheckpointEntry == nil {
			return "", errors.New("Unable to find the latest checkpoint entry")
		}

		if recordMessage.Descriptor.ParentID != latestCheckpointEntry.RecordID {
			return "", errors.New("The parent ID of the inbound message must match the latest checkpoint record ID.")
		}

		// If an existing CollectionsWrite entry linked to the Latest Checkpoint Entry IS NOT present and
		// the dateCreated value of the inbound message is greater than the Latest Checkpoint Entry,
		// store the message as the Latest Entry and cease processing, else discard the inbound message
		// and cease processing.
		recordMessageDateCreated, err := time.Parse(time.RFC3339, recordMessage.Descriptor.DateCreated)
		latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointEntry.Descriptor.DateCreated)
		if err != nil {
			return "", errors.New(ERR_INVALID_DATE_FORMAT)
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
				return "", err
			}

			parentRecord.LatestEntryID = latestEntry.MessageEntryID
			err = recordStore.SaveRecord(parentRecord)
			if err != nil {
				return "", err
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
				return "", errors.New(ERR_INVALID_DATE_FORMAT)
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
					return "", err
				}

				// I don't believe ww want to delete the message entry
				//err = recordStore.DeleteCollectionMessageEntry(latestCheckpointEntry)
				//if err != nil {
				//	return err
				//}

				parentRecord.LatestCheckpointEntryID = latestEntry.MessageEntryID
				err = recordStore.SaveRecord(parentRecord)
				if err != nil {
					return "", err
				}

			}

		}
	}

	return recordMessage.RecordID, nil

}

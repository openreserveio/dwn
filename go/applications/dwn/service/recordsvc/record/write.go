package record

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"time"
)

func RecordWrite(ctx context.Context, recordStore storage.RecordStore, recordMessage *model.Message) (string, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "recordsvc.record.RecordWrite")
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
	sp.AddEvent(fmt.Sprintf("Generated Message Entry ID, but kinda useless:  %s", entryId))

	//log.Info(" Entry ID: %s", entryId)
	//log.Info("Record ID: %s", recordMessage.RecordID)

	// TODO:  Come back to this.  This should match
	// if entryId == recordMessage.RecordID {
	// For now:  If there is no parent ID, ASSUME first entry
	sp.AddEvent("Is this the first entry?")
	var initialEntry bool = false
	existingRecord := recordStore.GetRecord(ctx, recordMessage.RecordID)
	if existingRecord == nil {
		initialEntry = true
	}

	if initialEntry {

		sp.AddEvent("THIS IS THE FIRST ENTRY")
		// This is the first entry of the record.  Create it and return
		// If there is an existing record id, there's a problem and return an error
		// Let caller know this is an intial entry via boolean
		record := storage.Record{
			DWNRecordID: recordMessage.RecordID,
			CreatorDID:  recordMessage.Processing.AuthorDID,
			OwnerDID:    recordMessage.Processing.RecipientDID,
			WriterDIDs:  []string{recordMessage.Processing.RecipientDID},
			ReaderDIDs:  []string{recordMessage.Processing.AuthorDID, recordMessage.Processing.RecipientDID},
		}

		recordMessageBytes, err := json.Marshal(recordMessage)
		entry := storage.MessageEntry{
			DWNRecordID: recordMessage.RecordID,
			Message:     recordMessageBytes,
		}

		sp.AddEvent("Creating record")
		recordStore.BeginTx(ctx)
		err = recordStore.CreateRecord(ctx, &record, &entry)
		if err != nil {
			recordStore.RollbackTx(ctx)
			sp.RecordError(err)
			return "", err
		}
		recordStore.CommitTx(ctx)

		// This was an initial entry!
		sp.AddEvent("CREATED AN INITIAL ENTRY")
		return record.DWNRecordID, nil

	} else {

		sp.AddEvent("Adding a WRITE to existing record")
		// the message may be an overwriting entry for the record; continue processing.
		// This is an attempt to overwrite a previous version.
		// So, let's get the parent version
		sp.AddEvent("Get the parent  record for write referenced in the parentID")

		// Ensure all immutable values from the Initial Entry remained unchanged if present in the
		// inbound message. If any have been mutated, discard the message and cease processing.
		sp.AddEvent("Get the parent's Initial Entry ID")
		initialMessageEntry := recordStore.GetMessageEntryByID(ctx, existingRecord.InitialEntryID)
		if initialMessageEntry == nil {
			return "", errors.New("Unable to find an initial entry")
		}

		sp.AddEvent("Ensuring immutable values stay immutable")
		var initialMessage model.Message
		err := json.Unmarshal(initialMessageEntry.Message, &initialMessage)

		if initialMessage.Descriptor.Protocol != recordMessage.Descriptor.Protocol ||
			initialMessage.Descriptor.ProtocolVersion != recordMessage.Descriptor.ProtocolVersion ||
			initialMessage.Descriptor.Schema != recordMessage.Descriptor.Schema {
			return "", errors.New(ERR_MUTATE_UMMUTABLE_VALUE)
		}

		// Retrieve the Latest Checkpoint Entry, which will be either the Initial Entry or the latest CollectionsDelete,
		// and compare the parentId value of the inbound message to the Entry ID of the
		// Latest Checkpoint Entry derived from running the Record ID Generation Process on it.
		// If the values match, proceed with processing, if the values do not match discard the message and cease processing.
		sp.AddEvent("Getting the latest checkpoint entry")
		latestCheckpointEntry := recordStore.GetMessageEntryByID(ctx, existingRecord.LatestCheckpointEntryID)
		if latestCheckpointEntry == nil {
			return "", errors.New("Unable to find the latest checkpoint entry")
		}

		// If an existing CollectionsWrite entry linked to the Latest Checkpoint Entry IS NOT present and
		// the dateCreated value of the inbound message is greater than the Latest Checkpoint Entry,
		// store the message as the Latest Entry and cease processing, else discard the inbound message
		// and cease processing.
		var latestCheckpointMessage model.Message
		err = json.Unmarshal(latestCheckpointEntry.Message, &latestCheckpointMessage)

		recordMessageDateCreated, err := time.Parse(time.RFC3339, recordMessage.Descriptor.DateCreated)
		latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointMessage.Descriptor.DateCreated)
		if err != nil {
			return "", errors.New(ERR_INVALID_DATE_FORMAT)
		}

		if latestCheckpointMessage.Descriptor.Method != model.METHOD_RECORDS_WRITE &&
			recordMessageDateCreated.After(latestCheckpointEntryDateCreated) {

			sp.AddEvent("Storing latest entry, adjusting latest entry ID to this WRITE")
			recordMessageBytes, err := json.Marshal(recordMessage)
			latestEntry := storage.MessageEntry{
				ID:                     uuid.NewString(),
				PreviousMessageEntryID: latestCheckpointEntry.ID,
				DWNRecordID:            initialMessage.RecordID,
				Message:                recordMessageBytes,
			}

			err = recordStore.AddMessageEntry(ctx, &latestEntry)
			if err != nil {
				return "", err
			}

			existingRecord.LatestEntryID = latestEntry.ID
			err = recordStore.SaveRecord(ctx, existingRecord)
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
		if latestCheckpointMessage.Descriptor.Method == model.METHOD_RECORDS_WRITE {

			recordMessageDateCreated, err := time.Parse(time.RFC3339, recordMessage.Descriptor.DateCreated)
			latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointMessage.Descriptor.DateCreated)
			if err != nil {
				return "", errors.New(ERR_INVALID_DATE_FORMAT)
			}
			if latestCheckpointEntryDateCreated.Equal(recordMessageDateCreated) ||
				recordMessageDateCreated.After(latestCheckpointEntryDateCreated) {

				// TODO:  How to compare lexicographically?  Will come back to this
				sp.AddEvent("Storing latest entry OVERRIDE previous WRITE, adjusting latest entry ID to this WRITE")
				recordMessageBytes, err := json.Marshal(recordMessage)
				latestEntry := storage.MessageEntry{
					ID:          uuid.NewString(),
					DWNRecordID: existingRecord.DWNRecordID,
					Message:     recordMessageBytes,
				}

				err = recordStore.AddMessageEntry(ctx, &latestEntry)
				if err != nil {
					return "", err
				}

				// I don't believe ww want to delete the message entry
				//err = recordStore.DeleteCollectionMessageEntry(latestCheckpointEntry)
				//if err != nil {
				//	return err
				//}

				existingRecord.LatestCheckpointEntryID = latestEntry.ID
				err = recordStore.SaveRecord(ctx, existingRecord)
				if err != nil {
					return "", err
				}

			}

		}
	}

	return recordMessage.RecordID, nil

}

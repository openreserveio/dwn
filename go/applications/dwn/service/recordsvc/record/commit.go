package record

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
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

func RecordCommit(ctx context.Context, recordStore storage.RecordStore, recordCommitMessage *model.Message) error {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "recordsvc.record.RecordCommit")
	defer sp.End()

	// Retrieve the currently active RecordsWrite entry for the recordId specified in the inbound RecordsCommit message.
	// If there is no currently active RecordsWrite entry, discard the inbound message and cease processing.

	sp.AddEvent("Get record for commit")
	logicalRecord := recordStore.GetRecord(ctx, recordCommitMessage.RecordID)
	if logicalRecord == nil {
		err := errors.New(ERR_COMMIT_TO_RECORD_NOT_FOUND)
		sp.RecordError(err)
		return err
	}

	sp.AddEvent("Get Message Entry from existing record")
	latestCheckpointEntry := recordStore.GetMessageEntryByID(ctx, logicalRecord.LatestCheckpointEntryID)
	if latestCheckpointEntry == nil {
		return errors.New(ERR_COMMIT_TO_RECORD_CHECKPOINT_ENTRY_NOT_FOUND)
	}

	// Ensure all immutable values from the Initial Entry remained unchanged if present in the inbound message.
	// If any have been mutated, discard the message and cease processing.
	// Ensure all immutable values from the Initial Entry remained unchanged if present in the
	// inbound message. If any have been mutated, discard the message and cease processing.
	sp.AddEvent("Get the initial entry from the existing record")
	initialMessageEntry := recordStore.GetMessageEntryByID(ctx, logicalRecord.InitialEntryID)
	if initialMessageEntry == nil {
		return errors.New("Unable to find an initial entry")
	}

	sp.AddEvent("Checking for immutables")
	var initialMessage model.Message
	json.Unmarshal(initialMessageEntry.Message, &initialMessage)

	if initialMessage.Descriptor.Protocol != recordCommitMessage.Descriptor.Protocol ||
		initialMessage.Descriptor.ProtocolVersion != recordCommitMessage.Descriptor.ProtocolVersion ||
		initialMessage.Descriptor.Schema != recordCommitMessage.Descriptor.Schema {
		return errors.New("Attempt to mutate an immutable value")
	}

	// If the currently active CollectionsWrite does not have a commitStrategy value, or the value does not
	// match the commitStrategy value specified in the inbound message, discard the message and cease processing.
	var latestCheckpointMessage model.Message
	json.Unmarshal(latestCheckpointEntry.Message, &latestCheckpointMessage)
	if latestCheckpointMessage.Descriptor.CommitStrategy != recordCommitMessage.Descriptor.CommitStrategy {
		return errors.New(ERR_MISMATCHED_COMMIT_STRATEGY)
	}

	// The parentId of the message MUST match the currently active CollectionsWrite message’s Entry ID or that of
	// another CollectionsCommit that descends from it. If the parentId does not match any of the messages in the
	// commit tree, discard the inbound message and cease processing.
	// This is done by way of searching for the parent ID

	// The inbound message’s entry dateCreated value is less than the dateCreated value of the message in the commit
	// tree its parentId references, discard the message and cease processing.
	recordCommitMessageDateCreated, err := time.Parse(time.RFC3339, recordCommitMessage.Descriptor.DateCreated)
	latestCheckpointEntryDateCreated, err := time.Parse(time.RFC3339, latestCheckpointMessage.Descriptor.DateCreated)
	if err != nil {
		return errors.New(ERR_INVALID_DATE_FORMAT)
	}

	if recordCommitMessageDateCreated.Before(latestCheckpointEntryDateCreated) {
		return errors.New(ERR_COMMIT_MESSAGE_CREATE_DATE_BEFORE_WRITE)
	}

	sp.AddEvent("Adding Message Entry")
	messageBytes, err := json.Marshal(recordCommitMessage)
	if err != nil {
		sp.RecordError(err)
		return err
	}

	commitMessageEntry := storage.MessageEntry{
		PreviousMessageEntryID: latestCheckpointEntry.ID,
		DWNRecordID:            recordCommitMessage.RecordID,
		Message:                messageBytes,
	}
	err = recordStore.AddMessageEntry(ctx, &commitMessageEntry)
	if err != nil {
		sp.RecordError(err)
		return err
	}

	// If all of the above steps are successful, store the message in relation to the record.
	// I think this means we set the latest checkpoint to the latest entry
	sp.AddEvent("Updating existing record")
	logicalRecord.LatestEntryID = logicalRecord.LatestCheckpointEntryID
	recordStore.SaveRecord(ctx, logicalRecord)

	return nil

}

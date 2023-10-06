package record

import (
	"context"
	"errors"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
)

type FindRecordResult struct {
	Status                string
	Error                 error
	Record                *storage.Record
	InitialEntry          *storage.MessageEntry
	LatestEntry           *storage.MessageEntry
	LatestCheckpointEntry *storage.MessageEntry
}

func FindRecordBySchemaAndRecordID(ctx context.Context, recordStore storage.RecordStore, schemaUri string, recordId string) (*FindRecordResult, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "recordsvc.record.FindRecordBySchemaAndRecordID")
	defer sp.End()

	collRecord := recordStore.GetRecord(ctx, recordId)
	if collRecord == nil {
		return &FindRecordResult{Status: "NOT_FOUND"}, nil
	}

	return &FindRecordResult{
		Status:                "OK",
		Error:                 nil,
		Record:                collRecord,
		InitialEntry:          recordStore.GetMessageEntryByID(ctx, collRecord.InitialEntryID),
		LatestEntry:           recordStore.GetMessageEntryByID(ctx, collRecord.LatestEntryID),
		LatestCheckpointEntry: recordStore.GetMessageEntryByID(ctx, collRecord.LatestCheckpointEntryID),
	}, nil

}

func FindRecordForCommit(ctx context.Context, recordStore storage.RecordStore, schemaUri string, logicalRecordId string) (*FindRecordResult, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "recordsvc.record.FindRecordForCommit")
	defer sp.End()

	result := FindRecordResult{}

	sp.AddEvent("Calling Record Store to get record for commit")
	record := recordStore.GetRecord(ctx, logicalRecordId)
	if record == nil {
		sp.AddEvent("Unable to find records for Commit")
		result.Status = "NOT_FOUND"
		result.Error = errors.New("Not Found")
		return &result, nil
	}

	messageEntry := recordStore.GetMessageEntryByID(ctx, record.LatestCheckpointEntryID)

	result.Status = "OK"
	result.LatestEntry = messageEntry
	result.Record = record

	return &result, nil

}

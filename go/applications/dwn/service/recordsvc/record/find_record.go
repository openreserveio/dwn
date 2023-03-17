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
	ctx, sp := observability.Tracer.Start(ctx, "FindRecordBySchemaAndRecordID")
	defer sp.End()

	collRecord := recordStore.GetRecord(recordId)
	if collRecord == nil {
		return &FindRecordResult{Status: "NOT_FOUND"}, nil
	}

	return &FindRecordResult{
		Status:                "OK",
		Error:                 nil,
		Record:                collRecord,
		InitialEntry:          recordStore.GetMessageEntryByID(collRecord.InitialEntryID),
		LatestEntry:           recordStore.GetMessageEntryByID(collRecord.LatestEntryID),
		LatestCheckpointEntry: recordStore.GetMessageEntryByID(collRecord.LatestCheckpointEntryID),
	}, nil

}

func FindRecordForCommit(ctx context.Context, recordStore storage.RecordStore, schemaUri string, parentRecordId string) (*FindRecordResult, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "FindRecordForCommit")
	defer sp.End()

	result := FindRecordResult{}

	record, messageEntry := recordStore.GetRecordForCommit(parentRecordId)
	if record == nil || messageEntry == nil {
		sp.AddEvent("Unable to find records for Commit")
		result.Status = "404"
		result.Error = errors.New("Not Found")
		return &result, nil
	}

	result.Status = "OK"
	result.LatestEntry = messageEntry
	result.Record = record

	return &result, nil

}

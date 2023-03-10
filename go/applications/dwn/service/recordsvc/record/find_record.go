package record

import (
	"context"
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

func FindRecordBySchemaAndRecordID(ctx context.Context, collectionStore storage.RecordStore, schemaUri string, recordId string) (*FindRecordResult, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "FindRecordBySchemaAndRecordID")
	defer sp.End()

	collRecord := collectionStore.GetRecord(recordId)
	if collRecord == nil {
		return &FindRecordResult{Status: "NOT_FOUND"}, nil
	}

	return &FindRecordResult{
		Status:                "OK",
		Error:                 nil,
		Record:                collRecord,
		InitialEntry:          collectionStore.GetMessageEntryByID(collRecord.InitialEntryID),
		LatestEntry:           collectionStore.GetMessageEntryByID(collRecord.LatestEntryID),
		LatestCheckpointEntry: collectionStore.GetMessageEntryByID(collRecord.LatestCheckpointEntryID),
	}, nil

}

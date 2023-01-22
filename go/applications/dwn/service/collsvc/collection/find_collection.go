package collection

import (
	"context"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
)

type FindCollectionResult struct {
	Status                string
	Error                 error
	Record                *storage.CollectionRecord
	InitialEntry          *storage.MessageEntry
	LatestEntry           *storage.MessageEntry
	LatestCheckpointEntry *storage.MessageEntry
}

func FindCollectionBySchemaAndRecordID(ctx context.Context, collectionStore storage.CollectionStore, schemaUri string, recordId string) (*FindCollectionResult, error) {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "FindCollectionBySchemaAndRecordID")
	defer sp.End()

	collRecord := collectionStore.GetCollectionRecord(recordId)
	if collRecord == nil {
		return &FindCollectionResult{Status: "NOT_FOUND"}, nil
	}

	return &FindCollectionResult{
		Status:                "OK",
		Error:                 nil,
		Record:                collRecord,
		InitialEntry:          collectionStore.GetMessageEntryByID(collRecord.InitialEntryID),
		LatestEntry:           collectionStore.GetMessageEntryByID(collRecord.LatestEntryID),
		LatestCheckpointEntry: collectionStore.GetMessageEntryByID(collRecord.LatestCheckpointEntryID),
	}, nil

}

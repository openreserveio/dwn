package collection

import "github.com/openreserveio/dwn/go/storage"

type FindCollectionResult struct {
	Status                string
	Error                 error
	Record                *storage.CollectionRecord
	InitialEntry          *storage.MessageEntry
	LatestEntry           *storage.MessageEntry
	LatestCheckpointEntry *storage.MessageEntry
}

func FindCollectionBySchemaAndRecordID(collectionStore storage.CollectionStore, schemaUri string, recordId string) (*FindCollectionResult, error) {

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

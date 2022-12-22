package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type CollectionStore interface {

	// Refactored methods here
	CreateCollectionRecord(record *CollectionRecord, initialEntry *MessageEntry) error
	AddCollectionMessageEntry(entry *MessageEntry) error
	GetMessageEntryByID(messageEntryID string) *MessageEntry
	GetCollectionRecord(recordId string) *CollectionRecord
}

type CollectionRecord struct {
	ID                      primitive.ObjectID `bson:"_id"`
	RecordID                string             `bson:"record_id"`
	CreatorDID              string             `bson:"creator_did"`
	OwnerDID                string             `bson:"owner_did"`
	WriterDIDs              []string           `bson:"writer_dids"`
	ReaderDIDs              []string           `bson:"reader_dids"`
	InitialEntryID          string             `bson:"initial_entry_id"`
	LatestEntryID           string             `bson:"latest_entry_id"`
	LatestCheckpointEntryID string             `bson:"latest_checkpoint_entry_id"`
}

type MessageEntry struct {
	ID             primitive.ObjectID `bson:"_id"`
	MessageEntryID string             `bson:"message_entry_id"`
	ParentEntryID  string             `bson:"parent_message_entry_id"`
	RecordID       string             `bson:"record_id"`
	Schema         string             `bson:"schema"`
	Method         string             `bson:"method"`
	Data           []byte             `bson:"data"`
	DataCID        string             `bson:"dataCID"`
	Protocol       string             `bson:"protocol"`
	CreatedDate    time.Time          `bson:"created_date"`
	PublishedDate  time.Time          `bson:"published_date"`
}

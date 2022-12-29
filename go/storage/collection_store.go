package storage

import (
	"github.com/openreserveio/dwn/go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type CollectionStore interface {

	// Refactored methods here
	CreateCollectionRecord(record *CollectionRecord, initialEntry *MessageEntry) error
	SaveCollectionRecord(record *CollectionRecord) error

	AddCollectionMessageEntry(entry *MessageEntry) error
	GetMessageEntryByID(messageEntryID string) *MessageEntry
	GetCollectionRecord(recordId string) *CollectionRecord
	DeleteCollectionMessageEntry(entry *MessageEntry) error
	DeleteCollectionMessageEntryByID(messageEntryId string) error
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
	model.Message
	ID             primitive.ObjectID `bson:"_id"`
	MessageEntryID string             `bson:"message_entry_id"`
}

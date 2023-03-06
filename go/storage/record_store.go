package storage

import (
	"github.com/openreserveio/dwn/go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type RecordStore interface {

	// Refactored methods here
	CreateRecord(record *Record, initialEntry *MessageEntry) error
	SaveRecord(record *Record) error

	AddMessageEntry(entry *MessageEntry) error
	GetMessageEntryByID(messageEntryID string) *MessageEntry
	GetRecord(recordId string) *Record
	DeleteMessageEntry(entry *MessageEntry) error
	DeleteMessageEntryByID(messageEntryId string) error
}

type Record struct {
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

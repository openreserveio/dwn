package storage

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

// Base interface for stores for collections
// Will probably need:
//   - methods for finding by content in the collection
type RecordStore interface {

	// Refactored methods here
	CreateRecord(ctx context.Context, record *Record, initialEntry *MessageEntry) error
	SaveRecord(ctx context.Context, record *Record) error

	AddMessageEntry(ctx context.Context, entry *MessageEntry) error
	GetMessageEntryByID(ctx context.Context, messageEntryID string) *MessageEntry
	GetRecord(ctx context.Context, recordId string) *Record
	GetRecordForCommit(ctx context.Context, recordId string) (*Record, *MessageEntry)
	DeleteMessageEntry(ctx context.Context, entry *MessageEntry) error
	DeleteMessageEntryByID(ctx context.Context, messageEntryId string) error

	BeginTx(ctx context.Context) error
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type Record struct {
	bun.BaseModel           `bun:"table:record"`
	ID                      string    `bun:"id,pk" json:"id"`
	RecordID                string    `bun:"record_id" json:"record_id"`
	CreatorDID              string    `bun:"creator_did" json:"creator_did"`
	OwnerDID                string    `bun:"owner_did" json:"owner_did"`
	WriterDIDs              []string  `bun:"writer_dids" json:"writer_dids"`
	ReaderDIDs              []string  `bun:"reader_dids" json:"reader_dids"`
	InitialEntryID          string    `bun:"initial_entry_id" json:"initial_entry_id"`
	LatestEntryID           string    `bun:"latest_entry_id" json:"latest_entry_id"`
	LatestCheckpointEntryID string    `bun:"latest_checkpoint_entry_id" json:"latest_checkpoint_entry_id"`
	CreateDate              time.Time `bun:"create_date" json:"create_date"`
}

type MessageEntry struct {
	bun.BaseModel          `bun:"table:message_entry"`
	ID                     string    `bun:"id,pk" json:"id"`
	MessageEntryID         string    `bun:"message_entry_id" json:"message_entry_id"`
	PreviousMessageEntryID string    `bun:"previous_message_entry_id" json:"previous_message_entry_id"`
	RecordID               string    `bun:"record_id" json:"record_id"`
	Message                []byte    `bun:"message" json:"message"`
	CreateDate             time.Time `bun:"create_date" json:"create_date"`
}

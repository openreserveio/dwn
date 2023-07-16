package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/erroring"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type RecordStorePostgres struct {
	DB       *bun.DB
	ActiveTx *bun.Tx
}

func NewRecordStorePostgres(pgConnString string) (*RecordStorePostgres, error) {

	// dwndb
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(pgConnString)))
	if sqldb == nil {
		return nil, &erroring.PostgresError{Msg: "failed to open DB connection"}
	}
	bundb := bun.NewDB(sqldb, pgdialect.New())

	return &RecordStorePostgres{
		DB: bundb,
	}, nil

}

func (recordStore *RecordStorePostgres) CreateRecord(ctx context.Context, record *storage.Record, initialEntry *storage.MessageEntry) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.CreateRecord")
	defer sp.End()

	// This must be a new record
	sp.AddEvent("Checking for existing record by ID or record_id")
	if record.ID != "" {
		err := &erroring.RecordError{Msg: "record already exists"}
		sp.RecordError(err)
		return err
	}

	// There cannot be one with the same RecordID for a create
	var recordCheck storage.Record
	err := recordStore.DB.NewSelect().Model(&recordCheck).Where("record_id = ?", record.RecordID).Scan(ctx, &recordCheck)
	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
		} else {
			sp.RecordError(err)
			return err
		}
	}
	if recordCheck.ID != "" {
		err := &erroring.RecordError{Msg: "record already exists"}
		sp.RecordError(err)
		return err
	}

	// Create a new initial entry, init the record and important fields
	sp.AddEvent("Creating Record and Initial Entry")
	initialEntry.ID = uuid.New().String()
	initialEntry.CreateDate = time.Now().UTC()
	record.ID = uuid.New().String()
	record.InitialEntryID = initialEntry.ID
	record.LatestEntryID = initialEntry.ID
	record.LatestCheckpointEntryID = initialEntry.ID
	record.CreateDate = time.Now().UTC()
	initialEntry.RecordID = record.RecordID

	err = nil
	if recordStore.ActiveTx != nil {
		sp.AddEvent("Creating Record and Initial Entry in transaction")
		_, err = recordStore.ActiveTx.NewInsert().Model(record).Exec(ctx)
		_, err = recordStore.ActiveTx.NewInsert().Model(initialEntry).Exec(ctx)
	} else {
		sp.AddEvent("Creating Record and Initial Entry without a transaction")
		_, err = recordStore.DB.NewInsert().Model(record).Exec(ctx)
		_, err = recordStore.DB.NewInsert().Model(initialEntry).Exec(ctx)
	}

	if err != nil {
		sp.RecordError(err)
		return err
	}

	return nil

}

func (recordStore *RecordStorePostgres) SaveRecord(ctx context.Context, record *storage.Record) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.SaveRecord")
	defer sp.End()

	var err error
	if recordStore.ActiveTx != nil {
		_, err = recordStore.ActiveTx.NewUpdate().Model(record).Exec(ctx)
	} else {
		_, err = recordStore.DB.NewUpdate().Model(record).Exec(ctx)
	}

	if err != nil {
		sp.RecordError(err)
		return err
	}

	return nil

}

func (recordStore *RecordStorePostgres) AddMessageEntry(ctx context.Context, entry *storage.MessageEntry) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.AddMessageEntry")
	defer sp.End()

	// Get the previous message entry
	sp.AddEvent("Getting previous message entry")
	previousEntry := recordStore.GetMessageEntryByID(ctx, entry.PreviousMessageEntryID)
	if previousEntry == nil {
		err := &erroring.RecordError{Msg: "previous message entry not found"}
		sp.RecordError(err)
		return err
	}

	// Get the associated Record
	sp.AddEvent("Getting associated record")
	rec := recordStore.GetRecord(ctx, entry.RecordID)
	if rec == nil {
		err := &erroring.RecordError{Msg: "Associated Record not found"}
		sp.RecordError(err)
		return err
	}

	var err error
	entry.ID = uuid.New().String()

	// This is generated upstream
	// entry.MessageEntryID = entry.ID
	entry.CreateDate = time.Now().UTC()
	if recordStore.ActiveTx != nil {
		_, err = recordStore.ActiveTx.NewInsert().Model(entry).Exec(ctx)
	} else {
		_, err = recordStore.DB.NewInsert().Model(entry).Exec(ctx)
	}

	if err != nil {
		sp.RecordError(err)
		return err
	}

	return nil

}

func (recordStore *RecordStorePostgres) GetMessageEntryByID(ctx context.Context, messageEntryID string) *storage.MessageEntry {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.GetMessageEntryByID")
	defer sp.End()

	var err error
	var entry storage.MessageEntry
	if recordStore.ActiveTx != nil {
		err = recordStore.ActiveTx.NewSelect().Model(&entry).Where("id = ?", messageEntryID).Scan(ctx, &entry)
	} else {
		err = recordStore.DB.NewSelect().Model(&entry).Where("id = ?", messageEntryID).Scan(ctx, &entry)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil
		} else {
			sp.RecordError(err)
			return nil
		}
	}

	return &entry

}

func (recordStore *RecordStorePostgres) GetRecord(ctx context.Context, recordId string) *storage.Record {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.GetRecord")
	defer sp.End()

	var err error
	var record storage.Record
	if recordStore.ActiveTx != nil {
		err = recordStore.ActiveTx.NewSelect().Model(&record).Where("record_id = ?", recordId).Scan(ctx, &record)
	} else {
		err = recordStore.DB.NewSelect().Model(&record).Where("record_id = ?", recordId).Scan(ctx, &record)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil
		} else {
			sp.RecordError(err)
			return nil
		}
	}

	return &record

}

func (recordStore *RecordStorePostgres) GetRecordForCommit(ctx context.Context, parentRecordId string) (*storage.Record, *storage.MessageEntry) {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.GetRecordForCommit")
	defer sp.End()

	// For this, the parent record ID is the record ID of the MESSAGE ENTRY, and then its associated Record
	// Generally this is for COMMIT a Write that is not the initial entry
	sp.AddEvent(fmt.Sprintf("Getting message entry by its record ID:  %s", parentRecordId))
	var err error
	var messageEntry storage.MessageEntry
	var record storage.Record
	if recordStore.ActiveTx != nil {
		err = recordStore.ActiveTx.NewSelect().Model(&messageEntry).Where("record_id = ?", parentRecordId).Scan(ctx, &messageEntry)
		err = recordStore.ActiveTx.NewSelect().Model(&record).Where("record_id = ?", parentRecordId).Scan(ctx, &record)
	} else {
		err = recordStore.DB.NewSelect().Model(&messageEntry).Where("record_id = ?", parentRecordId).Scan(ctx, &messageEntry)
		err = recordStore.DB.NewSelect().Model(&record).Where("record_id = ?", parentRecordId).Scan(ctx, &record)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil, nil
		} else {
			sp.RecordError(err)
			return nil, nil
		}
	}

	return &record, &messageEntry

}

func (recordStore *RecordStorePostgres) DeleteMessageEntry(ctx context.Context, entry *storage.MessageEntry) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "RecordStorePostgres.DeleteMessageEntry")
	defer sp.End()

	var err error
	if recordStore.ActiveTx != nil {
		_, err = recordStore.ActiveTx.NewDelete().Model(entry).Exec(ctx)
	} else {
		_, err = recordStore.DB.NewDelete().Model(entry).Exec(ctx)
	}

	if err != nil {
		sp.RecordError(err)
		return err
	}

	return nil

}

func (recordStore *RecordStorePostgres) DeleteMessageEntryByID(ctx context.Context, messageEntryId string) error {

	messageEntry := recordStore.GetMessageEntryByID(ctx, messageEntryId)
	if messageEntry == nil {
		return nil
	}

	return recordStore.DeleteMessageEntry(ctx, messageEntry)

}

func (recordStore *RecordStorePostgres) BeginTx(ctx context.Context) error {

	tx, err := recordStore.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	recordStore.ActiveTx = &tx
	return nil

}

func (recordStore *RecordStorePostgres) CommitTx(ctx context.Context) error {

	err := recordStore.ActiveTx.Commit()
	if err != nil {
		return err
	}

	recordStore.ActiveTx = nil
	return nil
}

func (recordStore *RecordStorePostgres) RollbackTx(ctx context.Context) error {
	err := recordStore.ActiveTx.Rollback()
	if err != nil {
		return err
	}

	recordStore.ActiveTx = nil
	return nil
}

package docdbstore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	RECORD_DATABASE = "dwn_records_db"

	COLLECTION_RECORD        = "records"
	COLLECTION_MESSAGE_ENTRY = "message_entries"

	COLLECTION_RECORD_ID_FIELD_NAME        = "record_id"
	COLLECTION_MESSAGE_ENTRY_ID_FIELD_NAME = "message_entry_id"
	COLLECTION_MESSAGE_ENTRY_RECORD_ID     = "message.record_id"
)

type RecordDocumentDBStore struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func CreateRecordDocumentDBStore(connectionUri string) (*RecordDocumentDBStore, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	dbStore := RecordDocumentDBStore{
		Client: client,
		DB:     client.Database(RECORD_DATABASE),
	}

	return &dbStore, nil

}

func (store *RecordDocumentDBStore) CreateRecord(record *storage.Record, initialEntry *storage.MessageEntry) error {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "CreateRecord")
	defer sp.End()

	initialEntry.RecordID = record.RecordID
	initialEntry.MessageEntryID = uuid.NewString()
	_, err := store.DB.Collection(COLLECTION_MESSAGE_ENTRY).InsertOne(context.Background(), initialEntry)
	if err != nil {
		return err
	}

	record.InitialEntryID = initialEntry.MessageEntryID
	record.LatestEntryID = initialEntry.MessageEntryID
	record.LatestCheckpointEntryID = initialEntry.MessageEntryID

	_, err = store.DB.Collection(COLLECTION_RECORD).InsertOne(context.Background(), record)
	if err != nil {
		return err
	}

	return nil

}

func (store *RecordDocumentDBStore) SaveRecord(record *storage.Record) error {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "SaveRecord")
	defer sp.End()

	_, err := store.DB.Collection(COLLECTION_RECORD).ReplaceOne(context.Background(), bson.D{{COLLECTION_RECORD_ID_FIELD_NAME, record.RecordID}}, record)
	if err != nil {
		return err
	}

	return nil
}

func (store *RecordDocumentDBStore) AddMessageEntry(entry *storage.MessageEntry) error {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "recordsvc.storage.docdbstore.AddMessageEntry")
	defer sp.End()

	sp.AddEvent("Inserting message entry for record")
	_, err := store.DB.Collection(COLLECTION_MESSAGE_ENTRY).InsertOne(context.Background(), entry)
	if err != nil {
		sp.RecordError(err)
		return err
	}

	return nil

}

func (store *RecordDocumentDBStore) GetMessageEntryByID(messageEntryID string) *storage.MessageEntry {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "GetMessageEntryByID")
	defer sp.End()

	res := store.DB.Collection(COLLECTION_MESSAGE_ENTRY).FindOne(context.Background(), bson.D{{COLLECTION_MESSAGE_ENTRY_ID_FIELD_NAME, messageEntryID}})
	if res.Err() != nil {
		log.Error("Error getting message entry by ID:  %v", res.Err())
		return nil
	}

	var messageEntry storage.MessageEntry
	err := res.Decode(&messageEntry)
	if err != nil {
		log.Error("Error decoding result of getting message entry by ID:  %v", err)
		return nil
	}

	return &messageEntry
}

func (store *RecordDocumentDBStore) GetRecord(recordId string) *storage.Record {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "GetRecord")
	defer sp.End()

	res := store.DB.Collection(COLLECTION_RECORD).FindOne(context.Background(), bson.D{{COLLECTION_RECORD_ID_FIELD_NAME, recordId}})
	if res.Err() != nil {
		if res.Err() != mongo.ErrNoDocuments {
			log.Error("Error getting record by ID:  %v", res.Err())
			return nil
		}

		log.Debug("No records found")
		return nil

	}

	var collectionRecord storage.Record
	err := res.Decode(&collectionRecord)
	if err != nil {
		log.Error("Error decoding result of getting collectionRecord by ID:  %v", err)
		return nil
	}

	return &collectionRecord
}

func (store *RecordDocumentDBStore) GetRecordForCommit(parentRecordId string) (*storage.Record, *storage.MessageEntry) {

	// tracing
	_, sp := observability.Tracer.Start(context.Background(), "recordsvc.storage.docdbstore.GetRecordForCommit")
	defer sp.End()

	// For this, the parent record ID is the record ID of the MESSAGE ENTRY, and then its associated Record
	// Generally this is for COMMIT a Write that is not the initial entry
	sp.AddEvent(fmt.Sprintf("Getting message entry by its record ID:  %s", parentRecordId))
	res := store.DB.Collection(COLLECTION_MESSAGE_ENTRY).FindOne(context.Background(), bson.D{{COLLECTION_MESSAGE_ENTRY_RECORD_ID, parentRecordId}})
	if res.Err() != nil {
		if res.Err() != mongo.ErrNoDocuments {
			sp.RecordError(res.Err())
			log.Error("Error getting record by ID:  %v", res.Err())
			return nil, nil
		}

		sp.AddEvent(fmt.Sprintf("No records found with record id:  %s", parentRecordId))
		log.Debug("No records found")
		return nil, nil

	}

	sp.AddEvent(fmt.Sprintf("Found message entry by its record ID:  %s", parentRecordId))
	var recordMessageEntry storage.MessageEntry
	err := res.Decode(&recordMessageEntry)
	if err != nil {
		log.Error("Error decoding result of getting record message entry by ID:  %v", err)
		sp.RecordError(err)
		return nil, nil
	}

	rec := store.GetRecord(recordMessageEntry.Descriptor.ParentID)
	if rec == nil {
		log.Error("We somehow found a message entry but no associated record for it?")
		return nil, nil
	}

	return rec, &recordMessageEntry

}

func (store *RecordDocumentDBStore) DeleteMessageEntry(entry *storage.MessageEntry) error {
	return store.DeleteMessageEntryByID(entry.MessageEntryID)
}

func (store *RecordDocumentDBStore) DeleteMessageEntryByID(messageEntryId string) error {

	_, err := store.DB.Collection(COLLECTION_MESSAGE_ENTRY).DeleteOne(context.Background(), bson.D{{COLLECTION_MESSAGE_ENTRY_ID_FIELD_NAME, messageEntryId}})
	if err != nil {
		return err
	}
	return nil
}

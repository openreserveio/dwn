package docdbstore

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionDocumentDBStore struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func CreateCollectionDocumentDBStore(connectionUri string) (*CollectionDocumentDBStore, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	dbStore := CollectionDocumentDBStore{
		Client: client,
		DB:     client.Database("dwn_collections_db"),
	}

	return &dbStore, nil

}

func (store *CollectionDocumentDBStore) CreateCollectionRecord(record *storage.CollectionRecord, initialEntry *storage.MessageEntry) error {

	initialEntry.RecordID = record.RecordID
	initialEntry.MessageEntryID = uuid.NewString()
	_, err := store.DB.Collection("dwn_collections").InsertOne(context.Background(), initialEntry)
	if err != nil {
		return err
	}

	record.InitialEntryID = initialEntry.MessageEntryID
	record.LatestEntryID = initialEntry.MessageEntryID
	record.LatestCheckpointEntryID = initialEntry.MessageEntryID

	_, err = store.DB.Collection("dwn_collections").InsertOne(context.Background(), record)
	if err != nil {
		return err
	}

	return nil

}

func (store *CollectionDocumentDBStore) AddCollectionMessageEntry(entry *storage.MessageEntry) error {

	// Get Record
	if entry.RecordID == "" {
		return errors.New("No Record ID")
	}
	collectionRecord := store.GetCollectionRecord(entry.RecordID)
	if collectionRecord == nil {
		return errors.New("No Record Found")
	}

	_, err := store.DB.Collection("dwn_collections").InsertOne(context.Background(), entry)
	if err != nil {
		return err
	}

	return nil

}

func (store *CollectionDocumentDBStore) GetMessageEntryByID(messageEntryID string) *storage.MessageEntry {

	res := store.DB.Collection("dwn_collections").FindOne(context.Background(), bson.D{{"message_entry_id", messageEntryID}})
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

func (store *CollectionDocumentDBStore) GetCollectionRecord(recordId string) *storage.CollectionRecord {

	res := store.DB.Collection("dwn_collections").FindOne(context.Background(), bson.D{{"record_id", recordId}})
	if res.Err() != nil {
		log.Error("Error getting record by ID:  %v", res.Err())
		return nil
	}

	var collectionRecord storage.CollectionRecord
	err := res.Decode(&collectionRecord)
	if err != nil {
		log.Error("Error decoding result of getting collectionRecord by ID:  %v", err)
		return nil
	}

	return &collectionRecord
}

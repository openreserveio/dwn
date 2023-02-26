package docdbstore

import (
	"context"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	HOOK_RECORD_COLLECTION       = "hook_records"
	HOOK_CONFIG_ENTRY_COLLECTION = "hook_config_entry"

	HOOK_RECORD_ID_FIELD_NAME = "hook_record_id"
)

type HookDocumentDBStore struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func CreateHookDocumentDBStore(connectionUri string) (*HookDocumentDBStore, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	dbStore := HookDocumentDBStore{
		Client: client,
		DB:     client.Database("dwn_hook_db"),
	}

	return &dbStore, nil

}

func (store *HookDocumentDBStore) CreateHookRecord(ctx context.Context, hookRecord *storage.HookRecord, initialConfiguration *storage.HookConfigurationEntry) error {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "hook_store.CreateHookRecord")
	defer sp.End()

	initialConfiguration.ConfigurationEntryID = uuid.NewString()
	initialConfiguration.HookRecordID = initialConfiguration.RecordID
	_, err := store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).InsertOne(ctx, initialConfiguration)
	if err != nil {
		return err
	}

	hookRecord.InitialHookConfigurationEntryID = initialConfiguration.ConfigurationEntryID
	hookRecord.LatestHookConfigurationEntryID = initialConfiguration.ConfigurationEntryID
	hookRecord.CreatorDID = initialConfiguration.Processing.AuthorDID
	_, err = store.DB.Collection(HOOK_RECORD_COLLECTION).InsertOne(ctx, hookRecord)
	if err != nil {
		return err
	}

	return nil

}

func (store *HookDocumentDBStore) UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *storage.HookConfigurationEntry) error {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "hook_store.UpdateHookRecord")
	defer sp.End()

	// Get current setup, error if not found
	res := store.DB.Collection(HOOK_RECORD_COLLECTION).FindOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if res.Err() != nil {
		log.Error("Error getting message entry by ID:  %v", res.Err())
		return nil
	}

}

func (store *HookDocumentDBStore) DeleteHookRecord(ctx context.Context, hookRecordId string) error {
	// tracing
	_, sp := observability.Tracer.Start(ctx, "hook_store.DeleteHookRecord")
	defer sp.End()

	//TODO implement me
	panic("implement me")
}

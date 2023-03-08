package docdbstore

import (
	"context"
	"errors"
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

	HOOK_RECORD_ID_FIELD_NAME       = "hook_record_id"
	HOOK_CONFIG_ENTRY_ID_FIELD_NAME = "configuration_entry_id"
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

func (store *HookDocumentDBStore) GetHookRecord(ctx context.Context, hookRecordId string) (*storage.HookRecord, *storage.HookConfigurationEntry, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.GetHookRecord")
	defer sp.End()

	// Get the Hook Record
	res := store.DB.Collection(HOOK_RECORD_COLLECTION).FindOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if res.Err() != nil {

		if res.Err() != mongo.ErrNoDocuments {
			log.Error("Error getting record by ID:  %v", res.Err())
			return nil, nil, res.Err()
		}

		log.Debug("No records found")
		return nil, nil, nil

	}

	var hookRecord storage.HookRecord
	var latestConfigEntry storage.HookConfigurationEntry

	err := res.Decode(&hookRecord)
	if err != nil {
		return nil, nil, err
	}

	// Get the latest configuration entry
	resConfig := store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).FindOne(ctx, bson.D{{HOOK_CONFIG_ENTRY_ID_FIELD_NAME, hookRecord.LatestHookConfigurationEntryID}})
	if resConfig.Err() != nil {

		if res.Err() != mongo.ErrNoDocuments {
			log.Error("Error getting record by ID:  %v", res.Err())
			return nil, nil, res.Err()
		}

		log.Debug("No records found - should always be at least 1 config entry!")
		return nil, nil, nil

	}

	err = resConfig.Decode(&latestConfigEntry)
	if err != nil {
		return nil, nil, err
	}

	return &hookRecord, &latestConfigEntry, nil

}

func (store *HookDocumentDBStore) GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*storage.HookRecord, []*storage.HookConfigurationEntry, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.GetHookRecordConfigurationEntries")
	defer sp.End()

	// Get the Hook Record
	res := store.DB.Collection(HOOK_RECORD_COLLECTION).FindOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if res.Err() != nil {

		if res.Err() != mongo.ErrNoDocuments {
			log.Error("Error getting record by ID:  %v", res.Err())
			return nil, nil, res.Err()
		}

		log.Debug("No records found")
		return nil, nil, nil

	}

	var hookRecord storage.HookRecord
	var configEntries []*storage.HookConfigurationEntry

	err := res.Decode(&hookRecord)
	if err != nil {
		return nil, nil, err
	}

	// Get all the configuration entries for this hook record
	cur, err := store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).Find(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	defer cur.Close(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Must be at least one
	if cur.RemainingBatchLength() == 0 {
		return nil, nil, nil
	}

	for cur.Next(ctx) {

		var entry storage.HookConfigurationEntry
		err = cur.Decode(&entry)
		if err != nil {
			return nil, nil, err
		}

		configEntries = append(configEntries, &entry)

	}

	return &hookRecord, configEntries, nil

}

func (store *HookDocumentDBStore) CreateHookRecord(ctx context.Context, hookRecord *storage.HookRecord, initialConfiguration *storage.HookConfigurationEntry) error {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.CreateHookRecord")
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
	hookRecord.FilterDataRecordID = initialConfiguration.Descriptor.Filter.RecordID
	_, err = store.DB.Collection(HOOK_RECORD_COLLECTION).InsertOne(ctx, hookRecord)
	if err != nil {
		return err
	}

	return nil

}

func (store *HookDocumentDBStore) UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *storage.HookConfigurationEntry) error {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.UpdateHookRecord")
	defer sp.End()

	// Get current setup, error if not found
	res := store.DB.Collection(HOOK_RECORD_COLLECTION).FindOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if res.Err() != nil {
		log.Error("Error getting message entry by ID:  %v", res.Err())
		return nil
	}

	var hookRecord storage.HookRecord
	err := res.Decode(&hookRecord)
	if err != nil {
		return err
	}

	// Give new ID, make sure it points to the hook record
	updatedConfiguration.ConfigurationEntryID = uuid.NewString()
	updatedConfiguration.RecordID = hookRecord.HookRecordID

	// Latest hook config should point to this new one
	hookRecord.LatestHookConfigurationEntryID = updatedConfiguration.ConfigurationEntryID

	// First, insert the new configuration
	_, err = store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).InsertOne(ctx, updatedConfiguration)
	if err != nil {
		return err
	}

	// Then, update the hook record
	updateRes, err := store.DB.Collection(HOOK_RECORD_COLLECTION).ReplaceOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecord.HookRecordID}}, &hookRecord)
	if err != nil {
		return err
	}
	if updateRes.ModifiedCount == 0 {
		return errors.New("Unable to modify the hook record")
	}

	return nil
}

func (store *HookDocumentDBStore) DeleteHookRecord(ctx context.Context, hookRecordId string) error {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.DeleteHookRecord")
	defer sp.End()

	// delete hook records
	_, err := store.DB.Collection(HOOK_RECORD_COLLECTION).DeleteOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if err != nil {
		log.Error("Error deleting record by Hook Record ID:  %v", err)
		return err
	}

	// Delete Entries BYE
	_, err = store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).DeleteOne(ctx, bson.D{{HOOK_RECORD_ID_FIELD_NAME, hookRecordId}})
	if err != nil {
		log.Error("Error deleting record config entries by Hook Record ID:  %v", err)
		return err
	}

	return nil

}

func (store *HookDocumentDBStore) FindHookRecordsForDataRecord(ctx context.Context, dataRecordId string) (map[*storage.HookRecord]*storage.HookConfigurationEntry, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "hook_store.FindHookRecordsForDataRecord")
	defer sp.End()

	// set up filter for mongodb
	filter := bson.D{{"filter_data_record_id", dataRecordId}}
	findCur, err := store.DB.Collection(HOOK_RECORD_COLLECTION).Find(ctx, filter)
	if err != nil {
		log.Error("Error while looking for hook records for data record:  %v", err)
		sp.RecordError(err)
		return nil, err
	}
	if findCur.Err() != nil {

		if findCur.Err() != mongo.ErrNoDocuments {
			log.Error("Error finding hook records by data record ID:  %v", findCur.Err())
			return nil, findCur.Err()
		}

		log.Debug("No records found")
		return nil, nil

	}

	res := make(map[*storage.HookRecord]*storage.HookConfigurationEntry)
	for findCur.Next(ctx) {
		var hookRecord storage.HookRecord
		findCur.Decode(&hookRecord)

		// Get the latest entry
		latestEntry := store.DB.Collection(HOOK_CONFIG_ENTRY_COLLECTION).FindOne(ctx, bson.D{{HOOK_CONFIG_ENTRY_ID_FIELD_NAME, hookRecord.LatestHookConfigurationEntryID}})
		if findCur.Err() != nil {

			if findCur.Err() != mongo.ErrNoDocuments {
				log.Error("Error finding hook configuration entries by latest entry ID:  %v", findCur.Err())
				return nil, findCur.Err()
			}

			log.Debug("No records found")
			return nil, nil

		}

		var latestConfigEntry storage.HookConfigurationEntry
		latestEntry.Decode(&latestConfigEntry)

		res[&hookRecord] = &latestConfigEntry
	}

	return res, nil
}

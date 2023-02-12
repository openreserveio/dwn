package docdbstore

import (
	"context"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	//TODO implement me
	panic("implement me")
}

func (store *HookDocumentDBStore) UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *storage.HookConfigurationEntry) {
	//TODO implement me
	panic("implement me")
}

package convoywebhook

import (
	"context"
	"github.com/openreserveio/dwn/go/storage"
)

type ConvoyWebhookStore struct {
}

func NewConvoyWebhookStore() *ConvoyWebhookStore {
	return &ConvoyWebhookStore{}
}

func (cws *ConvoyWebhookStore) GetHookRecord(ctx context.Context, hookRecordId string) (*storage.HookRecord, *storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*storage.HookRecord, []*storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) CreateHookRecord(ctx context.Context, hookRecord *storage.HookRecord, initialConfiguration *storage.HookConfigurationEntry) error {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *storage.HookConfigurationEntry) error {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) DeleteHookRecord(ctx context.Context, hookRecordId string) error {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) FindHookRecordsForDataRecord(ctx context.Context, dataRecordId string) (map[*storage.HookRecord]*storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (cws *ConvoyWebhookStore) FindHookRecordsForSchemaAndProtocol(ctx context.Context, schemaUri string, protocol string, protocolVersion string) (map[*storage.HookRecord]*storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

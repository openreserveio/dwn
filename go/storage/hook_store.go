package storage

import (
	"context"
	"github.com/openreserveio/dwn/go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base interface for storing Hook configuration from the user
type HookStore interface {
	GetHookRecord(ctx context.Context, hookRecordId string) (*HookRecord, *HookConfigurationEntry, error)
	GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*HookRecord, []*HookConfigurationEntry, error)
	CreateHookRecord(ctx context.Context, hookRecord *HookRecord, initialConfiguration *HookConfigurationEntry) error
	UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *HookConfigurationEntry) error
	DeleteHookRecord(ctx context.Context, hookRecordId string) error

	FindHookRecordsForDataRecord(ctx context.Context, dataRecordId string) (map[*HookRecord]*HookConfigurationEntry, error)
}

type HookRecord struct {
	ID                              primitive.ObjectID `bson:"_id"`
	HookRecordID                    string             `bson:"hook_record_id"`
	CreatorDID                      string             `bson:"creator_did"`
	InitialHookConfigurationEntryID string             `bson:"initial_hook_config_entry_id"`
	LatestHookConfigurationEntryID  string             `bson:"latest_hook_config_entry_id"`

	// For Indexing
	FilterDataRecordID string `bson:"filter_data_record_id"`
}

type HookConfigurationEntry struct {
	model.Message
	ID                   primitive.ObjectID `bson:"_id"`
	ConfigurationEntryID string             `bson:"configuration_entry_id"`
	HookRecordID         string             `bson:"hook_record_id"`
}

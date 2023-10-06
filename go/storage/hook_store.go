package storage

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

// Base interface for storing Hook configuration from the user
type HookStore interface {
	GetHookRecord(ctx context.Context, hookRecordId string) (*HookRecord, *HookConfigurationEntry, error)
	GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*HookRecord, []*HookConfigurationEntry, error)
	CreateHookRecord(ctx context.Context, hookRecord *HookRecord, initialConfiguration *HookConfigurationEntry) error
	UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *HookConfigurationEntry) error
	DeleteHookRecord(ctx context.Context, hookRecordId string) error

	FindHookRecordsForDataRecord(ctx context.Context, dataRecordId string) (map[*HookRecord]*HookConfigurationEntry, error)
	FindHookRecordsForSchemaAndProtocol(ctx context.Context, schemaUri string, protocol string, protocolVersion string) (map[*HookRecord]*HookConfigurationEntry, error)

	BeginTx(ctx context.Context) error
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type HookRecord struct {
	bun.BaseModel                   `bun:"table:hook_record"`
	ID                              string    `bun:"id,pk" json:"id"`
	HookRecordID                    string    `bun:"hook_record_id" json:"hook_record_id"`
	CreatorDID                      string    `bun:"creator_did" json:"creator_did"`
	OwnerDID                        string    `bun:"owner_did" json:"owner_did"`
	InitialHookConfigurationEntryID string    `bun:"initial_hook_config_entry_id" json:"initial_hook_config_entry_id"`
	LatestHookConfigurationEntryID  string    `bun:"latest_hook_config_entry_id" json:"latest_hook_config_entry_id"`
	CreateDate                      time.Time `bun:"create_date" json:"create_date"`

	// For Indexing
	FilterDataRecordID    string `bun:"filter_data_record_id" json:"filter_data_record_id"`
	FilterSchema          string `bun:"filter_schema" json:"filter_schema"`
	FilterProtocol        string `bun:"filter_protocol" json:"filter_protocol"`
	FilterProtocolVersion string `bun:"filter_protocol_version" json:"filter_protocol_version"`
}

type HookConfigurationEntry struct {
	bun.BaseModel        `bun:"table:hook_configuration_entry"`
	ID                   string    `bun:"id,pk" json:"id"`
	ConfigurationEntryID string    `bun:"configuration_entry_id" json:"configuration_entry_id"`
	HookRecordID         string    `bun:"hook_record_id" json:"hook_record_id"`
	Message              []byte    `bun:"message" json:"message"`
	CreateDate           time.Time `bun:"create_date" json:"create_date"`
}

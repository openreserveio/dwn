package pgsql

import (
	"context"
	"database/sql"
	"github.com/openreserveio/dwn/go/erroring"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type HookStorePostgres struct {
	DB       *bun.DB
	ActiveTx *bun.Tx
}

func NewHookStorePostgres(pgConnString string) (*HookStorePostgres, error) {

	// dwndb
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(pgConnString)))
	if sqldb == nil {
		return nil, &erroring.PostgresError{Msg: "failed to open DB connection"}
	}
	bundb := bun.NewDB(sqldb, pgdialect.New())

	return &HookStorePostgres{
		DB: bundb,
	}, nil
}

func (hookStore *HookStorePostgres) GetHookRecord(ctx context.Context, hookRecordId string) (*storage.HookRecord, *storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (hookStore *HookStorePostgres) GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*storage.HookRecord, []*storage.HookConfigurationEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (hookStore *HookStorePostgres) CreateHookRecord(ctx context.Context, hookRecord *storage.HookRecord, initialConfiguration *storage.HookConfigurationEntry) error {
	//TODO implement me
	panic("implement me")
}

func (hookStore *HookStorePostgres) UpdateHookRecord(ctx context.Context, hookRecordId string, updatedConfiguration *storage.HookConfigurationEntry) error {
	//TODO implement me
	panic("implement me")
}

func (hookStore *HookStorePostgres) DeleteHookRecord(ctx context.Context, hookRecordId string) error {
	//TODO implement me
	panic("implement me")
}

func (hookStore *HookStorePostgres) FindHookRecordsForDataRecord(ctx context.Context, dataRecordId string) (map[*storage.HookRecord]*storage.HookConfigurationEntry, error) {
	//TODO implement me
	var hookMap map[*storage.HookRecord]*storage.HookConfigurationEntry
	return hookMap, nil
}

func (hookStore *HookStorePostgres) FindHookRecordsForSchemaAndProtocol(ctx context.Context, schemaUri string, protocol string, protocolVersion string) (map[*storage.HookRecord]*storage.HookConfigurationEntry, error) {
	//TODO implement me
	var hookMap map[*storage.HookRecord]*storage.HookConfigurationEntry
	return hookMap, nil
}

func (hookStore *HookStorePostgres) BeginTx(ctx context.Context) error {

	tx, err := hookStore.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	hookStore.ActiveTx = &tx
	return nil

}

func (hookStore *HookStorePostgres) CommitTx(ctx context.Context) error {
	return hookStore.ActiveTx.Commit()
}

func (hookStore *HookStorePostgres) RollbackTx(ctx context.Context) error {
	return hookStore.ActiveTx.Rollback()
}

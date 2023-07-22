package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/erroring"
	"github.com/openreserveio/dwn/go/observability"
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

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.GetHookRecord")
	defer sp.End()

	sp.AddEvent("Getting Record by Record ID")
	hookRecord, err := hookStore.getHookRecordById(ctx, hookRecordId)
	if err != nil {
		sp.RecordError(err)
		return nil, nil, err
	}

	sp.AddEvent("Getting Configuration Entry by Record ID")
	hookConfigurationEntry, err := hookStore.getHookConfigurationEntryById(ctx, hookRecord.LatestHookConfigurationEntryID)
	if err != nil {
		sp.RecordError(err)
		return nil, nil, err
	}

	if hookConfigurationEntry == nil {
		err = errors.New("no configuration entry found for record, which should be impossible")
		sp.RecordError(err)
		return nil, nil, err
	}

	return hookRecord, hookConfigurationEntry, nil

}

func (hookStore *HookStorePostgres) GetHookRecordConfigurationEntries(ctx context.Context, hookRecordId string) (*storage.HookRecord, []*storage.HookConfigurationEntry, error) {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.GetHookRecordConfigurationEntries")
	defer sp.End()

	sp.AddEvent("Getting Record by Record ID")
	hookRecord, err := hookStore.getHookRecordById(ctx, hookRecordId)
	if err != nil {
		sp.RecordError(err)
		return nil, nil, err
	}

	sp.AddEvent("Getting Hook Config Entries")
	entries, err := hookStore.getHookConfigurationEntries(ctx, hookRecordId)
	if err != nil {
		sp.RecordError(err)
		return nil, nil, err
	}

	return hookRecord, entries, nil

}

func (hookStore *HookStorePostgres) CreateHookRecord(ctx context.Context, hookRecord *storage.HookRecord, initialConfiguration *storage.HookConfigurationEntry) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.CreateHookRecord")
	defer sp.End()

	sp.AddEvent("Checking for existing")
	existingHookRecord, err := hookStore.getHookRecordById(ctx, hookRecord.HookRecordID)
	if err != nil {
		sp.RecordError(err)
		return err
	}
	if existingHookRecord != nil {
		err = errors.New("hook record already exists")
		sp.RecordError(err)
		return err
	}

	sp.AddEvent("Inserting Hook Record")
	hookRecord.ID = uuid.NewString()
	initialConfiguration.ID = uuid.NewString()
	initialConfiguration.ConfigurationEntryID = initialConfiguration.ID
	initialConfiguration.HookRecordID = hookRecord.HookRecordID
	hookRecord.InitialHookConfigurationEntryID = initialConfiguration.ConfigurationEntryID

	return nil

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

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.BeginTx")
	defer sp.End()

	sp.AddEvent("BeginTx")
	tx, err := hookStore.DB.BeginTx(ctx, nil)
	if err != nil {
		sp.RecordError(err)
		return err
	}

	sp.AddEvent("ActiveTx started")
	hookStore.ActiveTx = &tx
	return nil

}

func (hookStore *HookStorePostgres) CommitTx(ctx context.Context) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.CommitTx")
	defer sp.End()

	sp.AddEvent("CommitTx")
	err := hookStore.ActiveTx.Commit()
	if err != nil {
		sp.RecordError(err)
		return err
	}

	sp.AddEvent("ActiveTx committed")
	hookStore.ActiveTx = nil
	return nil
}

func (hookStore *HookStorePostgres) RollbackTx(ctx context.Context) error {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.RollbackTx")
	defer sp.End()

	sp.AddEvent("RollbackTx")
	err := hookStore.ActiveTx.Rollback()
	if err != nil {
		sp.RecordError(err)
		return err
	}

	sp.AddEvent("ActiveTx rolled back")
	hookStore.ActiveTx = nil
	return nil
}

func (hookStore *HookStorePostgres) getHookRecordById(ctx context.Context, hookRecordId string) (*storage.HookRecord, error) {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.getHookRecordById")
	defer sp.End()

	sp.AddEvent("Getting Record by Record ID")
	var hookRecord storage.HookRecord
	var err error
	if hookStore.ActiveTx != nil {
		err = hookStore.ActiveTx.NewSelect().Model(&hookRecord).Where("record_id = ?", hookRecordId).Scan(ctx, &hookRecord)
	} else {
		err = hookStore.DB.NewSelect().Model(&hookRecord).Where("record_id = ?", hookRecordId).Scan(ctx, &hookRecord)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil, nil
		} else {
			// this is not fine
			sp.RecordError(err)
			return nil, err
		}
	}

	return &hookRecord, nil

}

func (hookStore *HookStorePostgres) getHookConfigurationEntryById(ctx context.Context, hookConfigurationEntryId string) (*storage.HookConfigurationEntry, error) {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.getHookConfigurationEntryById")
	defer sp.End()

	var err error
	var hookConfigurationEntry storage.HookConfigurationEntry
	if hookStore.ActiveTx != nil {
		err = hookStore.ActiveTx.NewSelect().Model(&hookConfigurationEntry).Where("configuration_entry_id = ?", hookConfigurationEntryId).Scan(ctx, &hookConfigurationEntry)
	} else {
		err = hookStore.DB.NewSelect().Model(&hookConfigurationEntry).Where("configuration_entry_id = ?", hookConfigurationEntryId).Scan(ctx, &hookConfigurationEntry)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil, nil
		} else {
			// this is not fine
			sp.RecordError(err)
			return nil, err
		}
	}

	return &hookConfigurationEntry, nil

}

func (hookStore *HookStorePostgres) getHookConfigurationEntries(ctx context.Context, hookRecordId string) ([]*storage.HookConfigurationEntry, error) {

	// Observability
	ctx, sp := observability.Tracer().Start(ctx, "HookStorePostgres.getHookConfigurationEntries")
	defer sp.End()

	var err error
	var hookConfigurationEntries []*storage.HookConfigurationEntry
	if hookStore.ActiveTx != nil {
		err = hookStore.ActiveTx.NewSelect().Model(&hookConfigurationEntries).Where("record_id = ?", hookRecordId).Order("create_date DESC").Scan(ctx, &hookConfigurationEntries)
	} else {
		err = hookStore.DB.NewSelect().Model(&hookConfigurationEntries).Where("record_id = ?", hookRecordId).Order("create_date DESC").Scan(ctx, &hookConfigurationEntries)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// This is fine
			return nil, nil
		} else {
			// this is not fine
			sp.RecordError(err)
			return nil, err
		}
	}

	return hookConfigurationEntries, nil

}

package hooksvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/framework/events"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/pgsql"
	"strings"
)

type HookService struct {
	services.UnimplementedHookServiceServer
	HookStore storage.HookStore
	EventHub  *events.EventHub
}

func CreateHookService(hookStoreConnectionURI string, queueConnUrl string) (*HookService, error) {

	// Setup Hook Store
	hookStore, err := pgsql.NewHookStorePostgres(hookStoreConnectionURI)
	if err != nil {
		log.Fatal("Unable to connect to hook store:  %v", err)
		return nil, err
	}

	// Event Hub
	eventHub, err := events.CreateEventHub(queueConnUrl)
	if err != nil {
		log.Fatal("Unable to connect to hook store:  %v", err)
		return nil, err
	}

	hookService := HookService{
		HookStore: hookStore,
		EventHub:  eventHub,
	}

	return &hookService, nil
}

func (hookService HookService) RegisterHook(ctx context.Context, request *services.RegisterHookRequest) (*services.RegisterHookResponse, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "HookService.RegisterHook")
	defer sp.End()

	response := services.RegisterHookResponse{}

	var message model.Message
	err := json.Unmarshal(request.Message, &message)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	// Is there already a hook record?
	hookRecord, _, err := hookService.HookStore.GetHookRecord(ctx, message.RecordID)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if hookRecord != nil {
		response.Status = &services.CommonStatus{Status: services.Status_CONFLICT, Details: "There is already a hook record ID present"}
		return &response, nil
	}

	// Create Initial Config Entry
	messageBytes, err := json.Marshal(&message)
	configEntry := storage.HookConfigurationEntry{
		Message:              messageBytes,
		ConfigurationEntryID: uuid.NewString(),
		HookRecordID:         message.RecordID,
	}

	// Create Hook Record!
	createdHookRecord := storage.HookRecord{
		HookRecordID:                    message.RecordID,
		CreatorDID:                      message.Processing.AuthorDID,
		OwnerDID:                        message.Processing.RecipientDID,
		InitialHookConfigurationEntryID: configEntry.ConfigurationEntryID,
		LatestHookConfigurationEntryID:  configEntry.ConfigurationEntryID,
		FilterProtocol:                  message.Descriptor.Filter.Protocol,
		FilterProtocolVersion:           message.Descriptor.Filter.ProtocolVersion,
		FilterSchema:                    message.Descriptor.Filter.Schema,
		FilterDataRecordID:              message.Descriptor.Filter.RecordID,
	}

	// Store!
	err = hookService.HookStore.CreateHookRecord(ctx, &createdHookRecord, &configEntry)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	return &response, nil

}

func (hookService HookService) UpdateHook(ctx context.Context, request *services.UpdateHookRequest) (*services.UpdateHookResponse, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "HookService.UpdateHook")
	defer sp.End()

	response := services.UpdateHookResponse{}

	var message model.Message
	err := json.Unmarshal(request.Message, &message)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	// Get the existing hook record, error if not found
	hookRecord, _, err := hookService.HookStore.GetHookRecord(ctx, message.RecordID)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if hookRecord == nil {
		response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND, Details: "No hook record ID found to update"}
		return &response, nil
	}

	// Create Config Entry
	messageBytes, _ := json.Marshal(&message)

	configEntry := storage.HookConfigurationEntry{
		Message:              messageBytes,
		ConfigurationEntryID: uuid.NewString(),
		HookRecordID:         message.RecordID,
	}

	// Store!
	err = hookService.HookStore.UpdateHookRecord(ctx, message.RecordID, &configEntry)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	return &response, nil
}

func (hookService HookService) GetHookByRecordId(ctx context.Context, request *services.GetHookByRecordIdRequest) (*services.GetHookByRecordIdResponse, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "HookService.GetHookByRecordId")
	defer sp.End()

	response := services.GetHookByRecordIdResponse{}

	sp.AddEvent("Get hook from store")
	if request.RecordId == "" {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "Hook Record ID is required"}
		return &response, nil
	}

	hookRecord, latestEntry, err := hookService.HookStore.GetHookRecord(ctx, request.RecordId)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if hookRecord == nil || latestEntry == nil {
		response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND}
		return &response, nil
	}

	var latestEntryMessage model.Message
	err = json.Unmarshal(latestEntry.Message, &latestEntryMessage)

	hookDef := services.HookDefinition{
		HookId: hookRecord.HookRecordID,
		Uri:    latestEntryMessage.Descriptor.URI,
	}
	if strings.Contains(hookDef.Uri, "http") {
		hookDef.HookChannel = services.HookDefinition_HTTP_CALLBACK
	} else if strings.Contains(hookDef.Uri, "grpc") {
		hookDef.HookChannel = services.HookDefinition_GRPC_CALLBACK
	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	response.HookDefinition = &hookDef

	return &response, nil

}

func (hookService HookService) GetHooksForRecord(ctx context.Context, request *services.GetHooksForRecordRequest) (*services.GetHooksForRecordResponse, error) {
	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "HookService.GetHooksForRecord")
	defer sp.End()

	response := services.GetHooksForRecordResponse{}

	sp.AddEvent("Get hook from store")
	if request.RecordId == "" {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "Record ID is required"}
		return &response, nil
	}

	hookRecord, latestEntry, err := hookService.HookStore.GetHookRecord(ctx, request.RecordId)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if hookRecord == nil || latestEntry == nil {
		response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND, Details: "No hook record ID found"}
		return &response, nil
	}

	var latestEntryMessage model.Message
	err = json.Unmarshal(latestEntry.Message, &latestEntryMessage)

	hookDef := services.HookDefinition{
		HookId: hookRecord.HookRecordID,
		Uri:    latestEntryMessage.Descriptor.URI,
	}
	if strings.Contains(hookDef.Uri, "http") {
		hookDef.HookChannel = services.HookDefinition_HTTP_CALLBACK
	} else if strings.Contains(hookDef.Uri, "grpc") {
		hookDef.HookChannel = services.HookDefinition_GRPC_CALLBACK
	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	response.HookDefinitions = append(response.HookDefinitions, &hookDef)

	return &response, nil
}

func (hookService HookService) NotifyHooksOfRecordEvent(ctx context.Context, request *services.NotifyHooksOfRecordEventRequest) (*services.NotifyHooksOfRecordEventResponse, error) {

	// tracing
	ctx, sp := observability.Tracer().Start(ctx, "HookService.NotifyHooksOfRecordEvent")
	defer sp.End()

	response := services.NotifyHooksOfRecordEventResponse{}

	// Get Hooks for Record
	sp.AddEvent("Get hooks by data record to see if there are notifications to send")
	if request.RecordId == "" {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "Record ID is required for notifications"}
		return &response, nil
	}

	hookRecords, err := hookService.HookStore.FindHookRecordsForDataRecord(ctx, request.RecordId)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if hookRecords == nil || len(hookRecords) == 0 {
		sp.AddEvent("No Hook Records by record id")
	} else {

		for _, latestEntry := range hookRecords {

			var latestEntryMessage model.Message
			err = json.Unmarshal(latestEntry.Message, &latestEntryMessage)

			sp.AddEvent(fmt.Sprintf("--->  HOOK RECORD URL:  %s", latestEntryMessage.Descriptor.URI))
			hookService.EventHub.RaiseNotifyCallbackEvent(latestEntryMessage.Descriptor.Schema, request.RecordId, latestEntryMessage.Descriptor.Protocol, latestEntryMessage.Descriptor.ProtocolVersion, latestEntryMessage.Descriptor.URI)
		}

	}

	// Get Hooks for schema & protocol
	sp.AddEvent("Get hooks by schema and protocol to see if there are notifications to send")
	spHookRecords, err := hookService.HookStore.FindHookRecordsForSchemaAndProtocol(ctx, request.Schema, request.Protocol, request.ProtocolVersion)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if spHookRecords == nil || len(spHookRecords) == 0 {
		sp.AddEvent("No Hook Records by schema and protocol")
	} else {

		for _, latestEntry := range spHookRecords {
			var latestEntryMessage model.Message
			err = json.Unmarshal(latestEntry.Message, &latestEntryMessage)

			sp.AddEvent(fmt.Sprintf("--->  HOOK RECORD URL:  %s", latestEntryMessage.Descriptor.URI))
			hookService.EventHub.RaiseNotifyCallbackEvent(latestEntryMessage.Descriptor.Schema, request.RecordId, latestEntryMessage.Descriptor.Protocol, latestEntryMessage.Descriptor.ProtocolVersion, latestEntryMessage.Descriptor.URI)
		}

	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	return &response, nil

}

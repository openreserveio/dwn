package hooksvc

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/framework/events"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HookService struct {
	services.UnimplementedHookServiceServer
	HookStore storage.HookStore
	EventHub  *events.EventHub
}

func CreateHookService(hookStoreConnectionURI string, queueConnUrl string) (*HookService, error) {

	// Setup Hook Store
	hookStore, err := docdbstore.CreateHookDocumentDBStore(hookStoreConnectionURI)
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
	_, sp := observability.Tracer.Start(ctx, "RegisterHook")
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
	configEntry := storage.HookConfigurationEntry{
		Message:              message,
		ID:                   primitive.NewObjectID(),
		ConfigurationEntryID: uuid.NewString(),
		HookRecordID:         message.RecordID,
	}

	// Create Hook Record!
	createdHookRecord := storage.HookRecord{
		ID:                              primitive.NewObjectID(),
		HookRecordID:                    message.RecordID,
		CreatorDID:                      message.Processing.AuthorDID,
		InitialHookConfigurationEntryID: configEntry.ConfigurationEntryID,
		LatestHookConfigurationEntryID:  configEntry.ConfigurationEntryID,
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
	_, sp := observability.Tracer.Start(ctx, "UpdateHook")
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
	configEntry := storage.HookConfigurationEntry{
		Message:              message,
		ID:                   primitive.NewObjectID(),
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

func (hookService HookService) GetHooksForCollection(ctx context.Context, request *services.GetHooksForCollectionRequest) (*services.GetHooksForCollectionResponse, error) {
	// tracing
	_, sp := observability.Tracer.Start(ctx, "GetHooksForCollection")
	defer sp.End()

	//TODO implement me
	panic("implement me")
}

func (hookService HookService) NotifyHooksOfCollectionEvent(ctx context.Context, request *services.NotifyHooksOfCollectionEventRequest) (*services.NotifyHooksOfCollectionEventResponse, error) {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "NotifyHooksOfCollectionEvent")
	defer sp.End()

	//TODO implement me
	panic("implement me")
}

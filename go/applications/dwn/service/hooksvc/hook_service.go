package hooksvc

import (
	"context"
	"github.com/openreserveio/dwn/go/framework/events"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
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
	//TODO implement me
	panic("implement me")
}

func (hookService HookService) GetHooksForCollection(ctx context.Context, request *services.GetHooksForCollectionRequest) (*services.GetHooksForCollectionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (hookService HookService) NotifyHooksOfCollectionEvent(ctx context.Context, in *services.NotifyHooksOfCollectionEventRequest) (*services.NotifyHooksOfCollectionEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

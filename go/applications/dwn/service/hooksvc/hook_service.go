package hooksvc

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
)

type HookService struct {
	services.UnimplementedHookServiceServer
}

func CreateHookService() (HookService, error) {
	return HookService{}, nil
}

func (hookService HookService) RegisterHook(ctx context.Context, request *services.RegisterHookRequest) (*services.RegisterHookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (hookService HookService) GetHooksForCollection(ctx context.Context, request *services.GetHooksForCollectionRequest) (*services.GetHooksForCollectionResponse, error) {
	//TODO implement me
	panic("implement me")
}

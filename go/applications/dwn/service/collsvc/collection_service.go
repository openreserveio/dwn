package collsvc

import (
	"context"
	"github.com/openreserveio/dwn/go/applications/dwn/service/collsvc/collection"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
)

type CollectionService struct {
	services.UnimplementedCollectionServiceServer
	CollectionStore storage.CollectionStore
}

func CreateCollectionService(collectionStoreConnectionURI string) (*CollectionService, error) {

	// Setup Collection Store
	colLStore, err := docdbstore.CreateCollectionDocumentDBStore(collectionStoreConnectionURI)
	if err != nil {
		log.Fatal("Unable to connect to collections store:  %v", err)
		return nil, err
	}

	collService := CollectionService{
		CollectionStore: colLStore,
	}

	return &collService, nil

}

func (c CollectionService) StoreCollection(ctx context.Context, request *services.StoreCollectionRequest) (*services.StoreCollectionResponse, error) {

	response := services.StoreCollectionResponse{}
	newOrExistingID, err := collection.StoreCollection(c.CollectionStore, request.SchemaURI, request.CollectionItemId, request.CollectionItem)
	if err != nil {
		log.Error("Storing the collection failed:  %v", err)
		response.Status.Status = services.Status_ERROR
		response.Status.Details = err.Error()
		return &response, nil
	}

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	response.CollectionId = newOrExistingID

	return &response, nil

}

func (c CollectionService) FindCollection(ctx context.Context, request *services.FindCollectionRequest) (*services.FindCollectionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CollectionService) CreateSchema(ctx context.Context, request *services.CreateSchemaRequest) (*services.CreateSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CollectionService) ValidateCollection(ctx context.Context, request *services.ValidateCollectionRequest) (*services.ValidateCollectionResponse, error) {

	return &services.ValidateCollectionResponse{
		Status: &services.CommonStatus{Status: services.Status_OK},
	}, nil

}

func (c CollectionService) InvalidateSchema(ctx context.Context, request *services.InvalidateSchemaRequest) (*services.InvalidateSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CollectionService) mustEmbedUnimplementedCollectionServiceServer() {
	//TODO implement me
	panic("implement me")
}

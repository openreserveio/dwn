package collsvc

import (
	"context"
	"encoding/json"
	"github.com/openreserveio/dwn/go/applications/dwn/service/collsvc/collection"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
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
	var collectionMessage model.Message
	err := json.Unmarshal(request.Message, &collectionMessage)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_WRITE ||
		collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_COMMIT ||
		collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_DELETE {

		result, err := collection.StoreCollection(c.CollectionStore, &collectionMessage)
		if err != nil {
			response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
			return &response, nil
		}

		if result.Status == "UNSUPPORTED_METHOD" {
			response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "UNSUPPORTED METHOD"}
			return &response, nil
		} else if result.Status == "ERROR" {
			response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: result.Error.Error()}
			return &response, nil
		}

		response.Status = &services.CommonStatus{Status: services.Status_OK}
		response.RecordId = result.RecordID

	}

	return &response, nil

}

func (c CollectionService) FindCollection(ctx context.Context, request *services.FindCollectionRequest) (*services.FindCollectionResponse, error) {

	response := services.FindCollectionResponse{}

	// TODO: We are only doing single record finds right now
	if request.QueryType != services.QueryType_SINGLE_COLLECTION_BY_ID_SCHEMA_URI {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "TODO: We are only doing single record finds right now"}
		return &response, nil
	}

	response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND, Details: "TODO: We are only doing single record finds right now"}
	return &response, nil

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

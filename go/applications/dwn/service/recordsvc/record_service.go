package recordsvc

import (
	"context"
	"encoding/json"
	"github.com/openreserveio/dwn/go/applications/dwn/service/recordsvc/record"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
)

type RecordService struct {
	services.UnimplementedRecordServiceServer
	CollectionStore storage.RecordStore
}

func CreateRecordService(recordStoreConnectionURI string) (*RecordService, error) {

	// Setup Collection Store
	colLStore, err := docdbstore.CreateRecordDocumentDBStore(recordStoreConnectionURI)
	if err != nil {
		log.Fatal("Unable to connect to collections store:  %v", err)
		return nil, err
	}

	collService := RecordService{
		CollectionStore: colLStore,
	}

	return &collService, nil

}

func (c RecordService) StoreRecord(ctx context.Context, request *services.StoreRecordRequest) (*services.StoreRecordResponse, error) {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "StoreRecord")
	defer sp.End()

	response := services.StoreRecordResponse{}
	var collectionMessage model.Message
	err := json.Unmarshal(request.Message, &collectionMessage)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_WRITE ||
		collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_COMMIT ||
		collectionMessage.Descriptor.Method == model.METHOD_COLLECTIONS_DELETE {

		result, err := record.StoreCollection(ctx, c.CollectionStore, &collectionMessage)
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

func (c RecordService) FindRecord(ctx context.Context, request *services.FindRecordRequest) (*services.FindRecordResponse, error) {

	// tracing
	_, sp := observability.Tracer.Start(ctx, "FindRecord")
	defer sp.End()

	response := services.FindRecordResponse{}

	// TODO: We are only doing single record finds right now
	if request.QueryType == services.QueryType_SINGLE_RECORD_BY_ID_SCHEMA_URI {

		result, err := record.FindCollectionBySchemaAndRecordID(ctx, c.CollectionStore, request.SchemaURI, request.RecordId)
		if err != nil {
			response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
			return &response, nil
		}

		switch result.Status {

		case "OK":
			response.Status = &services.CommonStatus{Status: services.Status_OK}
			response.SchemaURI = request.SchemaURI
			response.Writers = []string{result.LatestEntry.Processing.AuthorDID}
			response.Readers = []string{result.LatestEntry.Processing.AuthorDID, result.LatestEntry.Processing.RecipientDID}
			response.IsPublished = result.LatestEntry.Descriptor.Published

			latestEntryJsonBytes, _ := json.Marshal(result.LatestEntry)
			response.RecordItem = latestEntryJsonBytes

			return &response, nil

		case "NOT_FOUND":
			response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND}
			return &response, nil

		default:
			response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "Some kind of error"}
			return &response, nil
		}

	}

	response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: "Unsupport query type"}
	return &response, nil

}

func (c RecordService) CreateSchema(ctx context.Context, request *services.CreateSchemaRequest) (*services.CreateSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c RecordService) ValidateRecord(ctx context.Context, request *services.ValidateRecordRequest) (*services.ValidateRecordResponse, error) {

	return &services.ValidateRecordResponse{
		Status: &services.CommonStatus{Status: services.Status_OK},
	}, nil

}

func (c RecordService) InvalidateSchema(ctx context.Context, request *services.InvalidateSchemaRequest) (*services.InvalidateSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c RecordService) mustEmbedUnimplementedCollectionServiceServer() {
	//TODO implement me
	panic("implement me")
}

package recordsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/service/recordsvc/record"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RecordService struct {
	services.UnimplementedRecordServiceServer
	RecordStore storage.RecordStore
}

func CreateRecordService(recordStoreConnectionURI string) (*RecordService, error) {

	// Setup Collection Store
	colLStore, err := docdbstore.CreateRecordDocumentDBStore(recordStoreConnectionURI)
	if err != nil {
		log.Fatal("Unable to connect to collections store:  %v", err)
		return nil, err
	}

	collService := RecordService{
		RecordStore: colLStore,
	}

	return &collService, nil

}

func (c RecordService) StoreRecord(ctx context.Context, request *services.StoreRecordRequest) (*services.StoreRecordResponse, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "StoreRecord")
	defer sp.End()

	response := services.StoreRecordResponse{}
	var collectionMessage model.Message
	err := json.Unmarshal(request.Message, &collectionMessage)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if collectionMessage.Descriptor.Method == model.METHOD_RECORDS_WRITE ||
		collectionMessage.Descriptor.Method == model.METHOD_RECORDS_COMMIT ||
		collectionMessage.Descriptor.Method == model.METHOD_RECORDS_DELETE {

		result, err := record.StoreRecord(ctx, c.RecordStore, &collectionMessage)
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
		response.InitialEntry = result.InitialEntry

	}

	return &response, nil

}

func (c RecordService) FindRecord(ctx context.Context, request *services.FindRecordRequest) (*services.FindRecordResponse, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "FindRecord")
	defer sp.End()

	response := services.FindRecordResponse{}

	// TODO: We are only doing single record finds right now
	if request.QueryType == services.QueryType_SINGLE_RECORD_BY_ID_SCHEMA_URI {

		result, err := record.FindRecordBySchemaAndRecordID(ctx, c.RecordStore, request.SchemaURI, request.RecordId)
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

	} else if request.QueryType == services.QueryType_SINGLE_RECORD_BY_ID_FOR_COMMIT {

		result, err := record.FindRecordForCommit(ctx, c.RecordStore, request.SchemaURI, request.RecordId)
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

func (rs RecordService) Query(ctx context.Context, request *services.QueryRecordRequest) (*services.QueryRecordResponse, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordsvc.Query")
	defer sp.End()

	response := services.QueryRecordResponse{}

	var queryRecordMessage model.Message
	err := json.Unmarshal(request.Message, &queryRecordMessage)
	if err != nil {
		sp.RecordError(err)
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	sp.AddEvent("Execute Query Method for Record")
	queryResponseMessage, err := record.RecordQuery(ctx, rs.RecordStore, &queryRecordMessage)
	if err != nil {
		sp.RecordError(err)
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	if queryResponseMessage == nil {
		sp.AddEvent("No message found by query")
		response.Status = &services.CommonStatus{Status: services.Status_NOT_FOUND}
		return &response, nil
	}

	queryResponseMessageBytes, _ := json.Marshal(queryResponseMessage)
	sp.AddEvent("Query returned a message", trace.WithAttributes(attribute.String("message-json", string(queryResponseMessageBytes))))

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	response.Message = queryResponseMessageBytes

	return &response, nil

}

func (rs RecordService) Write(ctx context.Context, request *services.WriteRecordRequest) (*services.WriteRecordResponse, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordsvc.Write")
	defer sp.End()

	response := services.WriteRecordResponse{}

	var writeRecordMessage model.Message
	err := json.Unmarshal(request.Message, &writeRecordMessage)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	sp.AddEvent("Execute Write Method for Record")
	recordId, err := record.RecordWrite(ctx, rs.RecordStore, &writeRecordMessage)
	if err != nil {
		sp.RecordError(err)
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}
	sp.AddEvent(fmt.Sprintf("Write Method for Record Executed Successfully:  %s", recordId))

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	return &response, nil

}

func (rs RecordService) Commit(ctx context.Context, request *services.CommitRecordRequest) (*services.CommitRecordResponse, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordsvc.Commit")
	defer sp.End()

	response := services.CommitRecordResponse{}

	var commitRecordMessage model.Message
	err := json.Unmarshal(request.Message, &commitRecordMessage)
	if err != nil {
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}

	sp.AddEvent("Execute Commit Method for Record")
	err = record.RecordCommit(ctx, rs.RecordStore, &commitRecordMessage)
	if err != nil {
		sp.RecordError(err)
		response.Status = &services.CommonStatus{Status: services.Status_ERROR, Details: err.Error()}
		return &response, nil
	}
	sp.AddEvent("Commit Method for Record Executed Successfully")

	response.Status = &services.CommonStatus{Status: services.Status_OK}
	return &response, nil

}

func (rs RecordService) Delete(ctx context.Context, request *services.DeleteRecordRequest) (*services.DeleteRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (rs RecordService) mustEmbedUnimplementedRecordServiceServer() {
	//TODO implement me
	panic("implement me")
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

package collections

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func CollectionsQuery(ctx context.Context, collSvcClient services.CollectionServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "CollectionsQuery")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	// First, make sure authorizations are valid and correct for this message
	if !model.VerifyAuthorization(message) {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify authorization(s)."}
		return messageResultObj
	}

	// Next, find the schema and make sure it has been registered
	schemaUri := message.Descriptor.Filter.Schema
	if schemaUri == "" {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a CollectionsQuery"}
		return messageResultObj
	}

	// Is this for a specific record, or all records since context?
	if message.Descriptor.Filter.RecordID == "" {
		// TODO:  For now we are only going to allow single message access
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "TODO:  For now we are only going to allow single message access"}
		return messageResultObj
	}

	// Get the collection item
	req := services.FindCollectionRequest{
		QueryType:    services.QueryType_SINGLE_COLLECTION_BY_ID_SCHEMA_URI,
		RecordId:     message.Descriptor.Filter.RecordID,
		SchemaURI:    message.Descriptor.Filter.Schema,
		RequestorDID: message.Processing.AuthorDID,
	}

	findCollResponse, err := collSvcClient.FindCollection(ctx, &req)
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	switch findCollResponse.Status.Status {

	case services.Status_OK:
		messageResultEntry := model.MessageResultEntry{
			Result: findCollResponse.CollectionItem,
		}
		messageResultObj.Entries = append(messageResultObj.Entries, messageResultEntry)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
		return messageResultObj

	case services.Status_NOT_FOUND:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusNotFound, Detail: fmt.Sprintf("Unable to locate collection item with record ID %s", message.RecordID)}
		return messageResultObj

	case services.Status_INVALID_AUTHORIZATION:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized}
		return messageResultObj

	case services.Status_ERROR:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: findCollResponse.Status.Details}
		return messageResultObj

	default:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: "Default Error"}
		return messageResultObj

	}

	return messageResultObj

}

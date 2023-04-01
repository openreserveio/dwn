package records

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func RecordsQuery(ctx context.Context, recordServiceClient services.RecordServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, sp := observability.Tracer.Start(ctx, "api.records.RecordsQuery")
	defer sp.End()

	var messageResultObj model.MessageResultObject

	// First, make sure authorizations are valid and correct for this message
	sp.AddEvent("Verifying authorizations for message")
	if !model.VerifyAuthorization(message) {
		sp.AddEvent("Unable to verify authorizations.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify authorization(s)."}
		return messageResultObj
	}

	// Next, find the schema and make sure it has been registered
	sp.AddEvent(fmt.Sprintf("Verifying schema is registered: %s", message.Descriptor.Filter.Schema))
	schemaUri := message.Descriptor.Filter.Schema
	if schemaUri == "" {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a RecordsQuery"}
		return messageResultObj
	}

	// Is this for a specific record, or all records since context?
	if message.Descriptor.Filter.RecordID == "" {
		// TODO:  For now we are only going to allow single message access
		sp.AddEvent("No RecordID in DescriptorFilter")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "TODO:  For now we are only going to allow single message access.  Please specify a record ID in the DescriptorFilter for this query."}
		return messageResultObj
	}

	// Execute the Query
	sp.AddEvent("Calling Query at recordsvc")
	queryRecordBytes, _ := json.Marshal(message)
	queryRecordRequest := services.QueryRecordRequest{
		Message: queryRecordBytes,
	}

	queryRecordResponse, err := recordServiceClient.Query(ctx, &queryRecordRequest)
	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	// Process the response
	switch queryRecordResponse.Status.Status {

	case services.Status_OK:
		messageResultEntry := model.MessageResultEntry{
			Result: queryRecordResponse.Message,
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
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: queryRecordResponse.Status.Details}
		return messageResultObj

	default:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: "Default Error"}
		return messageResultObj

	}

	return messageResultObj

}

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

func RecordsWrite(ctx context.Context, recordSvcClient services.RecordServiceClient, hookServiceClient services.HookServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "RecordsWrite")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	// First, make sure attestations are valid and correct for this message
	// TODO:  Deal with whitelisting, blacklisting, authentication requirements
	if !model.VerifyAttestation(message) {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify attestation(s)."}
		return messageResultObj
	}

	// Make sure authorizations are valid for messages that are writes to existing records
	// Check for existing record
	childSpan.AddEvent("Looking for existing collection")
	findCollResp, err := recordSvcClient.FindRecord(ctx, &services.FindRecordRequest{QueryType: services.QueryType_SINGLE_RECORD_BY_ID_SCHEMA_URI, SchemaURI: message.Descriptor.Schema, RecordId: message.RecordID})
	if err != nil {
		childSpan.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}
	childSpan.AddEvent("Looked for existing collection")

	// If no record was found, then we don't need to authorize
	//var foundCollMessage model.Message
	//json.Unmarshal(findCollResp.CollectionItem, &foundCollMessage)

	switch findCollResp.Status.Status {

	case services.Status_NOT_FOUND:
	// No need to authorize!

	case services.Status_OK:
		// We found a record.  Must authorize
		if !model.VerifyAuthorization(message) {
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify authorization(s)."}
			return messageResultObj
		}

		authorized := false
		for _, v := range findCollResp.Writers {
			if v == message.Processing.AuthorDID {
				authorized = true
			}
		}

		if !authorized {
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Author is not authorized to write to this record."}
			return messageResultObj
		}

	}

	// Next, find the schema and make sure it has been registered
	schemaUri := message.Descriptor.Schema
	if schemaUri == "" {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a RecordsWrite"}
		return messageResultObj
	}

	// Validate given collection data validates against given schema
	validateCollRequest := services.ValidateRecordRequest{
		SchemaURI: schemaUri,
		Document:  []byte(message.Data),
	}

	validateCollResponse, err := recordSvcClient.ValidateRecord(ctx, &validateCollRequest)
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	if validateCollResponse.Status.Status != services.Status_OK {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: validateCollResponse.Status.Details}
		return messageResultObj
	}

	// Store and process the message if it passes schema validation!
	encodedMsg, _ := json.Marshal(message)
	storeReq := services.StoreRecordRequest{
		Message: encodedMsg,
	}
	storeResp, err := recordSvcClient.StoreRecord(ctx, &storeReq)
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if storeResp.Status.Status != services.Status_OK {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: storeResp.Status.Details}
		return messageResultObj
	}

	// Only notify others if this is an initial entry.
	// otherwise, wait for the commit to do a RECORD_CHANGED or DELETED
	childSpan.AddEvent(fmt.Sprintf("Is this an initial entry?  %v", storeResp.InitialEntry))
	if storeResp.InitialEntry {

		// Notify Others
		childSpan.AddEvent("INITIAL ENTRY.  Notify others of a RECORD_CREATED event!")
		notify := services.NotifyHooksOfRecordEventRequest{
			RecordId:        storeResp.RecordId,
			Protocol:        message.Descriptor.Protocol,
			ProtocolVersion: message.Descriptor.ProtocolVersion,
			Schema:          message.Descriptor.Schema,
			RecordEventType: services.RecordEventType_RECORD_CREATED,
		}
		hookServiceClient.NotifyHooksOfRecordEvent(ctx, &notify)

	}

	existingOrNewId := storeResp.RecordId
	messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
	messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(existingOrNewId)})

	return messageResultObj

}

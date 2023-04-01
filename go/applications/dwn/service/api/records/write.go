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
	ctx, sp := observability.Tracer.Start(ctx, "api.records.RecordsWrite")
	defer sp.End()

	var messageResultObj model.MessageResultObject

	// First, make sure attestations are valid and correct for this message
	// TODO:  Deal with whitelisting, blacklisting, authentication requirements
	sp.AddEvent("Verifying attestations for message")
	if !model.VerifyAttestation(message) {
		sp.AddEvent("Unable to verify attestations.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify attestation(s)."}
		return messageResultObj
	}

	// Make sure authorizations are valid for messages that are writes to existing records
	// Check for existing record
	sp.AddEvent("Looking for existing record")
	findRecordResp, err := recordSvcClient.FindRecord(ctx, &services.FindRecordRequest{QueryType: services.QueryType_SINGLE_RECORD_BY_ID_SCHEMA_URI, SchemaURI: message.Descriptor.Schema, RecordId: message.RecordID})
	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	// If no record was found, then we don't need to authorize
	//var foundCollMessage model.Message
	//json.Unmarshal(findRecordResp.CollectionItem, &foundCollMessage)

	switch findRecordResp.Status.Status {

	case services.Status_NOT_FOUND:
		sp.AddEvent("Record not found.  No need to authorize.")
	// No need to authorize!

	case services.Status_OK:
		// We found a record.  Must authorize
		sp.AddEvent("Record found.  Authorizing.")
		if !model.VerifyAuthorization(message) {
			sp.AddEvent("Unable to verify authorizations.")
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify authorization(s)."}
			return messageResultObj
		}

		sp.AddEvent("Authorizations verified.  Checking if author is authorized to write to this record.")

		authorized := false
		for _, v := range findRecordResp.Writers {
			if v == message.Processing.AuthorDID {
				sp.AddEvent("Author is authorized to write to this record.")
				authorized = true
			}
		}

		if !authorized {
			sp.AddEvent("Author is not authorized to write to this record.")
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Author is not authorized to write to this record."}
			return messageResultObj
		}

	}

	// Next, find the schema and make sure it has been registered
	sp.AddEvent("Ensuring schema is registered")
	schemaUri := message.Descriptor.Schema
	if schemaUri == "" {
		sp.AddEvent("Schema URI was not found to be registered.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a RecordsWrite"}
		return messageResultObj
	}

	// Validate given record data validates against given schema
	sp.AddEvent("Validating record data against schema")
	validateRecordRequest := services.ValidateRecordRequest{
		SchemaURI: schemaUri,
		Document:  []byte(message.Data),
	}

	validateRecordResponse, err := recordSvcClient.ValidateRecord(ctx, &validateRecordRequest)
	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	if validateRecordResponse.Status.Status != services.Status_OK {
		sp.AddEvent("Record data did not validate against schema.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: validateRecordResponse.Status.Details}
		return messageResultObj
	}

	// Store and process the message if it passes schema validation!
	sp.AddEvent("Execute Write method at recordsvc")
	encodedMsg, _ := json.Marshal(message)
	writeRecordRequest := services.WriteRecordRequest{
		Message: encodedMsg,
	}

	writeRecordResponse, err := recordSvcClient.Write(ctx, &writeRecordRequest)

	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if writeRecordResponse.Status.Status != services.Status_OK {
		sp.AddEvent(fmt.Sprintf("Unable to write record due to error:  %s", writeRecordResponse.Status.Details))
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: writeRecordResponse.Status.Details}
		return messageResultObj
	}
	sp.AddEvent("Record write was successful!")

	// Only notify others if this is an initial entry.
	// otherwise, wait for the commit to do a RECORD_CHANGED or DELETED
	sp.AddEvent("Notifying interested parties")
	notify := services.NotifyHooksOfRecordEventRequest{
		RecordId:        message.RecordID,
		Protocol:        message.Descriptor.Protocol,
		ProtocolVersion: message.Descriptor.ProtocolVersion,
		Schema:          message.Descriptor.Schema,
		RecordEventType: services.RecordEventType_RECORD_CREATED,
	}
	hookServiceClient.NotifyHooksOfRecordEvent(ctx, &notify)

	messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
	messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(message.RecordID)})

	return messageResultObj

}

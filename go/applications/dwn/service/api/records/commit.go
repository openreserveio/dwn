package records

import (
	"context"
	"encoding/json"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func RecordsCommit(ctx context.Context, collSvcClient services.RecordServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "RecordsCommit")
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
	findCollResp, err := collSvcClient.FindRecord(ctx, &services.FindRecordRequest{QueryType: services.QueryType_SINGLE_RECORD_BY_ID_SCHEMA_URI, SchemaURI: message.Descriptor.Schema, RecordId: message.Descriptor.ParentID})
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	// There can be no COMMIT to a record that does not exist or is erroring out
	switch findCollResp.Status.Status {

	case services.Status_NOT_FOUND:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Cannot COMMIT to a record that does not exist."}
		return messageResultObj

	case services.Status_ERROR:
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: findCollResp.Status.Details}
		return messageResultObj

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
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Author is not authorized to COMMIT to this record."}
			return messageResultObj
		}

	}

	// Store and process the message if it is authorized!
	encodedMsg, _ := json.Marshal(message)
	storeReq := services.StoreRecordRequest{
		Message: encodedMsg,
	}
	storeResp, err := collSvcClient.StoreRecord(ctx, &storeReq)
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if storeResp.Status.Status != services.Status_OK {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: storeResp.Status.Details}
		return messageResultObj
	}

	existingOrNewId := storeResp.RecordId
	messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
	messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(existingOrNewId)})

	return messageResultObj
}

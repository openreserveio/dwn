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

func RecordsCommit(ctx context.Context, recordServiceClient services.RecordServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, sp := observability.Tracer().Start(ctx, "api.records.RecordsCommit")
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
	// Check for existing record by the parent ID for commit
	sp.AddEvent(fmt.Sprintf("Looking for existing record for commit:  %s", message.RecordID))
	findRecordResponse, err := recordServiceClient.FindRecord(ctx, &services.FindRecordRequest{QueryType: services.QueryType_SINGLE_RECORD_BY_ID_FOR_COMMIT, SchemaURI: message.Descriptor.Schema, RecordId: message.Descriptor.ParentID})
	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	// There can be no COMMIT to a record that does not exist or is erroring out
	switch findRecordResponse.Status.Status {

	case services.Status_NOT_FOUND:
		sp.AddEvent("Record not found.  No need to authorize.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Cannot COMMIT to a record that does not exist."}
		return messageResultObj

	case services.Status_ERROR:
		sp.AddEvent("Error finding record.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: findRecordResponse.Status.Details}
		return messageResultObj

	case services.Status_OK:
		// We found a record.  Must authorize
		sp.AddEvent("Record found.  Verifying authorization(s).")
		if !model.VerifyAuthorization(message) {
			sp.AddEvent("Unable to verify authorization(s).")
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify authorization(s)."}
			return messageResultObj
		}

		// Make sure the author is authorized to write to this record
		authorized := false
		for _, v := range findRecordResponse.Writers {
			if v == message.Processing.AuthorDID {
				sp.AddEvent("Author is authorized to COMMIT to this record.")
				authorized = true
			}
		}

		if !authorized {
			sp.AddEvent("Author is not authorized to COMMIT to this record.")
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Author is not authorized to COMMIT to this record."}
			return messageResultObj
		}

	}

	// Store and process the message if it is authorized!
	encodedMsg, _ := json.Marshal(message)
	commitRecordRequest := services.CommitRecordRequest{
		Message: encodedMsg,
	}

	sp.AddEvent("Committing record...")
	commitRecordResponse, err := recordServiceClient.Commit(ctx, &commitRecordRequest)
	if err != nil {
		sp.RecordError(err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if commitRecordResponse.Status.Status != services.Status_OK {
		sp.AddEvent("Error committing record.")
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: commitRecordResponse.Status.Details}
		return messageResultObj
	}

	recordID := message.RecordID
	messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
	messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(recordID)})

	return messageResultObj
}

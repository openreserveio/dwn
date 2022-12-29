package collections

import (
	"context"
	"encoding/json"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
)

func CollectionsWrite(collSvcClient services.CollectionServiceClient, message *model.Message) model.MessageResultObject {

	var messageResultObj model.MessageResultObject

	// First, make sure attestations are valid and correct for this message
	if !model.VerifyAttestation(message) {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify attestation(s)."}
		return messageResultObj
	}

	// Next, find the schema and make sure it has been registered
	schemaUri := message.Descriptor.Schema
	if schemaUri == "" {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a CollectionsWrite"}
		return messageResultObj
	}

	// Validate given collection data validates against given schema
	validateCollRequest := services.ValidateCollectionRequest{
		SchemaURI: schemaUri,
		Document:  []byte(message.Data),
	}

	validateCollResponse, err := collSvcClient.ValidateCollection(context.Background(), &validateCollRequest)
	if err != nil {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError, Detail: err.Error()}
		return messageResultObj
	}

	if validateCollResponse.Status.Status != services.Status_OK {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: validateCollResponse.Status.Details}
		return messageResultObj
	}

	// Store the collection if it passes schema validation!
	encodedMsg, _ := json.Marshal(message)
	storeReq := services.StoreCollectionRequest{
		Message: encodedMsg,
	}
	storeResp, err := collSvcClient.StoreCollection(context.Background(), &storeReq)
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

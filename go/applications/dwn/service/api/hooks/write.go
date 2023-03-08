package hooks

import (
	"context"
	"encoding/json"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func HooksWrite(ctx context.Context, hookServiceClient services.HookServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "HooksWrite")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	// First, make sure attestations are valid and correct for this message
	// TODO:  Deal with whitelisting, blacklisting, authentication requirements
	if !model.VerifyAttestation(message) {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify attestation(s)."}
		return messageResultObj
	}

	// Make sure authorizations are valid for messages that are writes to existing records
	// TODO: Authorization

	// See if there's an existing Hook record and be sure we're authorized
	// otherwise, create a new one
	getHooksRequest := services.GetHooksForRecordRequest{
		RecordId:        message.RecordID,
		Protocol:        message.Descriptor.Protocol,
		ProtocolVersion: message.Descriptor.ProtocolVersion,
		Schema:          message.Descriptor.Schema,
	}
	getHooksResponse, err := hookServiceClient.GetHooksForRecord(ctx, &getHooksRequest)
	if err != nil {
		log.Error("Internal Server Error:  %v", err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if len(getHooksResponse.HookDefinitions) > 0 {

		// This is an update!
		messageBytes, _ := json.Marshal(&message)
		updateReq := services.UpdateHookRequest{Message: messageBytes}
		updateRes, err := hookServiceClient.UpdateHook(ctx, &updateReq)
		if err != nil {
			log.Error("Internal Server Error:  %v", err)
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
			return messageResultObj
		}

		if updateRes.Status.Status != services.Status_OK {
			messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: updateRes.Status.Details}
			return messageResultObj
		}

		messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
		messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(message.RecordID)})

		return messageResultObj

	}

	// We are creating!
	messageBytes, _ := json.Marshal(&message)
	registerHookReq := services.RegisterHookRequest{Message: messageBytes}
	registerHookRes, err := hookServiceClient.RegisterHook(ctx, &registerHookReq)
	if err != nil {
		log.Error("Internal Server Error:  %v", err)
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusInternalServerError}
		return messageResultObj
	}

	if registerHookRes.Status.Status != services.Status_OK {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: registerHookRes.Status.Details}
		return messageResultObj
	}

	messageResultObj.Status = model.ResponseStatus{Code: http.StatusOK}
	messageResultObj.Entries = append(messageResultObj.Entries, model.MessageResultEntry{Result: []byte(message.RecordID)})

	return messageResultObj

}

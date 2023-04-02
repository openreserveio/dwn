package api

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api/hooks"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api/records"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

type FeatureRouter struct {
	model.FeatureDetection
	MaximumTimeoutSeconds int
	RecordServiceClient   services.RecordServiceClient
	HookServiceClient     services.HookServiceClient
}

func CreateFeatureRouter(recordSvcClient services.RecordServiceClient, hooksvcClient services.HookServiceClient, maxTimeoutSeconds int) (*FeatureRouter, error) {

	return &FeatureRouter{
		FeatureDetection:      model.CurrentFeatureDetection,
		MaximumTimeoutSeconds: maxTimeoutSeconds,
		RecordServiceClient:   recordSvcClient,
		HookServiceClient:     hooksvcClient,
	}, nil

}

func (fr *FeatureRouter) Route(ctx context.Context, requestObject *model.RequestObject) (interface{}, error) {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "Route Operation")
	defer childSpan.End()

	// Setup Response Object
	responseObject := model.ResponseObject{}
	responseObject.Replies = make([]model.MessageResultObject, len(requestObject.Messages))

	// Process Messages and append responses to responseObject
	for idx, message := range requestObject.Messages {

		log.Info("Processing message %d", idx)
		res, err := fr.processMessage(ctx, idx, &message)
		if err != nil {
			return nil, err
		}
		responseObject.Replies[idx] = *res.MessageResult

	}

	responseObject.Status = model.ResponseStatus{Code: http.StatusOK}
	return &responseObject, nil

}

// Meant to be run in a goroutine
type MessageProcResult struct {
	Index         int
	MessageResult *model.MessageResultObject
}

// idx is for ordering, MessageProcResult wraps the messageresult object and the index for the responseobject
func (fr *FeatureRouter) processMessage(ctx context.Context, idx int, message *model.Message) (*MessageProcResult, error) {

	// Instrumentation
	ctx, sp := observability.Tracer.Start(ctx, "api.FeatureRouter.processMessage")
	defer sp.End()

	// Support Simple Test Messages
	if message.RecordID == "TEST" && message.Data == "TEST" {

		// This is a test message
		procComm := &MessageProcResult{
			Index: idx,
			MessageResult: &model.MessageResultObject{
				Status:  model.ResponseStatus{Code: http.StatusOK},
				Entries: nil,
			},
		}
		return procComm, nil

	}

	var messageResult model.MessageResultObject

	// Route logic based on Interface and Method
	switch message.Descriptor.Interface {

	case model.INTERFACE_RECORDS:

		switch message.Descriptor.Method {
		case model.METHOD_RECORDS_QUERY:
			sp.AddEvent("Start Records Query")
			messageResult = records.RecordsQuery(ctx, fr.RecordServiceClient, message)

		case model.METHOD_RECORDS_WRITE:
			sp.AddEvent("Start Records Write")
			messageResult = records.RecordsWrite(ctx, fr.RecordServiceClient, fr.HookServiceClient, message)

		case model.METHOD_RECORDS_COMMIT:
			sp.AddEvent("Start Records Commit")
			messageResult = records.RecordsCommit(ctx, fr.RecordServiceClient, message)

		case model.METHOD_RECORDS_DELETE:
			sp.AddEvent("Start Records Delete")
			messageResult = records.RecordsDelete(ctx, fr.RecordServiceClient, message)
		}

	case model.INTERFACE_HOOKS:

		switch message.Descriptor.Method {
		case model.METHOD_HOOKS_WRITE:
			sp.AddEvent("Start Hooks Write")
			messageResult = hooks.HooksWrite(ctx, fr.HookServiceClient, message)

		case model.METHOD_HOOKS_QUERY:
			sp.AddEvent("Start Hooks Query")
			messageResult = hooks.HooksQuery(ctx, fr.HookServiceClient, message)

		case model.METHOD_HOOKS_DELETE:
			sp.AddEvent("Start Hooks Delete")
			messageResult = hooks.HooksDelete(ctx, fr.HookServiceClient, message)

		default:
			sp.AddEvent("Bad Method")
			messageResult = model.MessageResultObject{Status: model.ResponseStatus{Code: http.StatusBadRequest, Detail: fmt.Sprintf("We do not yet support message method: %s", message.Descriptor.Method)}}

		}
	}

	// Reply back to channel
	procResult := MessageProcResult{
		Index:         idx,
		MessageResult: &messageResult,
	}
	return &procResult, nil

}

// Fills missing responses with generic timeouts
func (fr *FeatureRouter) genericTimeouts(responseObject *model.ResponseObject) {
	for idx, reply := range responseObject.Replies {
		if reply.Status.Code == 0 {
			// never got a status code of 200 or other error > 0
			responseObject.Replies[idx] = model.MessageResultObject{
				Status:  model.ResponseStatus{Code: http.StatusRequestTimeout, Detail: "This message could not be processed within a generous timeout"},
				Entries: nil,
			}
		}
	}
}

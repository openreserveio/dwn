package api

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api/collections"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api/hooks"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

type FeatureRouter struct {
	model.FeatureDetection
	MaximumTimeoutSeconds   int
	CollectionServiceClient services.CollectionServiceClient
	HookServiceClient       services.HookServiceClient
}

func CreateFeatureRouter(collsvcClient services.CollectionServiceClient, hooksvcClient services.HookServiceClient, maxTimeoutSeconds int) (*FeatureRouter, error) {

	return &FeatureRouter{
		FeatureDetection:        model.CurrentFeatureDetection,
		MaximumTimeoutSeconds:   maxTimeoutSeconds,
		CollectionServiceClient: collsvcClient,
		HookServiceClient:       hooksvcClient,
	}, nil

}

func (fr *FeatureRouter) Route(ctx context.Context, requestObject *model.RequestObject) (interface{}, error) {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "Route Operation")
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
	_, childSpan := observability.Tracer.Start(ctx, "processMessage")
	defer childSpan.End()

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

	// Route logic based on Method
	switch message.Descriptor.Method {

	case model.METHOD_COLLECTIONS_QUERY:
		childSpan.AddEvent("Start Collections Query")
		messageResult = collections.CollectionsQuery(ctx, fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_WRITE:
		childSpan.AddEvent("Start Collections Write")
		messageResult = collections.CollectionsWrite(ctx, fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_COMMIT:
		childSpan.AddEvent("Start Collections Commit")
		messageResult = collections.CollectionsCommit(ctx, fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_DELETE:
		childSpan.AddEvent("Start Collections Delete")
		messageResult = collections.CollectionsDelete(ctx, fr.CollectionServiceClient, message)

	case model.METHOD_HOOKS_WRITE:
		childSpan.AddEvent("Start Hooks Write")
		messageResult = hooks.HooksWrite(ctx, fr.HookServiceClient, message)

	default:
		childSpan.AddEvent("Bad Method")
		messageResult = model.MessageResultObject{Status: model.ResponseStatus{Code: http.StatusBadRequest, Detail: fmt.Sprintf("We do not yet support message method: %s", message.Descriptor.Method)}}

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

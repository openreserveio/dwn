package hooks

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
)

func HooksQuery(ctx context.Context, hookServiceClient services.HookServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "HooksQuery")
	defer childSpan.End()

	messageResultObj := model.MessageResultObject{}

	return messageResultObj
}

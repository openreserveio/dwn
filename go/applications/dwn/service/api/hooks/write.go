package hooks

import (
	"context"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
)

func HooksWrite(ctx context.Context, message *model.Message) model.MessageResultObject {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "CollectionsCommit")
	defer childSpan.End()

	return model.MessageResultObject{}

}

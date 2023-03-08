package hooks

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
)

func HooksQuery(ctx context.Context, hookServiceClient services.HookServiceClient, message *model.Message) model.MessageResultObject {

	messageResultObj := model.MessageResultObject{}

	return messageResultObj
}

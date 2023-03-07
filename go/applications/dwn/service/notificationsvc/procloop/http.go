package procloop

import (
	"github.com/openreserveio/dwn/go/generated/events"
	"github.com/openreserveio/dwn/go/log"
)

func NotifyHTTP(event *events.Event) error {

	log.Info("HTTP NOTIFY EVENT:  %v", event)
	return nil

}

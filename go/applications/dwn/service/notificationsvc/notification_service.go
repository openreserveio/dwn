package notificationsvc

import (
	"fmt"
	"github.com/openreserveio/dwn/go/framework/events"
	evtypes "github.com/openreserveio/dwn/go/generated/events"
	"github.com/openreserveio/dwn/go/log"
	"os"
)

type NotificationService struct {
	EventHub                   *events.EventHub
	NotificationQueueName      string
	NotificationErrorQueueName string
}

func CreateNotificationService(queueServiceUri string, notificationQueueName string) (*NotificationService, error) {

	eventHub, err := events.CreateEventHub(queueServiceUri)
	if err != nil {
		log.Fatal("Unable to create NotificationService due to problem with EventHub:  %v", err)
		os.Exit(1)
	}

	notifyService := NotificationService{
		EventHub:                   eventHub,
		NotificationQueueName:      notificationQueueName,
		NotificationErrorQueueName: fmt.Sprintf("%sErrors", notificationQueueName),
	}

	return &notifyService, nil

}

func (notifyService *NotificationService) StartProcessingLoop(processor func(event *evtypes.Event) error) error {

	notifyCallbackChan := notifyService.EventHub.Subscribe(notifyService.NotificationQueueName)

	for {

		msg := <-notifyCallbackChan

		// Decode Event Message
		eventMessage := notifyService.EventHub.DecodeEventMessage(msg.Data)
		if eventMessage == nil {
			log.Error("The event message received could not be decoded.")
			msg.Nak()
			continue
		}
		msg.Ack()

		// PROCESS!
		err := processor(eventMessage)
		if err != nil {
			// Send to error queue
			log.Error("Message unable to be processed, sending to Error Queue (%s):  %v", notifyService.NotificationErrorQueueName, err)
			notifyService.EventHub.Publish(notifyService.NotificationErrorQueueName, msg.Data)
		}

		// Next message

	}

}

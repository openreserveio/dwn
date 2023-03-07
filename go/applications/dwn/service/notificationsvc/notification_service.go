package notificationsvc

import (
	"github.com/openreserveio/dwn/go/framework/events"
	"github.com/openreserveio/dwn/go/log"
	"os"
)

type NotificationService struct {
	EventHub              *events.EventHub
	NotificationQueueName string
}

func CreateNotificationService(queueServiceUri string, notificationQueueName string) (*NotificationService, error) {

	eventHub, err := events.CreateEventHub(queueServiceUri)
	if err != nil {
		log.Fatal("Unable to create NotificationService due to problem with EventHub:  %v", err)
		os.Exit(1)
	}

	notifyService := NotificationService{
		EventHub: eventHub,
	}

	return &notifyService, nil

}

func (notifyService *NotificationService) StartProcessingLoop(processor func()) error {

}

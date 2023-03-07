package notificationsvc

import (
	"context"
	"errors"
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/applications/dwn/service/notificationsvc/procloop"
	"github.com/openreserveio/dwn/go/log"
	"os"
)

func Start(ctx context.Context, config configuration.Configuration) error {

	queueServiceUri := config.GetQueueServiceConnectionURI()
	notifyCallbackQueueName := config.GetNotifyCallbackQueueName()

	notificationService, err := CreateNotificationService(queueServiceUri, notifyCallbackQueueName)
	if err != nil {
		log.Fatal("An error occured while creating the notification service:  %v", err)
		os.Exit(1)
	}

	log.Info("Starting the Notification Service Loop")
	err = notificationService.StartProcessingLoop(procloop.ProcessMessage)
	if err != nil {
		log.Fatal("Error occurred in processing loop:  %v", err)
		return err
	}

	return errors.New("Exiting Processing Loop")

}

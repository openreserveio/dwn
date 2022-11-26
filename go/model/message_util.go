package model

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func CreateMessage(authorDID string, recipientDID string, dataFormat string, data []byte, methodName string, recordId string, schema string) *Message {

	// Verify Message Name

	// If there is data, base64 encode it in string form
	var encodedData string = ""
	if data != nil {
		encodedData = base64.URLEncoding.EncodeToString(data)
	}

	// Start the Message
	message := Message{
		RecordID: recordId,
		Data:     encodedData,
		Processing: MessageProcessing{
			Nonce:        uuid.NewString(),
			AuthorDID:    authorDID,
			RecipientDID: recipientDID,
		},
	}

	// create the descriptor
	var dataCID string = ""
	if message.Data != "" {
		dataCID = CreateDataCID(message.Data)
	}

	messageDesc := Descriptor{
		Method:     methodName,
		DataCID:    dataCID,
		DataFormat: dataFormat,
		Schema:     schema,
	}
	message.Descriptor = messageDesc

	return &message

}

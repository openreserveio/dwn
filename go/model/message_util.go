package model

import (
	"encoding/base64"
	"github.com/google/uuid"
	"time"
)

const (
	METHOD_COLLECTIONS_QUERY  = "CollectionsQuery"
	METHOD_COLLECTIONS_WRITE  = "CollectionsWrite"
	METHOD_COLLECTIONS_COMMIT = "CollectionsCommit"
	METHOD_COLLECTIONS_DELETE = "CollectionsDelete"
)

func CreateCollectionsWriteMessage(authorDID string, recipientDID string, protocol string, protocolVersion string, schema string, dataFormat string, data []byte) *Message {

	// If there is data, base64 encode it in string form
	var encodedData string = ""
	if data != nil {
		encodedData = base64.URLEncoding.EncodeToString(data)
	}

	// create the data CID if there is data
	var dataCID string = ""
	if encodedData != "" {
		dataCID = CreateDataCID(encodedData)
	}

	// Descriptor
	var messageDescriptorCID string = ""
	messageDesc := Descriptor{
		Method:          METHOD_COLLECTIONS_WRITE,
		DataCID:         dataCID,
		DataFormat:      dataFormat,
		ParentID:        "",
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		Schema:          schema,
		CommitStrategy:  "",
		DateCreated:     time.Now(),
	}
	messageDescriptorCID = CreateDescriptorCID(messageDesc)

	// Message Processing
	var processingCID string = ""
	messageProcessing := MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    authorDID,
		RecipientDID: recipientDID,
	}
	processingCID = CreateProcessingCID(messageProcessing)

	recordId := CreateRecordCID(messageDescriptorCID, processingCID)

	msg := Message{
		RecordID:   recordId,
		ContextID:  "",
		Data:       encodedData,
		Processing: messageProcessing,
		Descriptor: messageDesc,
	}

	return &msg

}

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

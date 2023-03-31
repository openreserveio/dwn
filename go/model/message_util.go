package model

import (
	"encoding/base64"
	"github.com/google/uuid"
	"time"
)

const (
	METHOD_RECORDS_QUERY  = "RecordsQuery"
	METHOD_RECORDS_WRITE  = "RecordsWrite"
	METHOD_RECORDS_COMMIT = "RecordsCommit"
	METHOD_RECORDS_DELETE = "RecordsDelete"

	METHOD_HOOKS_WRITE  = "HooksWrite"
	METHOD_HOOKS_QUERY  = "HooksQuery"
	METHOD_HOOKS_DELETE = "HooksDelete"
)

type ProtocolDefinition struct {
	ContextID       string
	Protocol        string
	ProtocolVersion string
}

func CreateQueryRecordsMessage(schemaUri string, recordId string, protocolDef *ProtocolDefinition, requestorDID string) *Message {

	queryDescriptor := Descriptor{
		Method: METHOD_RECORDS_QUERY,
		Filter: DescriptorFilter{
			RecordID:        recordId,
			Schema:          schemaUri,
			Protocol:        protocolDef.Protocol,
			ProtocolVersion: protocolDef.ProtocolVersion,
		},
	}

	queryMessageProcessing := MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    requestorDID,
		RecipientDID: requestorDID,
	}

	queryMessage := Message{
		ContextID:  protocolDef.ContextID,
		Processing: queryMessageProcessing,
		Descriptor: queryDescriptor,
	}

	return &queryMessage

}

func CreateUpdateRecordsWriteMessage(authorDID string, recipientDID string, previousRecordId string, protocolDef *ProtocolDefinition, schemaUri string, dataFormat string, data []byte) *Message {

	// TODO:  How to deal with Context IDs?

	// Encode your data
	dataEncoded := base64.RawURLEncoding.EncodeToString(data)
	dataCID := CreateDataCID(dataEncoded)

	descriptor := Descriptor{
		Method:          METHOD_RECORDS_WRITE,
		DataCID:         dataCID,
		DataFormat:      dataFormat,
		ParentID:        previousRecordId,
		Protocol:        protocolDef.Protocol,
		ProtocolVersion: protocolDef.ProtocolVersion,
		Schema:          schemaUri,
		CommitStrategy:  "",
		Published:       false,
		DateCreated:     time.Now().Format(time.RFC3339),
		DatePublished:   nil,
	}

	processing := MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    authorDID,
		RecipientDID: recipientDID,
	}

	descCID := CreateDescriptorCID(descriptor)
	mpCID := CreateProcessingCID(processing)
	recordCID := CreateRecordCID(descCID, mpCID)

	message := Message{
		RecordID:   recordCID,
		ContextID:  protocolDef.ContextID,
		Data:       dataEncoded,
		Processing: processing,
		Descriptor: descriptor,
	}

	return &message

}

func CreateRecordsCommitMessage(previousWriteRecordId string, schemaUrl string, committerDID string) *Message {

	// TODO:  How to deal with Context IDs?

	descriptor := Descriptor{
		Method:         METHOD_RECORDS_COMMIT,
		ParentID:       previousWriteRecordId,
		Schema:         schemaUrl,
		CommitStrategy: "",
		DateCreated:    time.Now().Format(time.RFC3339),
	}

	processing := MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    committerDID,
		RecipientDID: committerDID,
	}

	descCID := CreateDescriptorCID(descriptor)
	mpCID := CreateProcessingCID(processing)
	recordCID := CreateRecordCID(descCID, mpCID)

	message := Message{
		RecordID:   recordCID,
		Processing: processing,
		Descriptor: descriptor,
	}

	return &message

}

func CreateInitialRecordsWriteMessage(authorDID string, recipientDID string, protocolDef *ProtocolDefinition, schema string, dataFormat string, data []byte) *Message {

	// Encode your data
	dataEncoded := base64.RawURLEncoding.EncodeToString(data)
	dataCID := CreateDataCID(dataEncoded)

	descriptor := Descriptor{
		Method:          METHOD_RECORDS_WRITE,
		DataCID:         dataCID,
		DataFormat:      dataFormat,
		ParentID:        "",
		Protocol:        protocolDef.Protocol,
		ProtocolVersion: protocolDef.ProtocolVersion,
		Schema:          schema,
		CommitStrategy:  "",
		Published:       false,
		DateCreated:     time.Now().Format(time.RFC3339),
		DatePublished:   nil,
	}

	processing := MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    authorDID,
		RecipientDID: recipientDID,
	}

	descCID := CreateDescriptorCID(descriptor)
	mpCID := CreateProcessingCID(processing)
	recordCID := CreateRecordCID(descCID, mpCID)

	message := Message{
		RecordID:   recordCID,
		ContextID:  protocolDef.ContextID,
		Data:       dataEncoded,
		Processing: processing,
		Descriptor: descriptor,
	}

	return &message

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

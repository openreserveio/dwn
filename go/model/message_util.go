package model

import (
	"encoding/base64"
	"github.com/google/uuid"
	"time"
)

const (
	INTERFACE_RECORDS     = "Records"
	METHOD_RECORDS_QUERY  = "Query"
	METHOD_RECORDS_WRITE  = "Write"
	METHOD_RECORDS_COMMIT = "Commit"
	METHOD_RECORDS_DELETE = "Delete"

	INTERFACE_HOOKS     = "Hooks"
	METHOD_HOOKS_WRITE  = "Write"
	METHOD_HOOKS_QUERY  = "Query"
	METHOD_HOOKS_DELETE = "Delete"
)

type ProtocolDefinition struct {
	ContextID       string
	Protocol        string
	ProtocolVersion string
}

func CreateQueryRecordsMessage(schemaUri string, recordId string, protocolDef *ProtocolDefinition, requestorDID string) *Message {

	queryDescriptor := Descriptor{
		Interface: INTERFACE_RECORDS,
		Method:    METHOD_RECORDS_QUERY,
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

func CreateUpdateRecordsWriteMessage(authorDID string, recipientDID string, logicalRecordId string, protocolDef *ProtocolDefinition, schemaUri string, dataFormat string, data []byte) *Message {

	// TODO:  How to deal with Context IDs?

	// Encode your data
	dataEncoded := base64.RawURLEncoding.EncodeToString(data)
	dataCID := CreateDataCID(dataEncoded)

	descriptor := Descriptor{
		Interface:       INTERFACE_RECORDS,
		Method:          METHOD_RECORDS_WRITE,
		DataCID:         dataCID,
		DataFormat:      dataFormat,
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
	recordCID := CreateRecordCID(descCID)
	descriptor.ParentID = recordCID

	message := Message{
		RecordID:   logicalRecordId,
		ContextID:  protocolDef.ContextID,
		Data:       dataEncoded,
		Processing: processing,
		Descriptor: descriptor,
	}

	return &message

}

func CreateRecordsCommitMessage(logicalRecordId string, schemaUrl string, committerDID string) *Message {

	// TODO:  How to deal with Context IDs?

	descriptor := Descriptor{
		Interface:      INTERFACE_RECORDS,
		Method:         METHOD_RECORDS_COMMIT,
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
	recordCID := CreateRecordCID(descCID)
	descriptor.ParentID = recordCID

	message := Message{
		RecordID:   logicalRecordId,
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
		Interface:       INTERFACE_RECORDS,
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
	recordCID := CreateRecordCID(descCID)

	message := Message{
		RecordID:   recordCID,
		ContextID:  protocolDef.ContextID,
		Data:       dataEncoded,
		Processing: processing,
		Descriptor: descriptor,
	}

	return &message

}

func CreateMessage(authorDID string, recipientDID string, dataFormat string, data []byte, interfaceName string, methodName string, recordId string, schema string) *Message {

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
		Interface:  interfaceName,
		Method:     methodName,
		DataCID:    dataCID,
		DataFormat: dataFormat,
		Schema:     schema,
	}
	message.Descriptor = messageDesc

	return &message

}

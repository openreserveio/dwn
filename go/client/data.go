package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
	"net/http"
)

func (client *DWNClient) GetData(schemaUrl string, recordId string, requestorIdentity *Identity) ([]byte, string, error) {

	queryDescriptor := model.Descriptor{
		Method: model.METHOD_RECORDS_QUERY,
		Filter: model.DescriptorFilter{
			RecordID: recordId,
			Schema:   schemaUrl,
		},
	}

	queryMessageProcessing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    requestorIdentity.DID,
		RecipientDID: requestorIdentity.DID,
	}

	queryMessage := model.Message{
		ContextID:  "",
		Processing: queryMessageProcessing,
		Descriptor: queryDescriptor,
	}

	authorization := model.CreateAuthorization(&queryMessage, *requestorIdentity.Keypair.PrivateKey)
	attestation := model.CreateAttestation(&queryMessage, *requestorIdentity.Keypair.PrivateKey)
	queryMessage.Attestation = attestation
	queryMessage.Authorization = authorization

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, queryMessage)

	res, err := resty.New().R().
		SetBody(ro).
		SetHeader("Content-Type", "application/json").
		Post(client.DWNUrlBase)

	if err != nil {
		return nil, "", err
	}

	var responseObject model.ResponseObject
	err = json.Unmarshal(res.Body(), &responseObject)
	if err != nil {
		return nil, "", err
	}

	// TODO: Change this return object -- shouldn't be a message entry from storage package
	var entry storage.MessageEntry
	json.Unmarshal(responseObject.Replies[0].Entries[0].Result, &entry)
	data, err := base64.RawURLEncoding.DecodeString(entry.Data)
	if err != nil {
		return nil, "", err
	}

	return data, entry.Descriptor.DataFormat, nil

}

func (client *DWNClient) SaveData(schemaUrl string, data []byte, dataFormat string, dataAuthor *Identity, dataRecipient *Identity) (string, error) {

	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}

	recordsWriteMessage := model.CreateInitialRecordsWriteMessage(dataAuthor.DID, dataRecipient.DID, &protocolDef, schemaUrl, dataFormat, data)

	attestation := model.CreateAttestation(recordsWriteMessage, *dataAuthor.Keypair.PrivateKey)
	recordsWriteMessage.Attestation = attestation

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, *recordsWriteMessage)

	responseObject, err := client.CallDWNHTTP(ro)
	if err != nil {
		return "", err
	}

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

func (client *DWNClient) UpdateData(schemaUrl string, parentRecordId string, data []byte, dataFormat string, dataUpdater *Identity) (string, error) {

	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}
	message := model.CreateUpdateRecordsWriteMessage(dataUpdater.DID, dataUpdater.DID, parentRecordId, &protocolDef, schemaUrl, dataFormat, data)

	attestation := model.CreateAttestation(message, *dataUpdater.Keypair.PrivateKey)
	message.Attestation = attestation

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, *message)

	responseObject, err := client.CallDWNHTTP(ro)
	if err != nil {
		return "", err
	}

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

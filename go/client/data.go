package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
)

func (client *DWNClient) GetData(schemaUrl string, recordId string, requestorIdentity *Identity) (*model.Message, []byte, string, error) {

	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}
	queryMessage := model.CreateQueryRecordsMessage(schemaUrl, recordId, &protocolDef, requestorIdentity.DID)

	authorization := model.CreateAuthorization(queryMessage, *requestorIdentity.Keypair.PrivateKey)
	attestation := model.CreateAttestation(queryMessage, *requestorIdentity.Keypair.PrivateKey)
	queryMessage.Attestation = attestation
	queryMessage.Authorization = authorization

	responseObject, err := client.CallDWNHTTP(queryMessage)
	if err != nil {
		return nil, nil, "", err
	}

	// If the da†a was not found, return nil, "", nil
	if len(responseObject.Replies) > 0 {

		if responseObject.Replies[0].Status.Code == http.StatusNotFound {
			return nil, nil, "", nil
		}

	}

	var data []byte
	var dataFormat string
	var entry model.Message

	if len(responseObject.Replies[0].Entries) > 0 {

		// TODO: Change this return object -- shouldn't be a message entry from storage package
		json.Unmarshal(responseObject.Replies[0].Entries[0].Result, &entry)
		data, err = base64.RawURLEncoding.DecodeString(entry.Data)
		if err != nil {
			return nil, nil, "", err
		}
		dataFormat = entry.Descriptor.DataFormat

	}

	return &entry, data, dataFormat, nil

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

	authorization := model.CreateAuthorization(recordsWriteMessage, *dataAuthor.Keypair.PrivateKey)
	recordsWriteMessage.Authorization = authorization

	responseObject, err := client.CallDWNHTTP(recordsWriteMessage)
	if err != nil {
		return "", err
	}

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	if len(responseObject.Replies) < 1 {
		return "", errors.New("Wrong number of message replies.")
	}

	if responseObject.Replies[0].Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Replies[0].Status.Detail)
	}

	if len(responseObject.Replies[0].Entries) < 1 {
		return "", errors.New("No Reply Entries as expected")
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

func (client *DWNClient) UpdateData(schemaUrl string, umbrellaRecordId string, latestRecordWriteId string, data []byte, dataFormat string, dataUpdater *Identity) (string, error) {

	// Create a Write pointing back to the previous latest entry,
	// then do a commit on it
	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}

	// Query for the latest
	latestDataMessage, _, _, err := client.GetData(schemaUrl, umbrellaRecordId, dataUpdater)
	if err != nil {
		return "", err
	}
	if latestDataMessage == nil {
		return "", errors.New("No latest data message found")
	}

	writeMessage := model.CreateUpdateRecordsWriteMessage(dataUpdater.DID, dataUpdater.DID, latestRecordWriteId, &protocolDef, schemaUrl, dataFormat, data)
	writeAttestation := model.CreateAttestation(writeMessage, *dataUpdater.Keypair.PrivateKey)
	writeMessage.Attestation = writeAttestation
	writeAuthorization := model.CreateAuthorization(writeMessage, *dataUpdater.Keypair.PrivateKey)
	writeMessage.Authorization = writeAuthorization

	// Create the corresponding COMMIT
	commitMessage := model.CreateRecordsCommitMessage(writeMessage.RecordID, writeMessage.Descriptor.Schema, dataUpdater.DID)
	commitAttestation := model.CreateAttestation(commitMessage, *dataUpdater.Keypair.PrivateKey)
	commitMessage.Attestation = commitAttestation
	commitAuthorization := model.CreateAuthorization(commitMessage, *dataUpdater.Keypair.PrivateKey)
	commitMessage.Authorization = commitAuthorization

	responseObject, err := client.CallDWNHTTP(writeMessage, commitMessage)
	if err != nil {
		return "", err
	}

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	if len(responseObject.Replies) < 2 {
		return "", errors.New("Wrong number of message replies.")
	}

	if responseObject.Replies[1].Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Replies[1].Status.Detail)
	}

	if len(responseObject.Replies[0].Entries) < 1 {
		return "", errors.New("An updated record ID was not returned")
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

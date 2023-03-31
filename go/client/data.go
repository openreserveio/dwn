package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
	"net/http"
)

func (client *DWNClient) GetData(schemaUrl string, recordId string, requestorIdentity *Identity) ([]byte, string, error) {

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

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, *queryMessage)

	responseObject, err := client.CallDWNHTTP(ro)
	if err != nil {
		return nil, "", err
	}

	// If the da†a was not found, return nil, "", nil
	if len(responseObject.Replies) > 0 {

		if responseObject.Replies[0].Status.Code == http.StatusNotFound {
			return nil, "", nil
		}

	}

	var data []byte
	var dataFormat string
	if len(responseObject.Replies[0].Entries) > 0 {

		// TODO: Change this return object -- shouldn't be a message entry from storage package
		var entry storage.MessageEntry
		json.Unmarshal(responseObject.Replies[0].Entries[0].Result, &entry)
		data, err = base64.RawURLEncoding.DecodeString(entry.Data)
		if err != nil {
			return nil, "", err
		}
		dataFormat = entry.Descriptor.DataFormat

	}

	return data, dataFormat, nil

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

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, *recordsWriteMessage)

	responseObject, err := client.CallDWNHTTP(ro)
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

func (client *DWNClient) UpdateData(schemaUrl string, parentRecordId string, data []byte, dataFormat string, dataUpdater *Identity) (string, error) {

	// Create a Write pointing back to the previous latest entry,
	// then do a commit on it
	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}
	writeMessage := model.CreateUpdateRecordsWriteMessage(dataUpdater.DID, dataUpdater.DID, parentRecordId, &protocolDef, schemaUrl, dataFormat, data)
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

	// Append both
	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, *writeMessage)
	ro.Messages = append(ro.Messages, *commitMessage)

	responseObject, err := client.CallDWNHTTP(ro)
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

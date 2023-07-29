package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	requestorDIDDocument, err := model.ResolveDID(requestorIdentity.DID)
	requestorVerificationId := fmt.Sprintf("%s%s", requestorDIDDocument.VerificationMethod[0].Controller, requestorDIDDocument.VerificationMethod[0].ID)
	if err != nil {
		return nil, nil, "", errors.New(fmt.Sprintf("Unable to resolve requestor identity DID:  %v", err))
	}

	authorization := model.CreateAuthorization(queryMessage, requestorVerificationId, requestorIdentity.Keypair.PublicKey, requestorIdentity.Keypair.PrivateKey)
	attestation := model.CreateAttestation(queryMessage, requestorVerificationId, requestorIdentity.Keypair.PublicKey, requestorIdentity.Keypair.PrivateKey)
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

	dataAuthorDIDDocument, err := model.ResolveDID(dataAuthor.DID)
	dataAuthorVerificationId := fmt.Sprintf("%s%s", dataAuthorDIDDocument.VerificationMethod[0].Controller, dataAuthorDIDDocument.VerificationMethod[0].ID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to resolve data author DID:  %v", err))
	}

	attestation := model.CreateAttestation(recordsWriteMessage, dataAuthorVerificationId, dataAuthor.Keypair.PublicKey, dataAuthor.Keypair.PrivateKey)
	recordsWriteMessage.Attestation = attestation

	authorization := model.CreateAuthorization(recordsWriteMessage, dataAuthorVerificationId, dataAuthor.Keypair.PublicKey, dataAuthor.Keypair.PrivateKey)
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

func (client *DWNClient) UpdateData(schemaUrl string, logicalRecordId string, data []byte, dataFormat string, dataUpdater *Identity) (string, error) {

	// Create a Write pointing back to the previous latest entry,
	// then do a commit on it
	protocolDef := model.ProtocolDefinition{
		ContextID:       "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
	}

	// Query for the latest
	latestDataMessage, _, _, err := client.GetData(schemaUrl, logicalRecordId, dataUpdater)
	if err != nil {
		return "", err
	}
	if latestDataMessage == nil {
		return "", errors.New("No latest data message found")
	}

	// Resolve the data updater did
	dataUpdaterDIDDocument, err := model.ResolveDID(dataUpdater.DID)
	dataUpdaterVerificationId := fmt.Sprintf("%s%s", dataUpdaterDIDDocument.VerificationMethod[0].Controller, dataUpdaterDIDDocument.VerificationMethod[0].ID)

	writeMessage := model.CreateUpdateRecordsWriteMessage(dataUpdater.DID, dataUpdater.DID, logicalRecordId, &protocolDef, schemaUrl, dataFormat, data)
	writeAttestation := model.CreateAttestation(writeMessage, dataUpdaterVerificationId, dataUpdater.Keypair.PublicKey, dataUpdater.Keypair.PrivateKey)
	writeMessage.Attestation = writeAttestation
	writeAuthorization := model.CreateAuthorization(writeMessage, dataUpdaterVerificationId, dataUpdater.Keypair.PublicKey, dataUpdater.Keypair.PrivateKey)
	writeMessage.Authorization = writeAuthorization

	// Create the corresponding COMMIT
	commitMessage := model.CreateRecordsCommitMessage(logicalRecordId, writeMessage.Descriptor.Schema, dataUpdater.DID)
	commitAttestation := model.CreateAttestation(commitMessage, dataUpdaterVerificationId, dataUpdater.Keypair.PublicKey, dataUpdater.Keypair.PrivateKey)
	commitMessage.Attestation = commitAttestation
	commitAuthorization := model.CreateAuthorization(commitMessage, dataUpdaterVerificationId, dataUpdater.Keypair.PublicKey, dataUpdater.Keypair.PrivateKey)
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

package client

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
	"time"
)

func (client *DWNClient) SaveHookForSchemaAndProtocol(schemaUri string, protocol string, protocolVersion string, callbackUri string, requestor *Identity) (string, error) {

	log.Debug("Saving hook for schema %s, protocol & version %s %s, callback %s", schemaUri, protocol, protocolVersion, callbackUri)
	descriptor := model.Descriptor{
		Interface:       model.INTERFACE_HOOKS,
		Method:          model.METHOD_HOOKS_WRITE,
		ParentID:        "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
		Schema:          schemaUri,
		URI:             callbackUri,
		Published:       false,
		DateCreated:     time.Now().Format(time.RFC3339),
		DatePublished:   nil,

		// Filter for the data record ID
		Filter: model.DescriptorFilter{
			Schema:          schemaUri,
			Protocol:        protocol,
			ProtocolVersion: protocolVersion,
		},
	}

	processing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    requestor.DID,
		RecipientDID: requestor.DID,
	}

	message := model.Message{
		RecordID:   uuid.NewString(),
		ContextID:  "",
		Processing: processing,
		Descriptor: descriptor,
	}

	resolverDIDDocument, err := model.ResolveDID(requestor.DID)
	if err != nil {
		return "", errors.New("Unable to resolve requestor identity DID")
	}
	resolverDIDVerifyId := resolverDIDDocument.VerificationMethod[0].ID
	attestation := model.CreateAttestation(&message, resolverDIDVerifyId, requestor.Keypair.PublicKey, requestor.Keypair.PrivateKey)
	message.Attestation = attestation

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, message)

	res, err := resty.New().R().
		SetBody(ro).
		SetHeader(HEADER_CONTENT_TYPE_KEY, HEADER_CONTENT_TYPE_APPLICATION_JSON).
		Post(client.DWNUrlBase)

	if err != nil {
		return "", err
	}
	if !res.IsSuccess() {
		return "", errors.New("Unable to create hook")
	}

	var responseObject model.ResponseObject
	err = json.Unmarshal(res.Body(), &responseObject)

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	if len(responseObject.Replies) == 0 {
		return "", errors.New("No response replies")
	}

	if len(responseObject.Replies[0].Entries) == 0 {
		return "", errors.New("No response entries")
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

func (client *DWNClient) SaveHookForRecord(schemaUri string, dataRecordId string, callbackUri string, requestor *Identity) (string, error) {

	log.Debug("Saving hook for schema %s, data record %s, callback %s", schemaUri, dataRecordId, callbackUri)
	descriptor := model.Descriptor{
		Interface:       model.INTERFACE_HOOKS,
		Method:          model.METHOD_HOOKS_WRITE,
		ParentID:        "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
		Schema:          schemaUri,
		URI:             callbackUri,
		Published:       false,
		DateCreated:     time.Now().Format(time.RFC3339),
		DatePublished:   nil,

		// Filter for the data record ID
		Filter: model.DescriptorFilter{
			Schema:   schemaUri,
			RecordID: dataRecordId,
		},
	}

	processing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    requestor.DID,
		RecipientDID: requestor.DID,
	}

	message := model.Message{
		RecordID:   uuid.NewString(),
		ContextID:  "",
		Processing: processing,
		Descriptor: descriptor,
	}

	resolverDIDDocument, err := model.ResolveDID(requestor.DID)
	if err != nil {
		return "", errors.New("Unable to resolve requestor identity DID")
	}
	requestorDIDVerifyId := resolverDIDDocument.VerificationMethod[0].ID

	attestation := model.CreateAttestation(&message, requestorDIDVerifyId, requestor.Keypair.PublicKey, requestor.Keypair.PrivateKey)
	message.Attestation = attestation

	ro := model.RequestObject{}
	ro.Messages = append(ro.Messages, message)

	res, err := resty.New().R().
		SetBody(ro).
		SetHeader(HEADER_CONTENT_TYPE_KEY, HEADER_CONTENT_TYPE_APPLICATION_JSON).
		Post(client.DWNUrlBase)

	if err != nil {
		return "", err
	}
	if !res.IsSuccess() {
		return "", errors.New("Unable to create hook")
	}

	var responseObject model.ResponseObject
	err = json.Unmarshal(res.Body(), &responseObject)

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	if len(responseObject.Replies) == 0 {
		return "", errors.New("No response replies")
	}

	if len(responseObject.Replies[0].Entries) == 0 {
		return "", errors.New("No response entries")
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

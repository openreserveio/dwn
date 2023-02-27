package client

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
	"time"
)

func (client *DWNClient) SaveHook(schemaUri string, dataRecordId string, callbackUri string, requestor *Identity) (string, error) {

	descriptor := model.Descriptor{
		Method:          model.METHOD_HOOKS_WRITE,
		ParentID:        "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
		Schema:          schemaUri,
		Published:       false,
		DateCreated:     time.Now(),
		DatePublished:   nil,
	}

	processing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    requestor.DID,
		RecipientDID: requestor.DID,
	}

	message := model.Message{
		RecordID:   dataRecordId,
		ContextID:  "",
		Processing: processing,
		Descriptor: descriptor,
	}

	attestation := model.CreateAttestation(&message, *requestor.Keypair.PrivateKey)
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

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

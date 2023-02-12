package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
	"time"
)

const (
	HEADER_CONTENT_TYPE_KEY              = "Content-Type"
	HEADER_CONTENT_TYPE_APPLICATION_JSON = "application/json"
)

type DWNClient struct {
	DWNUrlBase      string
	Protocol        string
	ProtocolVersion string
}

func CreateDWNClient(urlBase string) *DWNClient {

	return &DWNClient{
		DWNUrlBase: urlBase,
	}

}

func CreateDWNClientForProtocol(urlBase string, protocol string, protocolVersion string) *DWNClient {

	return &DWNClient{
		DWNUrlBase:      urlBase,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
	}

}

func (client *DWNClient) SaveData(schemaUrl string, data []byte, dataFormat string, dataAuthor *Identity, dataRecipient *Identity) (string, error) {

	dataEncoded := base64.RawURLEncoding.EncodeToString(data)

	descriptor := model.Descriptor{
		Method:          model.METHOD_COLLECTIONS_WRITE,
		DataCID:         model.CreateDataCID(dataEncoded),
		DataFormat:      dataFormat,
		ParentID:        "",
		Protocol:        client.Protocol,
		ProtocolVersion: client.ProtocolVersion,
		Schema:          schemaUrl,
		CommitStrategy:  "",
		Published:       false,
		DateCreated:     time.Now(),
		DatePublished:   nil,
	}

	processing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    dataAuthor.DID,
		RecipientDID: dataRecipient.DID,
	}

	descriptorCID := model.CreateDescriptorCID(descriptor)
	processingCID := model.CreateProcessingCID(processing)
	recordId := model.CreateRecordCID(descriptorCID, processingCID)

	message := model.Message{
		RecordID:   recordId,
		ContextID:  "",
		Data:       dataEncoded,
		Processing: processing,
		Descriptor: descriptor,
	}

	attestation := model.CreateAttestation(&message, *dataAuthor.Keypair.PrivateKey)
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
		return "", errors.New("Unable to create data")
	}

	var responseObject model.ResponseObject
	err = json.Unmarshal(res.Body(), &responseObject)

	if responseObject.Status.Code != http.StatusOK {
		return "", errors.New(responseObject.Status.Detail)
	}

	return string(responseObject.Replies[0].Entries[0].Result), nil

}

func (client *DWNClient) UpdateData(schemaUrl string, primaryIdentifier string, data []byte, dataFormat string, dataUpdater *Identity) error {

	return nil
}

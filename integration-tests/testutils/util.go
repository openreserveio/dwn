package testutils

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
	"github.com/openreserveio/dwn/go/model"
)

func CreateMessage(authorDID string, recipientDID string, dataFormat string, data []byte, methodName string) *model.Message {

	// Verify Message Name

	// If there is data, base64 encode it in string form
	var encodedData string = ""
	if data != nil {
		encodedData = base64.URLEncoding.EncodeToString(data)
	}

	// Start the Message
	message := model.Message{
		Data: encodedData,
	}

	// create the descriptor
	var dataCID string = ""
	if message.Data != "" {
		dataCID = CreateDataCID(message.Data)
	}

	messageDesc := model.Descriptor{
		Method:     methodName,
		DataCID:    dataCID,
		DataFormat: dataFormat,
	}
	message.Descriptor = messageDesc

	processing := model.MessageProcessing{
		Nonce:        uuid.NewString(),
		AuthorDID:    authorDID,
		RecipientDID: recipientDID,
	}

	descriptorCID := CreateDescriptorCID(messageDesc)
	processingCID := CreateProcessingCID(processing)
	message.RecordID = CreateRecordCID(descriptorCID, processingCID)

	return &message

}

func CreateCIDFromNode(node datamodel.Node) cid.Cid {

	var buf bytes.Buffer
	dagcbor.Encode(node, &buf)

	cidPrefix := cid.Prefix{
		Version: 1,
	}
	cid, _ := cidPrefix.Sum(buf.Bytes())
	return cid

}

func CreateRecordCID(descriptorCID string, processingCID string) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "descriptorCid", qp.String(descriptorCID))
		qp.MapEntry(ma, "processingCid", qp.String(processingCID))
	})
	if err != nil {
		return ""
	}

	return CreateCIDFromNode(d).String()

}

func CreateDataCID(data string) string {

	d, err := qp.BuildList(basicnode.Prototype.Any, 1, func(la datamodel.ListAssembler) {
		qp.ListEntry(la, qp.String(data))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

func CreateDescriptorCID(descriptor model.Descriptor) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "method", qp.String(descriptor.Method))
		qp.MapEntry(ma, "schema", qp.String(descriptor.Schema))
		qp.MapEntry(ma, "dataCid", qp.String(descriptor.DataCID))
		qp.MapEntry(ma, "dataFormat", qp.String(descriptor.DataFormat))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

func CreateProcessingCID(mp model.MessageProcessing) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "nonce", qp.String(mp.Nonce))
		qp.MapEntry(ma, "author", qp.String(mp.AuthorDID))
		qp.MapEntry(ma, "recipient", qp.String(mp.RecipientDID))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

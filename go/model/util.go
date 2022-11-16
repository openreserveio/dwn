package model

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func CreateMessage(authorDID string, recipientDID string, dataFormat string, data []byte, methodName string) *Message {

	// Verify Message Name

	// If there is data, base64 encode it in string form
	var encodedData string = ""
	if data != nil {
		encodedData = base64.URLEncoding.EncodeToString(data)
	}

	// Start the Message
	message := Message{
		Data:       encodedData,
		Processing: MessageProcessing{},
		Descriptor: Descriptor{},
	}

	// create the descriptor
	var dataCID string = ""
	if message.Data != "" {
		d, err := qp.BuildList(basicnode.Prototype.Any, 1, func(la datamodel.ListAssembler) {
			qp.ListEntry(la, qp.String(message.Data))
		})
		if err != nil {
			return nil
		}

		var buf bytes.Buffer
		dagcbor.Encode(d, &buf)
		dataCID = buf.String()

	}
	messageDesc := Descriptor{
		Nonce:      uuid.New().String(),
		Method:     methodName,
		DataCID:    dataCID,
		DataFormat: dataFormat,
	}

	message.Descriptor = messageDesc

	return &message

}

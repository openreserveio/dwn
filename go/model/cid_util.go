package model

import (
	"bytes"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	mh "github.com/multiformats/go-multihash"
)

var DescriptorSchemaType schema.Type

func init() {

	descriptorSchemaBytes := []byte(`

type Descriptor struct {
  method String
  dataCid String
  dataFormat String
  filter DescriptorFilter
  parentId String
  protocol String
  protocolVersion String
  schema String
  commitStrategy String
  published Bool
  dateCreated String
  datePublished nullable String 
  uri String
}

type DescriptorFilter struct {
  schema String
  recordId String
  contextId String
  protocol String
  protocolVersion String
  dataFormat String
  dateSort String
}


`)

	ts, err := ipld.LoadSchemaBytes(descriptorSchemaBytes)
	if err != nil {
		panic(err)
	}
	DescriptorSchemaType = ts.TypeByName("Descriptor")

}

func CreateCIDFromNode(node datamodel.Node) cid.Cid {

	var buf bytes.Buffer
	dagcbor.Encode(node, &buf)

	cidBuilder := cid.V1Builder{
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	contentCID, err := cidBuilder.Sum(buf.Bytes())
	if err != nil {
		return cid.Cid{}
	}
	return contentCID

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

	d := bindnode.Wrap(&data, nil).Representation()

	return CreateCIDFromNode(d).String()

}

// See:
func CreateDescriptorCID(descriptor Descriptor) string {

	node := bindnode.Wrap(&descriptor, DescriptorSchemaType).Representation()
	if node.IsNull() {
		return ""
	}
	return CreateCIDFromNode(node).String()

}

func CreateProcessingCID(mp MessageProcessing) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "nonce", qp.String(mp.Nonce))
		qp.MapEntry(ma, "author", qp.String(mp.AuthorDID))
		qp.MapEntry(ma, "recipient", qp.String(mp.RecipientDID))
	})
	if err != nil {
		return ""
	}

	return CreateCIDFromNode(d).String()

}

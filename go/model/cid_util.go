package model

import (
	"bytes"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

func CreateRecordCID(descriptorCID string, processingCID string) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "descriptorCid", qp.String(descriptorCID))
		qp.MapEntry(ma, "processingCid", qp.String(processingCID))
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

func CreateDescriptorCID(descriptor Descriptor) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "method", qp.String(descriptor.Method))
		qp.MapEntry(ma, "schema", qp.String(descriptor.Schema))
		qp.MapEntry(ma, "dataCid", qp.String(descriptor.DataCID))
		qp.MapEntry(ma, "nonce", qp.String(descriptor.Nonce))
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

func CreateProcessingCID(mp MessageProcessing) string {

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

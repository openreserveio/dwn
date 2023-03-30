package model

//Given a type model.Descriptor object:
//
//type Descriptor struct {
//
//	// Base Required Fields per https://identity.foundation/decentralized-web-node/spec/#messages
//	Method     string `json:"method" bson:"method"`
//	DataCID    string `json:"dataCid,omitempty" bson:"dataCID"`
//	DataFormat string `json:"dataFormat,omitempty" bson:"dataFormat"`
//
//	// CollectionsQuery, HooksWrite per https://identity.foundation/decentralized-web-node/spec/#collectionsquery, https://identity.foundation/decentralized-web-node/spec/#hooks
//	Filter DescriptorFilter `json:"filter,omitempty"`
//
//	// CollectionsWrite, Delete, Commit per https://identity.foundation/decentralized-web-node/spec/#collectionswrite
//	ParentID        string     `json:"parentId,omitempty" bson:"parent_id"`
//	Protocol        string     `json:"protocol,omitempty" bson:"protocol"`
//	ProtocolVersion string     `json:"protocolVersion,omitempty" bson:"protocol_version"`
//	Schema          string     `json:"schema,omitempty" bson:"schema"`
//	CommitStrategy  string     `json:"commitStrategy,omitempty" bson:"commit_strategy"`
//	Published       bool       `json:"published,omitempty" bson:"published"`
//	DateCreated     time.Time  `json:"dateCreated,omitempty" bson:"date_created"`
//	DatePublished   *time.Time `json:"datePublished,omitempty" bson:"date_published"`
//
//	// HooksWrite per https://identity.foundation/decentralized-web-node/spec/#hooks
//	URI string `json:"uri,omitempty" bson:"uri"`
//}
//
//and a DescriptorFilter object:
//
//type DescriptorFilter struct {
//	Schema          string `json:"schema,omitempty"`
//	RecordID        string `json:"recordId,omitempty"`
//	ContextID       string `json:"contextId,omitempty"`
//	Protocol        string `json:"protocol,omitempty"`
//	ProtocolVersion string `json:"protocolVersion,omitempty"`
//	DataFormat      string `json:"dataFormat,omitempty"`
//	DateSort        string `json:"dateSort,omitempty"`
//}
//
//Write a public golang function that:
//- takes a Descriptor object reference as input
//- creates a DAG CBOR encoded byte array from the input
//- returns the Version 1 CID of the CBOR encoded byte array as a cid.Cid object

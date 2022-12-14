package model

import "time"

/**
Specification: https://identity.foundation/decentralized-web-node/spec/#request-objects
*/

type RequestObject struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	RecordID      string            `json:"recordId,omitempty" bson:"record_id"`
	ContextID     string            `json:"contextId,omitempty" bson:"context_id"`
	Data          string            `json:"data,omitempty" bson:"data"`
	Processing    MessageProcessing `json:"processing" bson:"processing"`
	Descriptor    Descriptor        `json:"descriptor" bson:"descriptor"`
	Attestation   DWNJWS            `json:"attestation,omitempty" bson:"attestation"`
	Authorization DWNJWS            `json:"authorization,omitempty" bson:"authorization"`
}

type Descriptor struct {
	// Base Required Fields per https://identity.foundation/decentralized-web-node/spec/#messages
	Method     string `json:"method" bson:"method"`
	DataCID    string `json:"dataCid,omitempty" bson:"dataCID"`
	DataFormat string `json:"dataFormat,omitempty" bson:"dataFormat"`

	// CollectionsQuery per https://identity.foundation/decentralized-web-node/spec/#collectionsquery
	Filter CollectionsQueryFilter `json:"filter,omitempty"`

	// CollectionsWrite, Delete, Commit per https://identity.foundation/decentralized-web-node/spec/#collectionswrite
	ParentID        string    `json:"parentId,omitempty" bson:"parent_id"`
	Protocol        string    `json:"protocol,omitempty" bson:"protocol"`
	ProtocolVersion string    `json:"protocolVersion,omitempty" bson:"protocol_version"`
	Schema          string    `json:"schema,omitempty" bson:"schema"`
	CommitStrategy  string    `json:"commitStrategy,omitempty" bson:"commit_strategy"`
	Published       bool      `json:"published,omitempty" bson:"published"`
	DateCreated     time.Time `json:"dateCreated,omitempty" bson:"date_created"`
	DatePublished   time.Time `json:"datePublished,omitempty" bson:"date_published"`
}

type MessageProcessing struct {
	Nonce        string `json:"nonce" bson:"nonce"`
	AuthorDID    string `json:"author" bson:"author_did"`
	RecipientDID string `json:"recipient" bson:"recipient_did"`
}

type DWNJWS struct {
	Payload    string      `json:"payload" bson:"payload"`
	Signatures []DWNJWSSig `json:"signatures" bson:"signatures"`
}

type DWNJWSSig struct {
	Protected string `json:"protected" bson:"protected"`
	Signature string `json:"signature" bson:"signature"`
}

type CollectionsQueryFilter struct {
	Schema          string `json:"schema,omitempty"`
	RecordID        string `json:"recordId,omitempty"`
	ContextID       string `json:"contextId,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	DataFormat      string `json:"dataFormat,omitempty"`
	DateSort        string `json:"dateSort,omitempty"`
}

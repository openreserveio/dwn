package model

import "time"

/**
Specification: https://identity.foundation/decentralized-web-node/spec/#request-objects
*/

type RequestObject struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	RecordID      string            `json:"recordId,omitempty"`
	ContextID     string            `json:"contextId,omitempty"`
	Data          string            `json:"data,omitempty"`
	Processing    MessageProcessing `json:"processing"`
	Descriptor    Descriptor        `json:"descriptor"`
	Attestation   DWNJWS            `json:"attestation,omitempty"`
	Authorization DWNJWS            `json:"authorization,omitempty"`
}

type Descriptor struct {
	// Base Required Fields per https://identity.foundation/decentralized-web-node/spec/#messages
	Method     string `json:"method"`
	DataCID    string `json:"dataCid,omitempty"`
	DataFormat string `json:"dataFormat,omitempty"`

	// CollectionsQuery per https://identity.foundation/decentralized-web-node/spec/#collectionsquery
	Filter CollectionsQueryFilter `json:"filter,omitempty"`

	// CollectionsWrite, Delete, Commit per https://identity.foundation/decentralized-web-node/spec/#collectionswrite
	ParentID        string    `json:"parentId,omitempty"`
	Protocol        string    `json:"protocol,omitempty"`
	ProtocolVersion string    `json:"protocolVersion,omitempty"`
	Schema          string    `json:"schema,omitempty"`
	CommitStrategy  string    `json:"commitStrategy,omitempty"`
	Published       bool      `json:"published,omitempty"`
	DateCreated     time.Time `json:"dateCreated,omitempty"`
	DatePublished   time.Time `json:"datePublished,omitempty"`
}

type MessageProcessing struct {
	Nonce        string `json:"nonce"`
	AuthorDID    string `json:"author"`
	RecipientDID string `json:"recipient"`
}

type DWNJWS struct {
	Payload    string      `json:"payload"`
	Signatures []DWNJWSSig `json:"signatures"`
}

type DWNJWSSig struct {
	Protected string `json:"protected"`
	Signature string `json:"signature"`
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

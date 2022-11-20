package model

/**
Specification: https://identity.foundation/decentralized-web-node/spec/#request-objects
*/

type RequestObject struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	RecordID      string            `json:"recordId"`
	Data          string            `json:"data,omitempty"`
	Processing    MessageProcessing `json:"processing"`
	Descriptor    Descriptor        `json:"descriptor"`
	Attestation   interface{}       `json:"attestation,omitempty"`
	Authorization interface{}       `json:"authorization,omitempty"`
}

type Descriptor struct {
	Nonce      string `json:"nonce"`
	Method     string `json:"method"`
	DataCID    string `json:"dataCid"`
	DataFormat string `json:"dataFormat"`
	Schema     string `json:"schema"`
}

type MessageProcessing struct {
	Nonce        string `json:"nonce"`
	AuthorDID    string `json:"author"`
	RecipientDID string `json:"recipient"`
}

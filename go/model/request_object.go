package model

/**
Specification: https://identity.foundation/decentralized-web-node/spec/#request-objects
*/

type RequestObject struct {
	TargetDID string    `json:"target"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	RecordID      string            `json:"recordId"`
	Processing    MessageProcessing `json:"processing"`
	Data          string            `json:"data,omitempty"`
	Descriptor    Descriptor        `json:"descriptor"`
	Attestation   interface{}       `json:"attestation,omitempty"`
	Authorization interface{}       `json:"authorization,omitempty"`
}

type Descriptor struct {
	Nonce      string `json:"nonce"`
	Method     string `json:"method"`
	DataCID    string `json:"dataCid"`
	DataFormat string `json:"dataFormat"`
}

type MessageProcessing struct {
	Nonce        string `json:"nonce"`
	AuthorDID    string `json:"author"`
	RecipientDID string `json:"recipient"`
}

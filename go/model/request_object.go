package model

/**
Specification: https://identity.foundation/decentralized-web-node/spec/#request-objects
*/

const (
	DATA_FORMAT_JSON   = "application/json"
	DATA_FORMAT_VC_JWT = "application/vc+jwt"
	DATA_FORMAT_VC_LDP = "application/vc+ldp"
)

type RequestObject struct {
	TargetDID string    `json:"target"`
	Messages  []Message `json:"messages"`
}

type Message struct {
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
	TargetDID    string `json:"target"`
	RecipientDID string `json:"recipient"`
}

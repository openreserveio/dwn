package model

type FeatureDetection struct {
	Type       string           `json:"type"`
	Interfaces FeatureInterface `json:"interfaces"`
}

type FeatureInterface struct {
	Records     RecordsFeatures     `json:"records,omitempty"`
	Hooks       HooksFeatures       `json:"hooks,omitempty"`
	Permissions PermissionsFeatures `json:"permissions,omitempty"`
	Messaging   MessagingFeatures   `json:"messaging,omitempty"`
	Protocols   ProtocolsFeatures   `json:"protocols,omitempty"`
}

type RecordsFeatures struct {
	RecordsQuery  bool `json:"RecordsQuery"`
	RecordsWrite  bool `json:"RecordsWrite"`
	RecordsCommit bool `json:"RecordsCommit"`
	RecordsDelete bool `json:"RecordsDelete"`
}

type HooksFeatures struct {
	HooksQuery  bool `json:"HooksQuery"`
	HooksWrite  bool `json:"HooksWrite"`
	HooksDelete bool `json:"HooksDelete"`
}

type PermissionsFeatures struct {
	PermissionsRequest bool `json:"PermissionsRequest"`
	PermissionsGrant   bool `json:"PermissionsGrant"`
	PermissionsRevoke  bool `json:"PermissionsRevoke"`
}

type MessagingFeatures struct {
	Batching bool `json:"batching"`
}

type ProtocolsFeatures struct {
	ProtocolsQuery     bool `json:"ProtocolsQuery"`
	ProtocolsConfigure bool `json:"ProtocolsConfigure"`
}

var CurrentFeatureDetection FeatureDetection = FeatureDetection{
	Type: "FeatureDetection",
	Interfaces: FeatureInterface{
		Records: RecordsFeatures{
			RecordsQuery:  true,
			RecordsWrite:  true,
			RecordsCommit: true,
			RecordsDelete: true,
		},
		Hooks: HooksFeatures{
			HooksQuery:  true,
			HooksWrite:  true,
			HooksDelete: true,
		},
		Permissions: PermissionsFeatures{
			PermissionsRequest: true,
			PermissionsGrant:   true,
			PermissionsRevoke:  true,
		},
		Messaging: MessagingFeatures{
			Batching: true,
		},
		Protocols: ProtocolsFeatures{
			ProtocolsQuery:     true,
			ProtocolsConfigure: true,
		},
	},
}

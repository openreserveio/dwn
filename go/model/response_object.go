package model

type ResponseStatus struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

type ResponseObject struct {
	Status  ResponseStatus        `json:"status"`
	Replies []MessageResultObject `json:"replies"`
}

type MessageResultObject struct {
	Status  ResponseStatus       `json:"status"`
	Entries []MessageResultEntry `json:"entries"`
}

type MessageResultEntry struct {
	Result []byte `json:"result"`
}

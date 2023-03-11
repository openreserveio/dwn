package client

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/openreserveio/dwn/go/model"
)

func (client *DWNClient) CallDWNHTTP(ro model.RequestObject) (model.ResponseObject, error) {

	respObj := model.ResponseObject{}

	res, err := resty.New().R().
		SetBody(ro).
		SetHeader(HEADER_CONTENT_TYPE_KEY, HEADER_CONTENT_TYPE_APPLICATION_JSON).
		Post(client.DWNUrlBase)

	if err != nil {
		return respObj, err
	}
	if !res.IsSuccess() {
		return respObj, errors.New("Unable to create data")
	}

	err = json.Unmarshal(res.Body(), &respObj)

	return respObj, nil

}

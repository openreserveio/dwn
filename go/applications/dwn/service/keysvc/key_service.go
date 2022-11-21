package keysvc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openreserveio/dwn/go/applications/dwn/service/keysvc/messages"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
)

type KeyService struct {
	services.UnimplementedKeyServiceServer
}

func CreateKeyService() (KeyService, error) {
	return KeyService{}, nil
}

func (k KeyService) VerifyMessageAttestation(ctx context.Context, request *services.VerifyMessageAttestationRequest) (*services.VerifyMessageAttestationResponse, error) {

	// reassemble message
	var message model.Message
	err := json.Unmarshal(request.Message, &message)
	if err != nil {
		return k.invalidAttestationResponse(err), nil
	}

	// Ensure attestation section
	if message.Attestation.Payload == "" {
		return k.invalidAttestationResponse(errors.New("No attestation present")), nil
	}

	if !messages.VerifyMessageAttestation(message) {
		return k.invalidAttestationResponse(errors.New("Invalid Attestation")), nil
	}

	return k.validAttestationResponse(), nil

}

func (k KeyService) validAttestationResponse() *services.VerifyMessageAttestationResponse {

	return &services.VerifyMessageAttestationResponse{
		Status: &services.CommonStatus{
			Status: services.Status_OK,
		},
	}

}

func (k KeyService) invalidAttestationResponse(err error) *services.VerifyMessageAttestationResponse {

	return &services.VerifyMessageAttestationResponse{
		Status: &services.CommonStatus{
			Status:  services.Status_INVALID_ATTESTATION,
			Details: err.Error(),
		},
	}

}

func (k KeyService) mustEmbedUnimplementedKeyServiceServer() {
	//TODO implement me
	panic("implement me")
}

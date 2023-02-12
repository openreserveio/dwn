package client

type DWNClient struct {
	DWNUrlBase      string
	Protocol        string
	ProtocolVersion string
}

func CreateDWNClient(urlBase string) *DWNClient {

	return &DWNClient{
		DWNUrlBase: urlBase,
	}

}

func CreateDWNClientForProtocol(urlBase string, protocol string, protocolVersion string) *DWNClient {

	return &DWNClient{
		DWNUrlBase:      urlBase,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
	}

}

func (client *DWNClient) SaveData(schemaUrl string, data []byte, dataAuthor *Identity, dataRecipient *Identity) error {

	return nil
}

func (client *DWNClient) UpdateData(schemaUrl string, primaryIdentifier string, data []byte, dataOwner *Identity) error {

	return nil
}

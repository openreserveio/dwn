package client

const (
	HEADER_CONTENT_TYPE_KEY              = "Content-Type"
	HEADER_CONTENT_TYPE_APPLICATION_JSON = "application/json"
)

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

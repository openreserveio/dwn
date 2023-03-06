package fileenv

import "github.com/spf13/viper"

const (
	ENV_PREFIX = "DWN"

	KEY_API_LISTEN_ADDRESS = "api.listenAddress"
	KEY_API_LISTEN_PORT    = "api.listenPort"

	KEY_RECORDSVC_LISTEN_ADDRESS       = "recordsvc.listenAddress"
	KEY_RECORDSVC_LISTEN_PORT          = "recordsvc.listenPort"
	KEY_RECORDSVC_EXTERNAL_ADDRESS     = "recordsvc.externalAddress"
	KEY_RECORDSVC_EXTERNAL_PORT        = "recordsvc.externalPort"
	KEY_RECORDSVC_DOCDB_CONNECTION_URI = "recordsvc.docdbConnectionURI"

	KEY_HOOKSVC_LISTEN_ADDRESS       = "hooksvc.listenAddress"
	KEY_HOOKSVC_LISTEN_PORT          = "hooksvc.listenPort"
	KEY_HOOKSVC_EXTERNAL_ADDRESS     = "hooksvc.externalAddress"
	KEY_HOOKSVC_EXTERNAL_PORT        = "hooksvc.externalPort"
	KEY_HOOKSVC_DOCDB_CONNECTION_URI = "hooksvc.docdbConnectionURI"

	KEY_QUEUE_SERVICE_CONNECTION_URI = "queueservice.connectionURI"
)

type FileEnvironmentConfiguration struct {
}

func CreateFileEnvironmentConfiguration() *FileEnvironmentConfiguration {

	viper.SetConfigName("dwn")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dwn")
	viper.AddConfigPath("/opt/dwn")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	return &FileEnvironmentConfiguration{}

}

func (config FileEnvironmentConfiguration) GetAPIListenAddress() string {
	return viper.GetString(KEY_API_LISTEN_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetAPIListenPort() int {
	return viper.GetInt(KEY_API_LISTEN_PORT)
}

func (config FileEnvironmentConfiguration) GetRecordServiceListenAddress() string {
	return viper.GetString(KEY_RECORDSVC_LISTEN_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetRecordServiceListenPort() int {
	return viper.GetInt(KEY_RECORDSVC_LISTEN_PORT)
}

func (config FileEnvironmentConfiguration) GetRecordServiceDocumentDBURI() string {
	return viper.GetString(KEY_RECORDSVC_DOCDB_CONNECTION_URI)
}

func (config FileEnvironmentConfiguration) GetRecordServiceExternalAddress() string {
	return viper.GetString(KEY_RECORDSVC_EXTERNAL_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetRecordServiceExternalPort() int {
	return viper.GetInt(KEY_RECORDSVC_EXTERNAL_PORT)
}

func (config FileEnvironmentConfiguration) GetHookServiceListenAddress() string {
	return viper.GetString(KEY_HOOKSVC_LISTEN_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetHookServiceListenPort() int {
	return viper.GetInt(KEY_HOOKSVC_LISTEN_PORT)
}

func (config FileEnvironmentConfiguration) GetHookServiceDocumentDBURI() string {
	return viper.GetString(KEY_HOOKSVC_DOCDB_CONNECTION_URI)
}

func (config FileEnvironmentConfiguration) GetHookServiceExternalAddress() string {
	return viper.GetString(KEY_HOOKSVC_EXTERNAL_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetHookServiceExternalPort() int {
	return viper.GetInt(KEY_HOOKSVC_EXTERNAL_PORT)
}

func (config FileEnvironmentConfiguration) GetQueueServiceConnectionURI() string {
	return viper.GetString(KEY_QUEUE_SERVICE_CONNECTION_URI)
}

package fileenv

import "github.com/spf13/viper"

const (
	ENV_PREFIX = "DWN"

	KEY_API_LISTEN_ADDRESS = "api.listenAddress"
	KEY_API_LISTEN_PORT    = "api.listenPort"

	KEY_COLLSVC_LISTEN_ADDRESS       = "collsvc.listenAddress"
	KEY_COLLSVC_LISTEN_PORT          = "collsvc.listenPort"
	KEY_COLLSVC_DOCDB_CONNECTION_URI = "collsvc.docdbConnectionURI"
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

func (config FileEnvironmentConfiguration) GetCollectionServiceListenAddress() string {
	return viper.GetString(KEY_COLLSVC_LISTEN_ADDRESS)
}

func (config FileEnvironmentConfiguration) GetCollectionServiceListenPort() int {
	return viper.GetInt(KEY_COLLSVC_LISTEN_PORT)
}

func (config FileEnvironmentConfiguration) GetCollectionServiceDocumentDBURI() string {
	return viper.GetString(KEY_COLLSVC_DOCDB_CONNECTION_URI)
}

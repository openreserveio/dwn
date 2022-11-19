package fileenv

import "github.com/spf13/viper"

const (
	ENV_PREFIX = "DWN"

	KEY_API_LISTEN_ADDRESS = "api.listenAddress"
	KEY_API_LISTEN_PORT    = "api.listenPort"

	KEY_DOCDB_CONNECTION_URI = "docdb.connection_uri"
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

func (config FileEnvironmentConfiguration) GetDocumentDBCollectionURI() string {
	return viper.GetString(KEY_DOCDB_CONNECTION_URI)
}

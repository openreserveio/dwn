package configuration

import (
	"github.com/openreserveio/dwn/go/applications/dwn/configuration/fileenv"
	"github.com/openreserveio/dwn/go/log"
	"os"
)

const (
	CONFIG_TYPE_FILEENV = "fileenv"
	CONFIG_TYPE_ETCD    = "etcd"
)

type Configuration interface {

	// API Configuration
	GetAPIListenAddress() string
	GetAPIListenPort() int

	// Record Service
	GetRecordServiceListenAddress() string
	GetRecordServiceListenPort() int
	GetRecordServiceExternalAddress() string
	GetRecordServiceExternalPort() int
	GetRecordServicePostgresURI() string

	// Hook Service
	GetHookServiceListenAddress() string
	GetHookServiceListenPort() int
	GetHookServiceExternalAddress() string
	GetHookServiceExternalPort() int
	GetHookServicePostgresURI() string

	// Queue Service
	GetQueueServiceConnectionURI() string

	// Queue Names
	GetNotifyCallbackQueueName() string
}

func Config() (Configuration, error) {

	configTypeFlag := os.Getenv("OR_CONFIG")

	var config Configuration
	switch configTypeFlag {

	case CONFIG_TYPE_FILEENV:
		log.Info("Configuration Type:  File/Environment")
		config = fileenv.CreateFileEnvironmentConfiguration()

	case CONFIG_TYPE_ETCD:
		log.Info("Configuration Type:  Distributed/ETCD")
		log.Fatal("Not yet supported")
		os.Exit(1)

	default:
		log.Info("Configuration Type:  File/Environment")
		config = fileenv.CreateFileEnvironmentConfiguration()

	}

	return config, nil

}

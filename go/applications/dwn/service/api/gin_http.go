package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/log"
)

type APIService struct {
	ListenAddress string
	ListenPort    int
	Gin           *gin.Engine
}

func CreateAPIService(options *framework.ServiceOptions) (*APIService, error) {

	apiService := APIService{
		ListenAddress: options.Address,
		ListenPort:    options.Port,
		Gin:           gin.Default(),
	}

	apiService.Gin.POST("/", apiService.HandleDWNRequest)

	return &apiService, nil

}

func (apiService APIService) Run() error {
	return apiService.Gin.Run(fmt.Sprintf("%s:%d", apiService.ListenAddress, apiService.ListenPort))
}

func (apiService APIService) HandleDWNRequest(ctx *gin.Context) {
	log.Info("Indeed")
}

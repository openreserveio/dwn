package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
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

	apiService.Gin.GET("/", apiService.HandleFeatureRequest)
	apiService.Gin.POST("/", apiService.HandleDWNRequest)

	return &apiService, nil

}

func (apiService APIService) Run() error {
	return apiService.Gin.Run(fmt.Sprintf("%s:%d", apiService.ListenAddress, apiService.ListenPort))
}

func (apiService APIService) HandleDWNRequest(ctx *gin.Context) {

	ro, err := apiService.GetRequestObject(ctx)
	if err != nil {
		log.Error("Error while parsing request object:  %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response := model.ResponseObject{
		Status: model.ResponseStatus{Code: 200, Detail: fmt.Sprintf("TargetDID:  %s", ro.TargetDID)},
	}
	ctx.JSON(200, &response)

}

func (apiService APIService) GetRequestObject(ctx *gin.Context) (*model.RequestObject, error) {

	var request model.RequestObject
	err := ctx.BindJSON(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil

}

func (apiService APIService) HandleFeatureRequest(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, model.CurrentFeatureDetection)
	return

}

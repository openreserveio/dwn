package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"google.golang.org/grpc"
	"net/http"
)

type APIService struct {
	ListenAddress string
	ListenPort    int
	Gin           *gin.Engine
	Router        *FeatureRouter
	CollSvcClient *services.CollectionServiceClient
}

func CreateAPIService(apiServiceOptions *framework.ServiceOptions, collSvcOptions *framework.ServiceOptions) (*APIService, error) {

	var err error
	var clientConn *grpc.ClientConn
	if collSvcOptions.SecureFlag {
		log.Fatal("Secure GRPC not yet supported - use Istio")
		return nil, errors.New("Secure GRPC not yet supported - use Istio")
	} else {
		clientConn, err = grpc.Dial(fmt.Sprintf("%s:%d", collSvcOptions.Address, collSvcOptions.Port))
		if err != nil {
			return nil, err
		}
	}

	collSvcClient := services.NewCollectionServiceClient(clientConn)
	fr, err := CreateFeatureRouter(collSvcClient, 15)
	apiService := APIService{
		ListenAddress: apiServiceOptions.Address,
		ListenPort:    apiServiceOptions.Port,
		Gin:           gin.Default(),
		CollSvcClient: &collSvcClient,
		Router:        fr,
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

	responseObject, err := apiService.Router.Route(ro)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, &responseObject)

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

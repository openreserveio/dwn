package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
		clientConn, err = grpc.Dial(fmt.Sprintf("%s:%d", collSvcOptions.Address, collSvcOptions.Port), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	}

	collSvcClient := services.NewCollectionServiceClient(clientConn)
	fr, err := CreateFeatureRouter(collSvcClient, 15)

	// Configure Tracing
	ginEngine := gin.Default()
	ginEngine.Use(otelgin.Middleware("DWN_API_SERVICE"))

	apiService := APIService{
		ListenAddress: apiServiceOptions.Address,
		ListenPort:    apiServiceOptions.Port,
		Gin:           ginEngine,
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

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "DWNRequest")
	defer childSpan.End()

	childSpan.AddEvent("Parsing Request Object")
	ro, err := apiService.GetRequestObject(ctx)
	if err != nil {
		log.Error("Error while parsing request object:  %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	childSpan.AddEvent("Request Object Parsed")

	childSpan.AddEvent("Routing Request")
	responseObject, err := apiService.Router.Route(ctx, ro)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	childSpan.AddEvent("Request routed with response")

	ctx.JSON(200, &responseObject)

}

func (apiService APIService) GetRequestObject(ctx *gin.Context) (*model.RequestObject, error) {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "GetRequestObject")
	defer childSpan.End()

	var request model.RequestObject
	err := ctx.BindJSON(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil

}

func (apiService APIService) HandleFeatureRequest(ctx *gin.Context) {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "HandleFeatureRequest")
	defer childSpan.End()

	childSpan.AddEvent("Current Feature Detection!")
	ctx.JSON(http.StatusOK, model.CurrentFeatureDetection)

	return

}

package presentation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	pubsubmessaging "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/rest"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/serverutils"
	"google.golang.org/api/healthcare/v1"
)

// ClinicalAllowedOrigins is a list of CORS origins allowed to interact with
// this service
var ClinicalAllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.co.ke",
	"http://localhost:5000",
	"http://localhost:4200",
	"https://clinical-staging.healthcloud.co.ke",
	"https://clinical-testing.healthcloud.co.ke",
	"https://clinical-prod.healthcloud.co.ke",
	"https://clinical-staging.bewell.co.ke",
	"https://clinical-testing.bewell.co.ke",
	"https://clinical-demo.bewell.co.ke",
	"https://clinical-prod.bewell.co.ke",
	"https://uat-emr.advantage.slade360.com",
	"https://advantage.slade360.com",
}

var (
	authServerEndpoint = serverutils.MustGetEnvVar("AUTHSERVER_ENDPOINT")
	clientID           = serverutils.MustGetEnvVar("CLIENT_ID")
	clientSecret       = serverutils.MustGetEnvVar("CLIENT_SECRET")
	username           = serverutils.MustGetEnvVar("AUTH_USERNAME")
	password           = serverutils.MustGetEnvVar("AUTH_PASSWORD")
	grantType          = serverutils.MustGetEnvVar("GRANT_TYPE")
)

// StartServer sets up gin
func StartServer(
	ctx context.Context,
	port int,
) {
	err := serverutils.Sentry()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	baseExtension := extensions.NewBaseExtensionImpl()

	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)

	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("unable to initialize pubsub client: %w", err))
	}

	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	datasetID := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	datasetLocation := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_LOCATION")
	fhirStoreID := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")

	hsv, err := healthcare.NewService(ctx)
	if err != nil {
		log.Panicf("unable to initialize new Google Cloud Healthcare Service: %s", err)
	}

	repo := fhirdataset.NewFHIRRepository(ctx, hsv, project, datasetID, datasetLocation, fhirStoreID)
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	upload := upload.NewServiceUpload(ctx)

	pubsubSvc, err := pubsubmessaging.NewServicePubSubMessaging(ctx, pubSubClient)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("failed to initialize pubsub messaging service: %w", err))
	}

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl, upload, pubsubSvc)

	authServerConfig := authutils.Config{
		AuthServerEndpoint: authServerEndpoint,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		GrantType:          grantType,
		Username:           username,
		Password:           password,
	}

	authclient, err := authutils.NewClient(authServerConfig)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	usecases := usecases.NewUsecasesInteractor(infrastructure)

	r := gin.Default()

	memoryStore := persist.NewMemoryStore(60 * time.Minute)

	SetupRoutes(r, memoryStore, authclient, usecases, infrastructure)

	addr := fmt.Sprintf(":%d", port)

	if err := r.Run(addr); err != nil {
		serverutils.LogStartupError(ctx, err)
	}
}

func SetupRoutes(r *gin.Engine, cacheStore persist.CacheStore, authclient *authutils.Client, usecases usecases.Interactor, infra infrastructure.Infrastructure) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: ClinicalAllowedOrigins,
		AllowMethods: []string{http.MethodPut, http.MethodPatch, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			"Accept",
			"Accept-Charset",
			"Accept-Language",
			"Accept-Encoding",
			"Origin",
			"Host",
			"User-Agent",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Authorization",
			"Clinical-Organization-ID",
			"Clinical-Facility-ID",
		},
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8000" || origin == "http://localhost:8080"
		},
		MaxAge:          12 * time.Hour,
		AllowWebSockets: true,
	}))

	handlers := rest.NewPresentationHandlers(usecases, infra.BaseExtension)

	graphQL := r.Group("/graphql")
	graphQL.Use(rest.AuthenticationGinMiddleware(cacheStore, *authclient))
	graphQL.Use(rest.TenantIdentifierExtractionMiddleware(infra.FHIR))
	graphQL.Any("", GQLHandler(usecases))

	// Unauthenticated routes
	ide := r.Group("/ide")
	ide.Any("", playgroundHandler())

	r.POST("/pubsub", handlers.ReceivePubSubPushMessage)

	apis := r.Group("/api")
	apis.Use(rest.AuthenticationGinMiddleware(cacheStore, *authclient))

	v1 := apis.Group("/v1")

	tenants := v1.Group("/tenants")
	tenants.POST("", handlers.RegisterTenant)

	facilities := v1.Group("/facilities")
	facilities.POST("", handlers.RegisterFacility)

	upload := v1.Group("/media")
	upload.Use(rest.AuthenticationGinMiddleware(cacheStore, *authclient))
	upload.Use(rest.TenantIdentifierExtractionMiddleware(infra.FHIR))
	upload.POST("", handlers.UploadMedia)
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(service usecases.Interactor) gin.HandlerFunc {
	resolver, err := graph.NewResolver(service)
	if err != nil {
		log.Panicf("failed to start graph resolver: %s", err)
	}

	server := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	return func(ctx *gin.Context) {
		server.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL IDE", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

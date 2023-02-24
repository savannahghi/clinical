package presentation

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	pubsubmessaging "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/rest"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
)

const serverTimeoutSeconds = 120

// ClinicalAllowedOrigins is a list of CORS origins allowed to interact with
// this service
var ClinicalAllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.co.ke",
	"http://localhost:5000",
	"https://clinical-staging.healthcloud.co.ke",
	"https://clinical-testing.healthcloud.co.ke",
	"https://clinical-prod.healthcloud.co.ke",
	"https://clinical-staging.bewell.co.ke",
	"https://clinical-testing.bewell.co.ke",
	"https://clinical-demo.bewell.co.ke",
	"https://clinical-prod.bewell.co.ke",
}

// ClinicalAllowedHeaders is a list of CORS allowed headers for the clinical
// service
var ClinicalAllowedHeaders = []string{
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
}

var (
	authServerEndpoint = serverutils.MustGetEnvVar("AUTHSERVER_ENDPOINT")
	clientID           = serverutils.MustGetEnvVar("CLIENT_ID")
	clientSecret       = serverutils.MustGetEnvVar("CLIENT_SECRET")
	username           = serverutils.MustGetEnvVar("AUTH_USERNAME")
	password           = serverutils.MustGetEnvVar("AUTH_PASSWORD")
	grantType          = serverutils.MustGetEnvVar("GRANT_TYPE")
)

// PrepareServer sets up the HTTP server
func PrepareServer(
	ctx context.Context,
	port int,
	allowedOrigins []string,
) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedHeaders(ClinicalAllowedHeaders),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(h, "application/json")
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
	log.Infof("Server running at port %v", addr)
	return srv
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	baseExtension := extensions.NewBaseExtensionImpl()

	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl)
	usecases := usecases.NewUsecasesInteractor(infrastructure)
	h := rest.NewPresentationHandlers(usecases)

	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(pubSubClient, baseExtension, infrastructure, usecases)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pubsub messaging service: %v", err)
	}

	authServerConfig := authutils.Config{
		AuthServerEndpoint: authServerEndpoint,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		GrantType:          grantType,
		Username:           username,
		Password:           password,
	}
	authClient, err := authutils.NewClient(authServerConfig)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))

	r.Path("/pubsub").Methods(
		http.MethodPost,
	).HandlerFunc(pubSub.ReceivePubSubPushMessages)

	// check server status.
	r.Path("/health").HandlerFunc(serverutils.HealthStatusCheck)

	r.Path("/delete_patient").
		Methods(http.MethodPost).
		HandlerFunc(h.DeleteFHIRPatientByPhone())

	// ISC routes. These are inter service route
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(authutils.SladeAuthenticationMiddleware(*authClient))
	isc.Path("/delete-patient").Methods(
		http.MethodDelete,
	).HandlerFunc(h.DeleteFHIRPatientByPhone())
	isc.Path("/create-fhir-organisation").Methods(
		http.MethodPost,
	).HandlerFunc(h.CreateFHIROrganization())

	//Authenticated routes
	gqlR := r.Path("/graphql").Subrouter()
	gqlR.Use(authutils.SladeAuthenticationMiddleware(*authClient))
	gqlR.Use(rest.TenantIdentifierExtractionMiddleware(usecases))
	gqlR.Methods(
		http.MethodPost, http.MethodGet, http.MethodOptions,
	).HandlerFunc(GQLHandler(ctx, usecases))
	return r, nil
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	service usecases.Interactor,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, service)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	server := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}
}

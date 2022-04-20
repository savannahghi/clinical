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
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	pubsubmessaging "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/rest"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
	"github.com/savannahghi/firebasetools"
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
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	baseExtension := extensions.NewBaseExtensionImpl(fc)

	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension)
	usecases := usecases.NewUsecasesInteractor(infrastructure)
	oclUseCase := ocl.NewUseCasesImpl(infrastructure)
	h := rest.NewPresentationHandlers(infrastructure, usecases)

	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(pubSubClient, baseExtension, infrastructure, usecases, usecases, oclUseCase)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pubsub messaging service: %v", err)
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

	//Authenticated routes
	gqlR := r.Path("/graphql").Subrouter()
	gqlR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
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

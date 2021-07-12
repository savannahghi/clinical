package graph

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
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
	return srv
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
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

	// check server status.
	r.Path("/health").HandlerFunc(serverutils.HealthStatusCheck)

	r.Path("/delete_patient").
		Methods(http.MethodPost).
		HandlerFunc(DeleteFHIRPatientByPhone(ctx))

	// Authenticated routes
	gqlR := r.Path("/graphql").Subrouter()
	gqlR.Use(base.AuthenticationMiddleware(firebaseApp))
	gqlR.Methods(
		http.MethodPost, http.MethodGet, http.MethodOptions,
	).HandlerFunc(graphqlHandler())
	return r, nil

}

func graphqlHandler() http.HandlerFunc {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: NewResolver(),
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}

// DeleteFHIRPatientByPhone handler exposes an endpoint that takes a
// patient's phone number and deletes the patient's FHIR compartment
func DeleteFHIRPatientByPhone(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := clinical.NewService()

		payload := &clinical.PhoneNumberPayload{}
		type errResponse struct {
			Err string `json:"error"`
		}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == "" {
			serverutils.WriteJSONResponse(
				w,
				errResponse{
					Err: "expected a phone number to be defined",
				},
				http.StatusBadRequest,
			)
			return
		}
		deleted, err := s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
		if err != nil {
			err := fmt.Sprintf("unable to delete patient: %v", err.Error())
			serverutils.WriteJSONResponse(
				w,
				errResponse{
					Err: err,
				},
				http.StatusInternalServerError,
			)
			return
		}

		type response struct {
			Deleted bool `json:"deleted"`
		}
		serverutils.WriteJSONResponse(
			w,
			response{Deleted: deleted},
			http.StatusOK,
		)
	}
}

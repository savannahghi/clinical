package main

import (
	"compress/gzip"
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/clinical/graph"
	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
)

const serverTimeoutSeconds = 120

var allowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"https://api-gateway-test.healthcloud.co.ke",
	"https://api-gateway-prod.healthcloud.co.ke",
	"https://clinical-testing-uyajqt434q-ew.a.run.app",
	"https://clinical-prod-uyajqt434q-ew.a.run.app",
}
var allowedHeaders = []string{
	"clinical", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}

func main() {
	ctx := context.Background()

	err := base.Sentry()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start up the router
	r, err := Router()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start the server
	addr := ":" + base.MustGetEnvVar("PORT")
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedHeaders(allowedHeaders),
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
	log.Fatal(srv.ListenAndServe())
}

// Router sets up the ginContext router
func Router() (*mux.Router, error) {
	ctx := context.Background()
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}
	clinicalService := clinical.NewService()
	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(base.RequestDebugMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/profiles/{id}").Methods(
		http.MethodGet, http.MethodOptions).HandlerFunc(
		clinical.PatientProfileHandlerFunc(clinicalService))
	r.Path("/visits/{id}").Methods(
		http.MethodGet, http.MethodOptions).HandlerFunc(
		clinical.LastVisitHandlerFunc(ctx, clinicalService))
	r.Path("/charts/{id}").Methods(
		http.MethodGet, http.MethodOptions).HandlerFunc(
		clinical.FullHistoryHandlerFunc(ctx, clinicalService))
	r.Path("/base.css").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(base.CSS())
	r.Path("/visit.css").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(base.VisitCSS())
	r.Path("/profile.css").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(base.ProfileCSS())
	r.Path("/history.css").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(base.HistoryCSS())
	r.Path("/invalid.css").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(base.InvalidCSS())

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
				Resolvers: graph.NewResolver(),
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}

package main

import (
	"compress/gzip"
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/clinical/graph"
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

	// check if must hace env variable exists
	// expects the server to die if this not explictly set
	base.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	base.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	base.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")

	// start up the router
	r, err := graph.Router()
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

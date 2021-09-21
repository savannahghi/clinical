package main

// TODO: restore
// import (
// 	"context"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"time"

// 	log "github.com/sirupsen/logrus"

// 	"github.com/savannahghi/clinical/graph"
// 	"github.com/savannahghi/serverutils"
// )

// const waitSeconds = 30

// func init() {
// 	// check if must have env variables exist
// 	// expects the server to die if this not explicitly set
// 	serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
// 	serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
// 	serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
// }

// func main() {
// 	ctx := context.Background()
// 	err := serverutils.Sentry()
// 	if err != nil {
// 		serverutils.LogStartupError(ctx, err)
// 	}

// 	port, err := strconv.Atoi(serverutils.MustGetEnvVar(serverutils.PortEnvVarName))
// 	if err != nil {
// 		serverutils.LogStartupError(ctx, err)
// 	}
// 	srv := graph.PrepareServer(ctx, port, graph.ClinicalAllowedOrigins)
// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil {
// 			serverutils.LogStartupError(ctx, err)
// 		}
// 	}()

// 	// Block until we receive a sigint (CTRL+C) signal
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	<-c

// 	// Create a deadline to wait for.
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*waitSeconds)
// 	defer cancel()

// 	// Doesn't block if no connections, but will otherwise wait until timeout
// 	err = srv.Shutdown(ctx)
// 	log.Printf("graceful shutdown started; the timeout is %d secs", waitSeconds)
// 	if err != nil {
// 		log.Printf("error during clean shutdown: %s", err)
// 		os.Exit(-1)
// 	}
// 	os.Exit(0)
// }

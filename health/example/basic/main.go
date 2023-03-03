package main

import (
	"net/http"

	"github.com/gardenbed/basil/health"
)

func main() {
	logger := &Logger{}
	dbClient := NewClient("db-client")
	queueClient := NewClient("queue-client")

	health.SetLogger(logger)
	health.RegisterChecker(dbClient, queueClient)

	logger.Infof("Listening on port 8080 ...")
	http.Handle("/health", health.HandlerFunc())
	_ = http.ListenAndServe(":8080", nil)
}

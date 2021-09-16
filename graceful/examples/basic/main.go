package main

import (
	"os"

	"github.com/gardenbed/basil/graceful"
)

func main() {
	logger := &Logger{}
	dbClient := NewClient("db-client")
	queueClient := NewClient("queue-client")
	apiServer := NewServer("api-server", 8080)
	infoServer := NewServer("info-server", 8081)

	graceful.SetLogger(logger)
	graceful.RegisterClient(dbClient, queueClient)
	graceful.RegisterServer(apiServer, infoServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}

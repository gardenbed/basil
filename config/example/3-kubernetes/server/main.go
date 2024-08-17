package main

import (
	"net/http"
	"sync"

	"github.com/gardenbed/basil/config"
	"github.com/gardenbed/basil/config/example/log"
)

var params = struct {
	sync.Mutex
	LogLevel string
}{
	// default value
	LogLevel: "info",
}

func main() {
	logger := new(log.Logger)
	logger.SetLevel("info")

	// Listening for any update to configurations
	ch := make(chan config.Update)
	go func() {
		for update := range ch {
			if update.Name == "LogLevel" {
				params.Lock()
				logger.SetLevel(params.LogLevel)
				params.Unlock()
			}
		}
	}()

	// Watching for configurations
	close, _ := config.Watch(&params, []chan config.Update{
		ch,
	})

	defer close()

	// HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("new request received")
		w.WriteHeader(http.StatusOK)
	})

	// Starting the HTTP server
	logger.Info("starting http server ...")
	_ = http.ListenAndServe(":8080", nil)
}

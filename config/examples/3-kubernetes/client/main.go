package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gardenbed/basil/config"
	"github.com/gardenbed/basil/config/examples/logger"
)

var params = struct {
	sync.Mutex
	LogLevel      string
	ServerAddress string
}{
	// default values
	LogLevel:      "info",
	ServerAddress: "http://localhost:8080",
}

func main() {
	logger := &logger.Logger{
		Level: logger.Info,
	}

	// Server address
	endpoint := "/"
	url := fmt.Sprintf("%s%s", params.ServerAddress, endpoint)

	// Listening for any update to configurations
	ch := make(chan config.Update)
	go func() {
		for update := range ch {
			switch update.Name {
			case "LogLevel":
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

	// Sending requests to server
	logger.Infof("start sending requests ...")

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}

		logger.Infof("response received from server: status_code: %d", resp.StatusCode)
	}
}

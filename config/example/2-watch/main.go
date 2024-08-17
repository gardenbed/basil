package main

import (
	"sync"
	"time"

	"github.com/gardenbed/basil/config"
	"github.com/gardenbed/basil/config/example/log"
)

// the single source of truth for all configurations
var params = struct {
	sync.Mutex
	LogLevel string
}{
	LogLevel: "info",
}

func main() {
	logger := new(log.Logger)
	logger.SetLevel("info")

	// Listening for configuration values and acting on them
	ch := make(chan config.Update, 1)
	go func() {
		for update := range ch {
			if update.Name == "LogLevel" {
				params.Lock()
				logger.SetLevel(params.LogLevel)
				params.Unlock()
			}
		}
	}()

	// Start watching for configurations values
	close, _ := config.Watch(&params, []chan config.Update{ch})
	defer close()

	// Simulate logging
	startLogging(logger)
}

func startLogging(logger *log.Logger) {
	wait := make(chan struct{}, 4)

	go func() {
		t1 := time.NewTicker(500 * time.Millisecond)
		for range t1.C {
			logger.Debug("Debugging ...")
		}
	}()

	go func() {
		t2 := time.NewTicker(1 * time.Second)
		for range t2.C {
			logger.Info("Informing ...")
		}
	}()

	go func() {
		t3 := time.NewTicker(2 * time.Second)
		for range t3.C {
			logger.Warn("Warning ...")
		}
	}()

	go func() {
		t4 := time.NewTicker(4 * time.Second)
		for range t4.C {
			logger.Error("Erroring ...")
		}
	}()

	<-wait
}

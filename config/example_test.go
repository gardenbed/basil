package config_test

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gardenbed/basil/config"
)

func ExamplePick() {
	// First, you need to define a struct.
	// Each field of the struct represents a configuration value.
	// Create an object from the struct with default values for its fields.
	var params = struct {
		LogLevel    string
		Environment string
		Region      string
		Timeout     time.Duration
		Replicas    []url.URL
	}{
		LogLevel:    "info",           // default
		Environment: "dev",            // default
		Region:      "local",          // default
		Timeout:     10 * time.Second, // default
	}

	// Second, pass the pointer to the struct object to the Pick method.
	// For each field, a value will be read either from flags, environment variables, or files.
	_ = config.Pick(&params)

	// Now, you can access the configuration values on the struct object.
	fmt.Printf("%+v\n", params)
}

func ExampleWatch() {
	// When using the Watch method, your struct needs to implement the sync.Locker interface.
	// You can simply achieve that by embedding the sync.Mutex type in your struct.
	var params = struct {
		sync.Mutex
		LogLevel string
	}{
		LogLevel: "info", // default
	}

	// For using the Watch method, you need to define a channel for receiving updates.
	// If a configuration value gets a new value (through files), you will get notified on this channel.
	ch := make(chan config.Update, 1)

	// In a separate goroutine, you can receive the new configuration values and re-configure your application accordingly.
	go func() {
		for update := range ch {
			if update.Name == "LogLevel" {
				params.Lock()
				fmt.Printf("Log level is updated to %s", params.LogLevel)
				params.Unlock()
			}
		}
	}()

	// You can now watch for configuration values.
	close, _ := config.Watch(&params, []chan config.Update{ch})
	defer close()
}

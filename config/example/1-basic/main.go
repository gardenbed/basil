package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/gardenbed/basil/config"
)

var params = struct {
	Port      uint16
	LogLevel  string
	Timeout   time.Duration
	Endpoints []url.URL
}{
	Port:     8080,            // default port
	LogLevel: "info",          // default log level
	Timeout:  2 * time.Minute, // default API call timeout
}

func main() {
	_ = config.Pick(&params)
	flag.Parse()

	fmt.Printf("\nParameters: %+v\n\n", params)
}

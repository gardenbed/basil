package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gardenbed/basil/config"
)

var params = struct {
	AuthToken string
}{}

func main() {
	_ = config.Pick(&params, config.Telepresence(), config.Debug(3))
	log.Printf("making service-to-service calls using this token: %s", params.AuthToken)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, World!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

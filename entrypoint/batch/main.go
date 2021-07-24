package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gae-go-sample/util/yaml"

	"gae-go-sample/di"
)

func main() {
	if os.Getenv("IS_LOCAL") == "true" {
		yaml.MustLoadLocalEnv("/app/entrypoint/batch/app.yaml")
	}

	mux := http.DefaultServeMux
	di.ResolveBatchHandler()(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("running server on port: %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("failed running server, err=%+v", err)
	}
}
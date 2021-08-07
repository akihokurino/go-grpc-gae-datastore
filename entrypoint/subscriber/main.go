package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gae-go-recruiting-server/util/yaml"

	"gae-go-recruiting-server/di"
)

func main() {
	if os.Getenv("IS_LOCAL") == "true" {
		yaml.MustLoadLocalEnv("/app/entrypoint/subscriber/app.yaml")
	}

	mux := http.DefaultServeMux
	di.ResolveSubscriberHandler()(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("running server on port: %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("failed running server, err=%+v", err)
	}
}

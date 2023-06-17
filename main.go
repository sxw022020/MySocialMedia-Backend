package main

import (
	"MySocialMedia-Backend/backend"
	"MySocialMedia-Backend/handler"
	"MySocialMedia-Backend/util"

	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Service started!")

	// load application configuration from deploy.yml
	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
	if err != nil {
		panic(err)
	}

	// Do the Elasticsearch Initialization when starting the program
	backend.InitElasticsearchBackend(config.ElasticsearchConfig)

	// Do the GCS Initialization when starting the program
	backend.InitGCSBackend(config.GCSConfig)

	/**
	An HTTP server is started,
	1. listening on port 8080 (request sent from Postman) and
	2. using the router initialized by the handler.InitRouter() function to handle incoming requests.
	3. If the server encounters any errors, the log.Fatal function is used to log the error message and exit the program with a non-zero status code.
	*/
	// The `config.TokenConfig` argument passed to `InitRouter` suggests that your routing configuration may depend on some sort of token information (possibly for authentication).
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig)))
}

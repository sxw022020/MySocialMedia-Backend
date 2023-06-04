package main

import (
	"MySocialMedia-Backend/backend"
	"MySocialMedia-Backend/handler"

	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Service started!")

	// Do the Elasticsearch Initialization when starting the program
	backend.InitElasticsearchBackend()

	/**
	An HTTP server is started,
	1. listening on port 8080 and
	2. using the router initialized by the handler.InitRouter() function to handle incoming requests.
	3. If the server encounters any errors, the log.Fatal function is used to log the error message and exit the program with a non-zero status code.
	*/
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}

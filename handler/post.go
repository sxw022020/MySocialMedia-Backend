package handler

import (
	"MySocialMedia-Backend/model"
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. Handle incoming requests for uploading posts
/*
`w http.ResponseWriter`:
	- An interface that allows you to send an HTTP response to the client
`r *http.Request`:
	- A pointer to an HTTP request object containing all the information about the incoming request
*/
func postUploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse from body of request to get a json object
	fmt.Println("Recieved one post uploading request!")

	// creates a JSON decoder using json.NewDecoder(r.Body) to parse the request body
	decoder := json.NewDecoder(r.Body)

	// initializes a `model.Post` struct variable `p`
	var p model.Post

	// decodes the JSON request body into the `model.Post` struct using `decoder Decode(&p)`
	// If there is an error during decoding, the function will panic and stop execution
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	// writes a response to the client using `fmt.Fprintf(w, "Post received: %s\n", p.Message)`, including the message from the uploaded post
	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

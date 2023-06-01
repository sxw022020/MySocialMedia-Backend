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

/*
1. Why pass reference as a request?
  - to avoid deep copy the request, just shallow copy
    --> improve the performance

2. Why pass an object as a response?
  - `ResponseWrite` is an interface which could be declared, not be implemented
  - When used as a parameter in a function or method,
    it allows the function to work with any concrete implementation
    of the ResponseWriter interface, which makes the code more modular and flexible.
  - `net/http“ package takes care of creating a concrete implementation of the `ResponseWriter`
    --> you don't need to worry about the specific implementation details and
    can focus on handling the HTTP request and constructing the response
*/
func postUploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse from body of request to get a json object
	fmt.Println("Recieved 1 post uploading request!")

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

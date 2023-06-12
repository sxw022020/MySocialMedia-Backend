package handler

import (
	"MySocialMedia-Backend/model"
	"MySocialMedia-Backend/service"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
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

	p := model.Post{
		Id:      uuid.New().String(),
		User:    r.FormValue("user"),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %w\n", err)
		return
	}

	fmt.Println("Post is saved successfully")

	/* Original Implementation, left here for review:

	// creates a JSON decoder using json.NewDecoder(r.Body) to parse the request body
	decoder := json.NewDecoder(r.Body)

	// initializes a `model.Post` struct variable `p`
	var p model.Post

	// decodes the JSON request body into the `model.Post` struct using `decoder.Decode(&p)`
	// If there is an error during decoding, the function will panic and stop execution
	if err := decoder.Decode(&p); err != nil {
		fmt.Println("Error: ", err)
	}

	// writes a response to the client using `fmt.Fprintf(w, "Post received: %s\n", p.Message)`, including the message from the uploaded post
	fmt.Fprintf(w, "Post received: %s\n", p.Message)

	*/
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}

package handler

import (
	"github.com/gorilla/mux"
)

// Unlike Spring, go cannot suport annotation-based HTTP routing,
// so we need our own method to map request URL to HTTP handler

/*
Return Type: `*mux.Router` - a pointer to a `mux.Router` instance.
*/
func InitRouter() *mux.Router {
	// Create a new Gorilla Mux router instance
	router := mux.NewRouter()

	// 1. **Register** the `postUploadHandler` function as an HTTP handler for the "/upload" endpoint with the "POST" method
	router.HandleFunc("/upload", postUploadHandler).Methods("POST")

	return router
}

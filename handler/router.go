package handler

import (
	"github.com/gorilla/mux"
)

/*
*
Return Type: `*mux.Router` - a pointer to a `mux.Router` instance.
*/
func InitRouter() *mux.Router {
	// Create a new Gorilla Mux router instance
	router := mux.NewRouter()

	// 1. **Register** the `postUploadHandler` function as an HTTP handler for the "/upload" endpoint with the "POST" method
	router.HandleFunc("/upload", postUploadHandler).Methods("POST")

	return router
}

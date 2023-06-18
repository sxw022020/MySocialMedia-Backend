package handler

import (
	"MySocialMedia-Backend/util"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var mySigningKey []byte

func InitRouter(config *util.TokenInfo) http.Handler {
	// `[]byte(config.Secret)` converts the string `config.Secret` into a slice of bytes.
	mySigningKey = []byte(config.Secret)
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()

	// Endpoints not requiring authentication
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/signin", signinHandler).Methods("POST")

	// Apply the middleware to endpoints requiring authentication
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(jwtMiddleware.Handler)

	authRouter.HandleFunc("/upload", postUploadHandler).Methods("POST")
	authRouter.HandleFunc("/search", postSearchHandler).Methods("GET")

	// accept requests from any origin
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	// applies the CORS rules to the router and returns it
	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}

package handler

import (
	"MySocialMedia-Backend/model"
	"MySocialMedia-Backend/service"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

// `w`: an interface allowing you to form a response to the request
// `r`: a pointer to the request received by the server
func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	// sets the "Content-Type" of the HTTP response to "text/plain".
	w.Header().Set("Content-Type", "text/plain")

	// If the HTTP method of the request is "OPTIONS" (which is typically used for CORS preflight requests), the function returns immediately.
	// This is done because the OPTIONS requests do not require further processing.
	if r.Method == "OPTIONS" {
		return
	}

	// Get `User` information from client, read the request body (which is assumed to be JSON) into `model.User`
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}

	exists, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
		return
	}

	// A conditional check in Go that checks if the exists variable is `false`
	//     - `service.CheckUser` return `boolean, error`
	if !exists {
		http.Error(w, "User doesn't exist or wrong password", http.StatusUnauthorized)
		fmt.Printf("User doesn't exist or wrong password\n")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// adds a claim where the claim key is "username" and the claim value is the username of the user who is signing in. This means that the JWT will include information about the user's username.
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// signs the token using the secret key "your-secret-key" and returns the final JWT in its string form.
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v\n", err)
		return
	}

	//  If no error occurred, this line sends the JWT to the client in the response body. This JWT should be used by the client for subsequent authenticated requests.
	_, err = w.Write([]byte(tokenString))
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "text/plain")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}

	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password\n")
		return
	}

	success, err := service.AddUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
		return
	}

	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Println("User already exists")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v\n", err)
		return
	}

	w.Write([]byte(tokenString))
}

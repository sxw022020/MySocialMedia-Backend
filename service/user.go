package service

import (
	"MySocialMedia-Backend/backend"
	"MySocialMedia-Backend/constants"
	"MySocialMedia-Backend/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// CheckUser : validate a user's credentials against an Elasticsearch index.
func CheckUser(username, password string) (bool, error) {

	// creating a map in Go, which represents an Elasticsearch boolean query, which can be converted to JSON.
	// The query will search for documents where both the "username" field matches the given username and the "password" field matches the given password.
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				// "must": This represents a list of conditions that must all be true for a document to match the query.
				"must": []map[string]interface{}{
					// "term" query, which looks for an exact match of a term (in this case, the username and password) in a specific field.
					{"term": map[string]interface{}{"username": username}},
					{"term": map[string]interface{}{"password": password}},
				},
			},
		},
	}

	// The query map is encoded into JSON and stored in a bytes.Buffer
	//     - A Buffer is a variable-sized buffer of bytes with Read and Write methods. It's part of Go's standard library.
	var buf bytes.Buffer
	//     - `json.NewEncoder(&buf)` creates a new JSON encoder that writes to buf. The encoder takes a pointer to a bytes.Buffer, which is why &buf is used.
	//     - `Encode(query)` encodes the query map into JSON format and writes it into buf.
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return false, err
	}

	// creates an Elasticsearch Search API request
	searchRequest := esapi.SearchRequest{
		Index:          []string{constants.USER_INDEX},
		Body:           &buf,
		TrackTotalHits: true, // track the total number of hits (matching results)
		// Setting TrackTotalHits to true instructs Elasticsearch to track and return the exact total number of hits, no matter how large that number is.
	}

	// executes the search request using the Elasticsearch client
	searchResponse, err := searchRequest.Do(context.Background(), backend.ESBackend.Client)
	if err != nil {
		return false, err
	}
	// ensures the response body stream is properly closed after we're done with it.
	// make sure the response body stream is closed as soon as the CheckUser function finishes, whether it finishes by returning normally or due to an error.
	defer searchResponse.Body.Close()

	if searchResponse.IsError() {
		return false, fmt.Errorf("Error getting response: %s", searchResponse.String())
	}

	// This line takes the body of the HTTP response from Elasticsearch (which is in JSON format), and decodes it into the res map.
	//     - initializes a new map to hold the response data
	var res map[string]interface{}
	//     - The `json.NewDecoder` function returns a new decoder that reads from searchResponse.Body
	//     -  The `Decode(&res)` method then reads the next JSON-encoded value from its input and stores it in the res map.
	if err := json.NewDecoder(searchResponse.Body).Decode(&res); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	}

	// Extract hits, which is an array of documents that matched our search criteria
	//    - We use type assertions (.(map[string]interface{}) and .([]interface{})) to convert the interface{} values to the actual types that we expect. This is necessary because Go is statically typed and does not automatically convert between types.
	hits := res["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		// For each "hit", we're interested in the "_source" field.
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		u := model.User{
			// access the value associated with the key "username" in the source map, and assert that it is a string
			Username: source["username"].(string),
			Password: source["password"].(string),
			// Add any other fields from your user model as necessary
		}

		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}

	return false, nil
}

func AddUser(user *model.User) (bool, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{"username": user.Username},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return false, err
	}

	searchRequest := esapi.SearchRequest{
		Index:          []string{constants.USER_INDEX},
		Body:           &buf,
		TrackTotalHits: true,
	}

	searchResponse, err := searchRequest.Do(context.Background(), backend.ESBackend.Client)
	if err != nil {
		return false, err
	}
	defer searchResponse.Body.Close()

	// Decode searchResponse
	var r map[string]interface{}
	if err := json.NewDecoder(searchResponse.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return false, err
	}

	// Extract total hits
	hits, found := r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	if !found {
		return false, fmt.Errorf("Error parsing total hits")
	}

	// If user already exists, return false
	if hits > 0 {
		return false, nil
	}

	data, err := json.Marshal(user)
	if err != nil {
		return false, err
	}

	indexRequest := esapi.IndexRequest{
		Index:      constants.USER_INDEX,
		DocumentID: user.Username,
		Body:       bytes.NewReader(data),
		Refresh:    "true", // Refresh the index to make the document immediately searchable
	}

	indexResponse, err := indexRequest.Do(context.Background(), backend.ESBackend.Client)
	if err != nil {
		return false, err
	}
	defer indexResponse.Body.Close()

	if indexResponse.IsError() {
		return false, fmt.Errorf("Error indexing document: %s", indexResponse.String())
	}

	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}

package backend

import (
	"MySocialMedia-Backend/constants"
	"MySocialMedia-Backend/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"net/http"
)

/**
This file: a Go program for interacting with an Elasticsearch backend
**/

// ESBackend a global variable `ESBackend` of type `*ElasticsearchBackend`.
// This is the Elasticsearch client that will be used to communicate with the Elasticsearch server.
var (
	// ESBackend a pointer of ElasticsearchBackend
	ESBackend *ElasticsearchBackend
)

// ElasticsearchBackend represents a backend service for Elasticsearch.
// It contains a pointer to an `elastic`.
// `Client` object which is a client connection to an Elasticsearch server.
type ElasticsearchBackend struct {
	// client a pointer of es7.Client
	Client *es7.Client
}

// InitElasticsearchBackend 0. Initialization of ElasticsearchBackend
func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	fmt.Println("Initialization of Elasticsearch!")

	cfg := es7.Config{
		Addresses: []string{
			config.Address,
		},
		Username: config.Username,
		Password: config.Password,
	}
	client, err := es7.NewClient(cfg)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	/*
		Once the client is successfully created,
		it's encapsulated in an ElasticsearchBackend instance and
		assigned to the global ESBackend.
	*/
	ESBackend = &ElasticsearchBackend{
		Client: client,
	}

	fmt.Println("Creation of Index!")
	// mapping: JSON string that defines the structure of the index
	// "keyword": `keyword`` fields are only searchable by their exact value
	// 	- if you want to search for a specific user by username or password,
	//    you'll need to provide the **exact** username or password,
	//	  cannot use partial info to do search
	// `"index": false`: it cannot be searched, but can still appear in the results
	checkAndCreateIndex(constants.POST_INDEX, `{
		"mappings": {
			"properties": {
				"id": 		{ "type": "keyword" },
				"user": 	{ "type": "keyword" },
				"message": 	{ "type": "text" },
				"url": 		{ "type": "keyword", "index": false },
				"type": 	{ "type": "keyword", "index": false }
			}
		}
	}`)

	checkAndCreateIndex(constants.USER_INDEX, `{
		"mappings": {
			"properties": {
				"username": 	{ "type": "keyword" },
				"password": 	{ "type": "keyword" },
				"age": 			{ "type": "long", "index": false },
				"gender": 		{ "type": "keyword", "index": false }
			}
		}
	}`)
}

// ReadFromES The function takes two arguments: query (a map of type string to interface{}) and index (a string).
// - The query represents the Elasticsearch query in a structured format, and the index is the name of the Elasticsearch index to search.
func (backend *ElasticsearchBackend) ReadFromES(query map[string]interface{}, index string) (*esapi.Response, error) {
	var (
		res *esapi.Response
		err error
	)

	// `Marshal(query)` takes Go data, in this case, `query`, and converts it to some other format (like JSON or XML).
	// It is commonly used when you want to send data over a network or save it to a file.
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// creating an instance of the `esapi.SearchRequest` struct from the Elasticsearch Go client library
	searchRequest := esapi.SearchRequest{
		// `index`: This field represents the name of the Elasticsearch index where the search will be performed
		Index:  []string{index},
		Body:   bytes.NewReader(queryJSON),
		Pretty: true,
	}

	// `backend.Client`:
	// This is an HTTP client that sends the request.
	// The client handles making the request and returning the response.
	// The client is usually created during the initialization of the Elasticsearch client.
	res, err = searchRequest.Do(context.Background(), backend.Client)
	// The Do method of searchRequest returns two values:
	//    - An `*http.Response` which is a pointer to a `http.Response` struct that contains the server's response to the HTTP request.
	//    - An `error` which is a built-in interface type for representing error conditions.

	if err != nil {
		return nil, err
	}

	return res, nil
}

////////// Helper Function //////////

/*
This function checks:
if a particular Elasticsearch index exists and
if it doesn't, creates it.
*/
func checkAndCreateIndex(indexName string, mapping string) {
	// check if the index exists
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), ESBackend.Client)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// `res.StatusCode == 404` means that the index does not exist
	// `http.StatusNotFound` is 404
	if res.StatusCode == http.StatusNotFound {
		// Index does not exist, so create it
		// - creating a new `IndicesCreateRequest` struct from the esapi package
		req := esapi.IndicesCreateRequest{
			Index: indexName,
			// sets the body of the request
			// - converting the mapping (which is a string) into a byte slice and
			//   then creating a new Reader for those bytes.
			//   This Reader can then be read by the Elasticsearch client to send the data
			Body: bytes.NewReader([]byte(mapping)),
		} // executing the request to create the index
		// - using the Do method of the IndicesCreateRequest struct,
		//   which sends the HTTP request to the Elasticsearch instance and
		//   returns the response.
		res, err := req.Do(context.Background(), ESBackend.Client)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// checks if the HTTP response indicates an error
		// (i.e., if the HTTP status code is 4xx or 5xx).
		if res.IsError() {
			//  declares a variable `e`` of type `map[string]interface{}``.
			// This map will be used to store the JSON error response from the Elasticsearch instance
			var e map[string]interface{}
			// reading the body of the error response from the Elasticsearch instance,
			// decoding the JSON into the e map, and checking if there was an error doing so
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				fmt.Println("Error: ", err)
			}

			// `%v` verb to print the entire error map
			fmt.Printf("[%s] Error creating index, response: %v", res.Status(), e)
		}

		err = res.Body.Close()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Indexes are created!")
	} else {
		fmt.Println("Index: " + indexName + " has already existed!")
	}
}

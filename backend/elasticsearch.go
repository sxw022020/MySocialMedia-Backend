package backend

import (
	"MySocialMedia-Backend/constants"
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
func InitElasticsearchBackend() {
	fmt.Println("Initialization of Elasticsearch!")

	cfg := es7.Config{
		Addresses: []string{
			constants.ES_URL,
		},
		Username: constants.ES_USERNAME,
		Password: constants.ES_PASSWORD,
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
	ESBackend = &ElasticsearchBackend{Client: client}

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
// - `map[string]interface{}`: a map of type string to interface{}
func (backend *ElasticsearchBackend) ReadFromES(query map[string]interface{}, index string) (*esapi.Response, error) {
	var (
		res *esapi.Response
		err error
	)

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	searchRequest := esapi.SearchRequest{
		Index:  []string{index},
		Body:   bytes.NewReader(queryJSON),
		Pretty: true,
	}

	res, err = searchRequest.Do(context.Background(), backend.Client)

	if err != nil {
		return nil, err
	}

	return res, nil
}

////////// Helper Functions //////////

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

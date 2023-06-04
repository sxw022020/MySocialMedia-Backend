package backend

import (
	"MySocialMedia-Backend/constants"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

/**
This file: a Go program for interacting with an Elasticsearch backend
**/

// Declares a global variable ESBackend of type *ElasticsearchBackend.
// This is the Elasticsearch client that will be used to communicate with the Elasticsearch server.
var (
	ESBackend *ElasticsearchBackend
)

// `ElasticsearchBackend struct` represents a backend service for Elasticsearch.
// It contains a pointer to an `elastic`.
// `Client“ object which is a client connection to an Elasticsearch server.
type ElasticsearchBackend struct {
	Client *es7.Client
}

func InitElasticsearchBackend() {
	cfg := es7.Config{
		Addresses: []string{
			constants.ES_URL,
		},
		Username: constants.ES_USERNAME,
		Password: constants.ES_PASSWORD,
	}
	client, err := es7.NewClient(cfg)
	if err != nil {
		/*
			When the panic function is called,
			it immediately stops the execution of the current function and
			begins unwinding the stack of the goroutine,
			running any deferred functions along the way.
			If that unwinding reaches the top of the goroutine's stack, the program dies.
		*/
		panic(err)
	}

	/*
		Once the client is successfully created,
		it's encapsulated in an ElasticsearchBackend instance and
		assigned to the global ESBackend.
	*/
	ESBackend = &ElasticsearchBackend{Client: client}

	// mapping: JSON string that defines the structure of the index
	// "keyword": `keyword`` fields are only searchable by their exact value
	// 	- if you want to search for a specific user by username or password,
	//    you'll need to provide the **exact** username or password,
	//	  cannot use partial info to do search
	// `"index": false`: it cannot be searched, but can still sppear in the results
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

	fmt.Println("Indexes are created!")
}

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
		panic(err)
	}
	// halting execution and reports the error
	// then close the body
	defer res.Body.Close()

	// `res.StatusCode == 404` means that the index does not exist
	// `http.StatusNotFound` is 404
	if res.StatusCode == 404 {
		// Index does not exist, so create it
		// - creating a new `IndicesCreateRequest` struct from the esapi package
		req := esapi.IndicesCreateRequest{
			Index: indexName,
			// sets the body of the request
			// - converting the mapping (which is a string) into a byte slice and
			//   then creating a new Reader for those bytes.
			//   This Reader can then be read by the Elasticsearch client to send the data
			Body: bytes.NewReader([]byte(mapping)),
		}

		// executing the request to create the index
		// - using the Do method of the IndicesCreateRequest struct,
		//   which sends the HTTP request to the Elasticsearch instance and
		//   returns the response.
		res, err := req.Do(context.Background(), ESBackend.Client)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close() // defer the Close() call here

		// checks if the HTTP response indicates an error
		// (i.e., if the HTTP status code is 4xx or 5xx).
		if res.IsError() {
			//  declares a variable `e`` of type `map[string]interface{}``.
			// This map will be used to store the JSON error response from the Elasticsearch instance
			var e map[string]interface{}
			// reading the body of the error response from the Elasticsearch instance,
			// decoding the JSON into the e map, and checking if there was an error doing so
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				panic(err)
			}

			fmt.Printf("[%s] Error indexing document ID=%d", res.Status(), e)
		}
	}
}

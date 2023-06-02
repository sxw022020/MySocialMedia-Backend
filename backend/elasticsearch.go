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
		panic(err)
	}

	ESBackend = &ElasticsearchBackend{Client: client}

	checkAndCreateIndex(constants.POST_INDEX, `{
		"mappings": {
			"properties:" {
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
			"properties:" {
				"username": 	{ "type": "keyword" },
				"password": 	{ "type": "keyword" },
				"age": 			{ "type": "long", "index": false },
				"gender": 		{ "type": "keyword", "index": false }
			}
		}
	}`)

	fmt.Println("Indexes are created!")
}

func checkAndCreateIndex(indexName string, mapping string) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), ESBackend.Client)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// Index does not exist, so create it
		req := esapi.IndicesCreateRequest{
			Index: indexName,
			Body:  bytes.NewReader([]byte(mapping)),
		}

		res, err := req.Do(context.Background(), ESBackend.Client)
		if err != nil {
			panic(err)
		}

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				panic(err)
			}

			fmt.Printf("[%s] Error indexing document ID=%d", res.Status(), e)
		}

		res.Body.Close()
	}
}

package backend

import (
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

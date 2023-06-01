package backend

import "github.com/olivere/elastic/v7"

var (
	ESBackend *ElasticsearchBackend
)

// `ElasticsearchBackend struct` represents a backend service for Elasticsearch.
// It contains a pointer to an `elastic`.
// `Client“ object which is a client connection to an Elasticsearch server.
type ElasticsearchBackend struct {
	Client *elastic.Client
}

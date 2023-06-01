package backend

import "github.com/olivere/elastic/v7"

var (
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
	Client *elastic.Client
}

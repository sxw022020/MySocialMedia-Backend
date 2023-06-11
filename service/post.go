package service

import (
	"MySocialMedia-Backend/constants"
	"MySocialMedia-Backend/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
)

func SearchPostsByUser(user string) ([]model.Post, error) {
	var es = &elasticsearch.Client{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"user": user,
			},
		},
	}
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(constants.POST_INDEX),
		es.Search.WithBody(strings.NewReader(buf.String())),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(res), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	var es = &elasticsearch.Client{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"user": keywords,
			},
		},
	}

	if keywords == "" {
		query["query"].(map[string]interface{})["term"].(map[string]interface{})["user"] = "all"
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(constants.POST_INDEX),
		es.Search.WithBody(strings.NewReader(buf.String())),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(res), nil
}

func getPostFromSearchResult(res *esapi.Response) []model.Post {
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	}

	var posts []model.Post
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var post model.Post
		postMap := hit.(map[string]interface{})["_source"].(map[string]interface{})
		postBytes, _ := json.Marshal(postMap)
		err := json.Unmarshal(postBytes, &post)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		posts = append(posts, post)
	}
	return posts
}

package service

import (
	"MySocialMedia-Backend/backend"
	"MySocialMedia-Backend/constants"
	"MySocialMedia-Backend/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func SearchPostsByUser(user string) ([]model.Post, error) {

	es := backend.ESBackend.Client

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

	req := esapi.SearchRequest{
		Index: []string{constants.POST_INDEX},
		Body:  strings.NewReader(buf.String()),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] Error searching posts by user: %v", res.Status(), e)
	}

	return getPostFromSearchResult(res), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	es := backend.ESBackend.Client

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"user": keywords,
			},
		},
	}

	if keywords == "" {
		query["query"].(map[string]interface{})["match"].(map[string]interface{})["user"] = "all"
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	req := esapi.SearchRequest{
		Index: []string{constants.POST_INDEX},
		Body:  strings.NewReader(buf.String()),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] Error searching posts by keywords: %v", res.Status(), e)
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

func SavePost(post *model.Post, file multipart.File) error {
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink

	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}

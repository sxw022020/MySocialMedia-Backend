package backend

import (
	"MySocialMedia-Backend/constants"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	Client *storage.Client
	Bucket string
}

func InitGCSBackend() {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Println("Error: ", err)
	}

	GCSBackend = &GoogleCloudStorageBackend{
		Client: client,
		Bucket: constants.GCS_BUCKET,
	}
}

func (gcsbackend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := gcsbackend.Client.Bucket(GCSBackend.Bucket).Object(objectName)

	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
	var b strings.Builder
	err := json.NewEncoder(&b).Encode(i)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(b.String()),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), backend.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error: %s", res.String())
	}

	return nil
}

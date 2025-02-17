// Package elastic provides a client for interacting with Elasticsearch.
package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client represents an Elasticsearch client.
type Client struct {
	es *elasticsearch.Client
}

var (
	instance *Client
	once     sync.Once
)

// GetClient returns the singleton instance of the Elasticsearch client.
func GetClient(url string) *Client {
	once.Do(func() {
		es, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{url},
		})
		if err != nil {
			log.Fatal(err)
		}

		instance = &Client{es: es}
	})

	return instance
}

// Get retrieves a document from Elasticsearch.
func (c *Client) Get(index string, id string) (*esapi.Response, error) {
	res, err := c.es.Get(index, id, c.es.Get.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create creates a new document in Elasticsearch.
func (c *Client) Create(index string, id string, doc interface{}) (*esapi.Response, error) {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	docWithTimestamp := struct {
		Timestamp string      `json:"@timestamp"`
		Data      interface{} `json:"data"`
	}{
		Timestamp: timestamp,
		Data:      doc,
	}

	data, err := json.Marshal(docWithTimestamp)
	if err != nil {
		return nil, err
	}

	res, err := c.es.Index(index, bytes.NewReader(data), c.es.Index.WithDocumentID(id), c.es.Index.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update updates an existing document in Elasticsearch.
func (c *Client) Update(index string, id string, doc interface{}) (*esapi.Response, error) {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	docWithTimestamp := struct {
		Timestamp string      `json:"@timestamp"`
		Data      interface{} `json:"data"`
	}{
		Timestamp: timestamp,
		Data:      doc,
	}

	data, err := json.Marshal(docWithTimestamp)
	if err != nil {
		return nil, err
	}

	res, err := c.es.Update(index, id, bytes.NewReader(data), c.es.Update.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete deletes a document from Elasticsearch.
func (c *Client) Delete(index string, id string) (*esapi.Response, error) {
	res, err := c.es.Delete(index, id, c.es.Delete.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}

	return res, nil
}

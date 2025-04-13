package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"

	config "factorysphere/devicesphere/engines/notificationengine/pkg/config"
	log "factorysphere/devicesphere/engines/notificationengine/pkg/utils/loggers"
)

// ElasticsearchRepository represents an Elasticsearch repository.
type ElasticsearchRepository struct {
	client *elasticsearch.Client
}

// NewElasticsearchRepository creates a new Elasticsearch repository.
func NewElasticsearchRepository(config *config.Config) (*ElasticsearchRepository, error) {
	// Create Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: config.Elasticsearch.Addresses,
	}

	// Add authentication if provided
	if config.Elasticsearch.Username != "" && config.Elasticsearch.Password != "" {
		cfg.Username = config.Elasticsearch.Username
		cfg.Password = config.Elasticsearch.Password
	}

	// Add CloudID if provided
	if config.Elasticsearch.CloudID != "" {
		cfg.CloudID = config.Elasticsearch.CloudID
	}

	// Add APIKey if provided
	if config.Elasticsearch.APIKey != "" {
		cfg.APIKey = config.Elasticsearch.APIKey
	}

	// Create Elasticsearch client
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Test the connection
	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error connecting to Elasticsearch: %s", res.String())
	}

	log.Debug("Connected to Elasticsearch!")

	return &ElasticsearchRepository{client: client}, nil
}

// Create creates a new document in the specified index.
func (r *ElasticsearchRepository) Create(index string, data map[string]interface{}) (string, error) {
	// Generate a random ID if not provided
	id, ok := data["id"].(string)
	if !ok || id == "" {
		id = uuid.New().String()
		data["id"] = id
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create the document
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true", // Ensure the document is immediately available for search
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("error indexing document: %s", res.String())
	}

	// Parse the response
	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		return "", err
	}

	log.Debug("Document created with ID: %s", id)
	return id, nil
}

// Get retrieves a document by ID from the specified index.
func (r *ElasticsearchRepository) Get(index string, id string) (map[string]interface{}, error) {
	// Get the document
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, nil // Document not found
		}
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	// Parse the response
	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		return nil, err
	}

	// Extract the document source
	source, ok := resMap["_source"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format")
	}

	return source, nil
}

// Update updates an existing document in the specified index.
func (r *ElasticsearchRepository) Update(index string, id string, data map[string]interface{}) error {
	// Prepare the update document
	updateDoc := map[string]interface{}{
		"doc": data,
	}

	// Convert update document to JSON
	jsonData, err := json.Marshal(updateDoc)
	if err != nil {
		return err
	}

	// Update the document
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true", // Ensure the update is immediately visible
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}

	log.Debug("Document updated with ID: %s", id)
	return nil
}

// Delete deletes a document by ID from the specified index.
func (r *ElasticsearchRepository) Delete(index string, id string) error {
	// Delete the document
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true", // Ensure the deletion is immediately visible
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	log.Debug("Document deleted with ID: %s", id)
	return nil
}

// Search searches for documents in the specified index.
func (r *ElasticsearchRepository) Search(index string, query map[string]interface{}) ([]map[string]interface{}, error) {
	// Convert query to JSON
	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Search for documents
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  bytes.NewReader(jsonQuery),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	// Parse the response
	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		return nil, err
	}

	// Extract the hits
	hits, err := extractHits(resMap)
	if err != nil {
		return nil, err
	}

	return hits, nil
}

// BulkCreate creates multiple documents in the specified index.
func (r *ElasticsearchRepository) BulkCreate(index string, documents []map[string]interface{}) error {
	if len(documents) == 0 {
		return nil
	}

	// Prepare the bulk request body
	var buf bytes.Buffer
	for _, doc := range documents {
		// Generate a random ID if not provided
		id, ok := doc["id"].(string)
		if !ok || id == "" {
			id = uuid.New().String()
			doc["id"] = id
		}

		// Add the index action
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
				"_id":    id,
			},
		}

		// Add the metadata line
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}

		// Add the document line
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return err
		}
	}

	// Execute the bulk request
	res, err := r.client.Bulk(bytes.NewReader(buf.Bytes()), r.client.Bulk.WithIndex(index), r.client.Bulk.WithRefresh("true"))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in bulk operation: %s", res.String())
	}

	// Parse the response to check for errors
	var bulkResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return err
	}

	// Check if there were any errors in the items
	if hasErrors, ok := bulkResponse["errors"].(bool); ok && hasErrors {
		return errors.New("some documents failed to be indexed")
	}

	log.Debug("Bulk create operation completed successfully")
	return nil
}

// DeleteByQuery deletes documents matching the query in the specified index.
func (r *ElasticsearchRepository) DeleteByQuery(index string, query map[string]interface{}) error {
	// Convert query to JSON
	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return err
	}

	// Delete documents by query
	req := esapi.DeleteByQueryRequest{
		Index: []string{index},
		Body:  bytes.NewReader(jsonQuery),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting documents by query: %s", res.String())
	}

	log.Debug("Documents deleted by query")
	return nil
}

// CreateIndex creates a new index with the specified settings and mappings.
func (r *ElasticsearchRepository) CreateIndex(index string, settings map[string]interface{}) error {
	// Convert settings to JSON
	jsonSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// Create the index
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  bytes.NewReader(jsonSettings),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// If the error is because the index already exists, we can ignore it
		if strings.Contains(res.String(), "resource_already_exists_exception") {
			log.Debug("Index %s already exists", index)
			return nil
		}
		return fmt.Errorf("error creating index: %s", res.String())
	}

	log.Debug("Index %s created", index)
	return nil
}

// DeleteIndex deletes the specified index.
func (r *ElasticsearchRepository) DeleteIndex(index string) error {
	// Delete the index
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// If the error is because the index doesn't exist, we can ignore it
		if res.StatusCode == 404 {
			log.Debug("Index %s does not exist", index)
			return nil
		}
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	log.Debug("Index %s deleted", index)
	return nil
}

// Helper function to extract hits from search response
func extractHits(response map[string]interface{}) ([]map[string]interface{}, error) {
	hitsObj, ok := response["hits"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format: missing hits object")
	}

	hitsArray, ok := hitsObj["hits"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: missing hits array")
	}

	var results []map[string]interface{}
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Add the document ID to the source
		if id, ok := hitMap["_id"].(string); ok {
			source["id"] = id
		}

		results = append(results, source)
	}

	return results, nil
}

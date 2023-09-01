package main

import (
	"context"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Replace with your Elasticsearch server address
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	indexName := "braingain2"

	// Query using cosine similarity
	query := `
	{
		"query": {
			"script_score": {
				"query": {
					"match_all": {}
				},
				"script": {
					"source": "cosineSimilarity(params.queryVector, 'vector_field') + 1.0",
					"params": {
						"queryVector": [0.1, 0.2, 0.0]
					}
				}
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error querying index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error querying index: %s", res.String())
	}

	log.Println("Query result:")
	log.Println(res.String())
}

//package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/elastic/go-elasticsearch/v8"
//	"github.com/elastic/go-elasticsearch/v8/esapi"
//	"log"
//	"strings"
//)
//
//func main() {
//	log.SetFlags(log.LstdFlags | log.Lshortfile)
//
//	cfg := elasticsearch.Config{
//		Addresses: []string{"http://localhost:9200"}, // Replace with your Elasticsearch server address
//	}
//
//	client, err := elasticsearch.NewClient(cfg)
//	if err != nil {
//		log.Fatalf("Error creating the client: %s", err)
//	}
//
//	indexName := "braingain2"
//
//	mapping := `
//	{
//		"mappings": {
//			"properties": {
//				"vector_field": {
//					"type": "dense_vector",
//					"dims": 3
//				}
//			}
//		}
//	}`
//
//	req := esapi.IndicesCreateRequest{
//		Index: indexName,
//		Body:  strings.NewReader(mapping),
//	}
//
//	res, err := req.Do(context.Background(), client)
//	if err != nil {
//		log.Fatalf("Error creating index: %s", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		log.Fatalf("Error creating index: %s", res.String())
//	}
//
//	fmt.Println("Index mapping created successfully.")
//
//	// Add data to the index
//	document := `
//	{
//		"vector_field": [0.1, 0.2, 0.3]
//	}`
//
//	req1 := esapi.IndexRequest{
//		Index:      indexName,
//		DocumentID: "1", // Replace with the desired document ID
//		Body:       strings.NewReader(document),
//		Refresh:    "true",
//	}
//
//	res, err = req1.Do(context.Background(), client)
//	if err != nil {
//		log.Fatalf("Error indexing document: %s", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		log.Fatalf("Error indexing document: %s", res.String())
//	}
//
//	log.Println("Document indexed successfully.")
//}

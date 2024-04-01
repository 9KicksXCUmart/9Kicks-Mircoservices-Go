package dao

import (
	"9Kicks/config"
	"9Kicks/model/product"
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

var (
	client    = config.GetOpenSearchClient()
	indexName = "ninekicks_products"
)

func Search(from, size int, keyword string) (product.SearchResponse, error) {
	var searchResponse product.SearchResponse
	searchQuery := product.SearchQuery{
		From: from,
		Size: size,
		Query: struct {
			Match struct {
				Name string `json:"name"`
			} `json:"match"`
		}{
			Match: struct {
				Name string `json:"name"`
			}{
				Name: keyword,
			},
		},
	}

	searchQueryBytes, _ := json.Marshal(searchQuery)
	query := strings.NewReader(string(searchQueryBytes))

	// Search for products
	search := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body:  query,
	}

	resp, _ := search.Do(context.Background(), client)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return searchResponse, err
	}
	respString := buf.String()
	searchResponse, err = loadSearchResponse(respString)

	return searchResponse, err
}

func loadSearchResponse(dataString string) (product.SearchResponse, error) {
	var searchResponse product.SearchResponse
	err := json.Unmarshal([]byte(dataString), &searchResponse)
	if err != nil {
		return searchResponse, err
	}
	return searchResponse, nil
}

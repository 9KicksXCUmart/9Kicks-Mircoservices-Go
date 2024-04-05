package dao

import (
	"9Kicks/config"
	"9Kicks/model/product"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

var (
	client    = config.GetOpenSearchClient()
	indexName = "ninekicks_products"
)

func Filter(from, size int, boolQuery product.BoolQuery) (product.SearchResponse, error) {
	var searchResponse product.SearchResponse
	searchQuery := product.SearchQuery{From: from, Size: size, Query: boolQuery}

	searchQueryBytes, _ := json.Marshal(searchQuery)
	query := strings.NewReader(string(searchQueryBytes))

	log.Println(string(searchQueryBytes))

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

func CreateProduct(productInfo product.ProductInfo) error {
	docId := productInfo.ID
	productInfoBytes, _ := json.Marshal(productInfo)

	log.Println(string(productInfoBytes))

	document := strings.NewReader(string(productInfoBytes))
	req := opensearchapi.IndexRequest{
		Index:      indexName,
		DocumentID: docId,
		Body:       document,
	}

	insertResponse, err := req.Do(context.Background(), client)

	log.Println(insertResponse)

	if err != nil {
		return err
	}
	return nil
}

func GetProductDetailByID() {
	// TODO: Get product detail by ID
}

func loadSearchResponse(dataString string) (product.SearchResponse, error) {
	var searchResponse product.SearchResponse
	err := json.Unmarshal([]byte(dataString), &searchResponse)
	if err != nil {
		return searchResponse, err
	}
	return searchResponse, nil
}

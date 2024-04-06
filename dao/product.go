package dao

import (
	"9Kicks/config"
	"9Kicks/model/product"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

var (
	client       = config.GetOpenSearchClient()
	indexName    = "ninekicks_products"
	sizeFieldMap = map[string]string{
		"5.5":  "Size5_5",
		"6.0":  "Size6_0",
		"6.5":  "Size6_5",
		"7.0":  "Size7_0",
		"7.5":  "Size7_5",
		"8.0":  "Size8_0",
		"8.5":  "Size8_5",
		"9.0":  "Size9_0",
		"9.5":  "Size9_5",
		"10.0": "Size10_0",
		"10.5": "Size10_5",
		"11.0": "Size11_0",
		"11.5": "Size11_5",
		"12.0": "Size12_0",
		"12.5": "Size12_5",
		"13.0": "Size13_0",
		"14.0": "Size14_0",
	}
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

func UpdateProduct(productInfo product.ProductInfo) error {
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

func GetProductDetailByID(productId string) (product.DocumentResponse, error) {
	var documentResponse product.DocumentResponse
	resp, err := client.Get(indexName, productId)
	if err != nil {
		return documentResponse, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return documentResponse, err
	}
	respString := buf.String()
	documentResponse, err = loadDocumentResponse(respString)

	return documentResponse, err
}

func DeleteProduct(productId string) error {
	req := opensearchapi.DeleteRequest{
		Index:      indexName,
		DocumentID: productId,
	}

	resp, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}

func GetRemainingStock(productId string, size string) (int, error) {
	documentResponse, err := GetProductDetailByID(productId)
	if err != nil {
		return 0, err
	}

	productInfo := documentResponse.Source
	// Get the field name from the map
	fieldName, ok := sizeFieldMap[size]
	if !ok {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	// Use reflection to get the corresponding field
	sizeStruct := reflect.ValueOf(&productInfo.Size).Elem()
	field := sizeStruct.FieldByName(fieldName)
	if field.IsValid() {
		return int(field.Int()), nil
	}

	return 0, fmt.Errorf("invalid field: %s", fieldName)
}

func UpdateStock(productId string, size string, sold int) error {
	documentResponse, err := GetProductDetailByID(productId)
	if err != nil {
		return err
	}

	productInfo := documentResponse.Source
	// Get the field name from the map
	fieldName, ok := sizeFieldMap[size]
	if !ok {
		return fmt.Errorf("invalid size: %s", size)
	}

	// Use reflection to update the corresponding field
	sizeStruct := reflect.ValueOf(&productInfo.Size).Elem()
	field := sizeStruct.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		if field.Int() < int64(sold) {
			return fmt.Errorf("not enough stock for size: %s", size)
		}
		field.SetInt(field.Int() - int64(sold))
	} else {
		return fmt.Errorf("invalid field: %s", fieldName)
	}

	// Update the buy count
	productInfo.BuyCount += sold

	err = UpdateProduct(productInfo)
	if err != nil {
		return err
	}

	return nil
}

func loadSearchResponse(dataString string) (product.SearchResponse, error) {
	var searchResponse product.SearchResponse
	err := json.Unmarshal([]byte(dataString), &searchResponse)
	if err != nil {
		return searchResponse, err
	}
	return searchResponse, nil
}

func loadDocumentResponse(dataString string) (product.DocumentResponse, error) {
	var documentResponse product.DocumentResponse
	err := json.Unmarshal([]byte(dataString), &documentResponse)
	if err != nil {
		return documentResponse, err
	}
	return documentResponse, nil
}

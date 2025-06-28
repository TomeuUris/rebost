package main

import (
	"encoding/json"
	"net/http"
)

// OpenFoodFactsClient is a client for the Open Food Facts API.
// It is used to get information about products.
type OpenFoodFactsClient struct {
	BaseURL string
}

// NewOpenFoodFactsClient creates a new OpenFoodFactsClient.
func NewOpenFoodFactsClient() *OpenFoodFactsClient {
	return &OpenFoodFactsClient{
		BaseURL: "https://es.openfoodfacts.org",
	}
}

// GetProduct gets a product from the Open Food Facts API.
// It takes a barcode as input and returns an APIProduct struct.
func (c *OpenFoodFactsClient) GetProduct(barcode string) (*APIProduct, error) {
	url := c.BaseURL + "/api/v0/product/" + barcode + ".json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Product APIProduct `json:"product"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Product, nil
}

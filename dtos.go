package main

import "time"

// ProductDTO represents a product in the API.
type ProductDTO struct {
	Barcode     string        `json:"barcode"`
	ProductName string        `json:"product_name"`
	Brand       string        `json:"brand"`
	Ingredients string        `json:"ingredients"`
	Nutriments  NutrimentsDTO `json:"nutriments"`
}

// NutrimentsDTO represents the nutrition facts of a product in the API.
type NutrimentsDTO struct {
	Per100g    *RationDTO `json:"per_100g,omitempty"`
	PerServing *RationDTO `json:"per_serving,omitempty"`
}

// RationDTO represents the nutrition facts of a product for a specific ration in the API.
type RationDTO struct {
	EnergyKcal    float64 `json:"energy-kcal"`
	Fat           float64 `json:"fat"`
	SaturatedFat  float64 `json:"saturated-fat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Sugars        float64 `json:"sugars"`
	Fiber         float64 `json:"fiber"`
	Proteins      float64 `json:"proteins"`
	Salt          float64 `json:"salt"`
	Sodium        float64 `json:"sodium"`
}

// InventoryItemDTO represents an inventory item in the API.
type InventoryItemDTO struct {
	ID             uint       `json:"id"`
	Product        ProductDTO `json:"product"`
	ExpirationDate time.Time  `json:"expiration_date"`
	Quantity       int        `json:"quantity"`
	CreatedAt      time.Time  `json:"created_at"`
}

package main

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// InventoryItem represents an item in the inventory.
// It contains information about the product and the quantity.
type InventoryItem struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	ProductID      uint
	Product        Product
	ExpirationDate time.Time `json:"expiration_date"`
	Quantity       int       `json:"quantity"`
}

// InventoryItemRequest is used to bind the request body when adding a new inventory item.
type InventoryItemRequest struct {
	Barcode        string `json:"barcode"`
	ExpirationDate string `json:"expiration_date"`
	Quantity       int    `json:"quantity"`
}

// Product is a struct that represents a product from the Open Food Facts API.
// It contains information about the product, such as its name, brand, and ingredients.
type Product struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Barcode      string         `gorm:"unique"`
	ProductName  string
	Brand        string
	Ingredients  string
	Nutriments   Nutriments
	NutrimentsID uint
}

// APIProduct is a struct that represents a product from the Open Food Facts API.
type APIProduct struct {
	Barcode     string     `json:"id"`
	ProductName string     `json:"product_name"`
	Brand       string     `json:"brands"`
	Ingredients string     `json:"ingredients_text"`
	Nutriments  Nutriments `json:"nutriments"`
}

// Nutriments is a struct that represents the nutrition facts of a product.
// It contains information about the product's energy, fat, carbohydrates, and protein.
type Nutriments struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Per100g      *Ration        `json:"per_100g,omitempty"`
	Per100gID    *uint
	PerServing   *Ration `json:"per_serving,omitempty"`
	PerServingID *uint
}

// Ration is a struct that represents the nutrition facts of a product for a specific ration.
// It contains information about the product's energy, fat, carbohydrates, and protein.
type Ration struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	EnergyKcal    float64        `json:"energy-kcal"`
	Fat           float64        `json:"fat"`
	SaturatedFat  float64        `json:"saturated-fat"`
	Carbohydrates float64        `json:"carbohydrates"`
	Sugars        float64        `json:"sugars"`
	Fiber         float64        `json:"fiber"`
	Proteins      float64        `json:"proteins"`
	Salt          float64        `json:"salt"`
	Sodium        float64        `json:"sodium"`
}

// UnmarshalJSON is a custom unmarshaler for the Nutriments struct.
// It is used to parse the data from the Open Food Facts API into the new structure.
func (n *Nutriments) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	getFloat := func(key string) float64 {
		if val, ok := raw[key]; ok {
			if fVal, ok := val.(float64); ok {
				return fVal
			}
		}
		return 0
	}

	if _, ok := raw["energy-kcal_100g"]; ok {
		n.Per100g = &Ration{
			EnergyKcal:    getFloat("energy-kcal_100g"),
			Fat:           getFloat("fat_100g"),
			SaturatedFat:  getFloat("saturated-fat_100g"),
			Carbohydrates: getFloat("carbohydrates_100g"),
			Sugars:        getFloat("sugars_100g"),
			Fiber:         getFloat("fiber_100g"),
			Proteins:      getFloat("proteins_100g"),
			Salt:          getFloat("salt_100g"),
			Sodium:        getFloat("sodium_100g"),
		}
	}

	if _, ok := raw["energy-kcal_serving"]; ok {
		n.PerServing = &Ration{
			EnergyKcal:    getFloat("energy-kcal_serving"),
			Fat:           getFloat("fat_serving"),
			SaturatedFat:  getFloat("saturated-fat_serving"),
			Carbohydrates: getFloat("carbohydrates_serving"),
			Sugars:        getFloat("sugars_serving"),
			Fiber:         getFloat("fiber_serving"),
			Proteins:      getFloat("proteins_serving"),
			Salt:          getFloat("salt_serving"),
			Sodium:        getFloat("sodium_serving"),
		}
	}

	return nil
}

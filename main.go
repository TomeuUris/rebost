package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "rebost/docs"
)

// @title Rebost API
// @version 1.0
// @description This is a sample server for a food inventory tracker.
// @host localhost:8080
// @BasePath /

type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Get a product by its barcode
// @Description get product by barcode
// @ID get-product-by-barcode
// @Accept  json
// @Produce  json
// @Param   barcode      path   string     true  "Product Barcode"
// @Success 200 {object} ProductDTO
// @Failure 500 {object} ErrorResponse
// @Router /v1/product/{barcode} [get]
func getProduct(c *gin.Context, client *OpenFoodFactsClient) {
	barcode := c.Param("barcode")
	apiProduct, err := client.GetProduct(barcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToProductDTOFromAPIProduct(*apiProduct))
}

func main() {
	r := gin.Default()

	db, err := InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	offenFoodFactsClient := NewOpenFoodFactsClient()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/product/:barcode", func(c *gin.Context) {
			getProduct(c, offenFoodFactsClient)
		})
		v1.POST("/inventory", func(c *gin.Context) {
			addInventoryItem(c, db, offenFoodFactsClient)
		})
		v1.GET("/inventory", func(c *gin.Context) {
			getInventory(c, db)
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// @Summary Add an item to the inventory
// @Description add an item to the inventory
// @ID add-inventory-item
// @Accept  json
// @Produce  json
// @Param   inventoryItem   body    InventoryItemRequest   true  "Inventory Item"
// @Success 200 {object} InventoryItemDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/inventory [post]
func addInventoryItem(c *gin.Context, db *gorm.DB, client *OpenFoodFactsClient) {
	var req InventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	expirationDate, err := time.Parse("2006/01/02", req.ExpirationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date format for expiration_date. Use YYYY/MM/DD."})
		return
	}

	// Check if product exists
	var product Product
	if err := db.First(&product, "barcode = ?", req.Barcode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Product not found, fetch from Open Food Facts
			apiProduct, err := client.GetProduct(req.Barcode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch product from Open Food Facts"})
				return
			}

			// Create product in db
			product = ToProductFromAPIProduct(*apiProduct)

			if err := db.Create(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save product"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
			return
		}
	}

	// Create inventory item
	inventoryItem := InventoryItem{
		ProductID:      product.ID,
		Product:        product,
		ExpirationDate: expirationDate,
		Quantity:       req.Quantity,
	}

	if err := db.Create(&inventoryItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save inventory item"})
		return
	}

	c.JSON(http.StatusOK, ToInventoryItemDTO(inventoryItem))
}

// @Summary Get the inventory
// @Description get all inventory items
// @ID get-inventory
// @Produce  json
// @Success 200 {array} InventoryItemDTO
// @Failure 500 {object} ErrorResponse
// @Router /v1/inventory [get]
func getInventory(c *gin.Context, db *gorm.DB) {
	var inventoryItems []InventoryItem
	if err := db.Preload("Product.Nutriments.Per100g").Preload("Product.Nutriments.PerServing").Preload("Product").Find(&inventoryItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve inventory"})
		return
	}

	var inventoryDTOs []InventoryItemDTO
	for _, item := range inventoryItems {
		inventoryDTOs = append(inventoryDTOs, ToInventoryItemDTO(item))
	}

	c.JSON(http.StatusOK, inventoryDTOs)
}
